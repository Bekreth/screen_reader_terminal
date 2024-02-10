package terminal

import (
	"fmt"
	"io"
	"testing"

	"github.com/bekreth/screen_reader_terminal/buffer"
	"github.com/bekreth/screen_reader_terminal/history"
	"github.com/bekreth/screen_reader_terminal/utils"
	"github.com/bekreth/screen_reader_terminal/window"
	"github.com/stretchr/testify/assert"
)

// test interface implmenetations

type testFile struct {
	written []byte
}

func (window *testFile) Write(input []byte) (int, error) {
	window.written = append(window.written, input...)
	return len(input), nil
}

// tests
func TestDrawLine(t *testing.T) {
	trials := []struct {
		description      string
		previousData     string
		previousPosition int
		currentValue     string
		currentPosition  int
		expectedOutput   string
	}{
		{
			description:      "Two empty strings",
			previousData:     "",
			previousPosition: 0,
			currentValue:     "",
			currentPosition:  0,
			expectedOutput:   "",
		},
		{
			description:      "Current string is appended",
			previousData:     "hello worl",
			previousPosition: 10,
			currentValue:     "hello world",
			currentPosition:  11,
			expectedOutput:   "d",
		},
		{
			description:      "Current string is appended, cursor not at end",
			previousData:     "hello worl",
			previousPosition: 6,
			currentValue:     "hello world",
			currentPosition:  6,
			expectedOutput:   "world",
		},
		{
			description:      "Current string has a character insert",
			previousData:     "hello world",
			previousPosition: 11,
			currentValue:     "hello HERE world",
			currentPosition:  6,
			expectedOutput:   "HERE world",
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			actualOutput := lineDiff(
				trial.previousData, trial.previousPosition,
				trial.currentValue, trial.currentPosition,
			)
			assert.Equal(tt, trial.expectedOutput, actualOutput)
		})
	}
}

type testWindow struct {
	window.Window
	io.Writer
}

func TestDraw(t *testing.T) {
	trials := []struct {
		description      string
		previousPrefix   string
		prefix           string
		currentValue     string
		currentPosition  int
		previousValue    string
		previousPosition int
		expectedOutput   string
	}{
		{
			description:      "Draw empty buffer",
			currentValue:     "",
			currentPosition:  0,
			previousValue:    "",
			previousPosition: 0,
			expectedOutput:   "",
		},
		{
			description:      "Draw new buffer",
			currentValue:     "Hello world",
			currentPosition:  11,
			previousValue:    "",
			previousPosition: 0,
			expectedOutput: fmt.Sprintf("%v",
				"Hello world",
			),
		},
		{
			description:      "Add character to buffer end",
			currentValue:     "Hello world",
			currentPosition:  11,
			previousValue:    "Hello worl",
			previousPosition: 10,
			expectedOutput: fmt.Sprintf("%v",
				"d",
			),
		},
		{
			description:      "Remove character from buffer end",
			currentValue:     "Hello worl",
			currentPosition:  10,
			previousValue:    "Hello world",
			previousPosition: 11,
			expectedOutput: fmt.Sprintf("%v%v",
				fmt.Sprintf("%v%v%v", window.CSI, 1, "D"),                     // Move cursor left by 1
				fmt.Sprintf("%v%v%v", window.CSI, window.CURSOR_FORWARD, "K"), // Remove the text
			),
		},
		{
			description:      "Insert text in middle of line",
			currentValue:     "Hello, world",
			currentPosition:  6,
			previousValue:    "Hello world",
			previousPosition: 5,
			expectedOutput: fmt.Sprintf("%v%v",
				", world",
				fmt.Sprintf("%v%v%v", window.CSI, 6, "D"), // Move cursor left by 3
			),
		},
		{
			description:      "Move cursor back 3",
			currentValue:     "Hello world",
			currentPosition:  7,
			previousValue:    "Hello world",
			previousPosition: 10,
			expectedOutput: fmt.Sprintf("%v%v%v",
				window.CSI, 3, "D", // Move cursor left by 3
			),
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			file := testFile{
				written: []byte{},
			}
			win := window.NewWindow().
				SetWindowSize(window.WindowSize{
					Height: 20,
					Width:  20,
				}).
				SetWriter(&file)

			buf := buffer.NewBuffer()
			buf.SetCurrentValues(buffer.BufferValues{
				Prefix:   trial.prefix,
				Value:    trial.currentValue,
				Position: trial.currentPosition,
			})
			buf.SetPreviousValues(buffer.BufferValues{
				Prefix:   trial.previousPrefix,
				Value:    trial.previousValue,
				Position: trial.previousPosition,
			})

			terminalUnderTest := NewTerminal(
				win,
				&buf,
				utils.TestLogger{
					TestPrefix: trial.description[0:10],
					Tester:     tt,
				},
			)
			newHistory := history.NewBufferHistory()
			terminalUnderTest.history = &newHistory

			terminalUnderTest.Draw()
			assert.Equal(tt, []byte(trial.expectedOutput), file.written)
		})
	}
}

func TestDetermineRows(t *testing.T) {
	veryLongLine := fmt.Sprintf("18. %v\n25. %v\n35. %v\n45. %v\n55. %v",
		"456789012345678",
		"4567890123456789012345",
		"45678901234567890123456789012345",
		"456789012345678901234567890123456789012345",
		"45678901234567890123456789012345678901234567890",
	)

	trials := []struct {
		description string

		previousPrefix   string
		prefix           string
		previousValue    string
		previousPosition int
		currentValue     string
		currentPosition  int

		expectedRow    int
		expectedCursor int
		expectedOffset int
	}{
		{
			description:      "Empty buffer, should see 1,1",
			prefix:           "",
			previousPrefix:   "",
			currentPosition:  0,
			currentValue:     "",
			previousPosition: 0,
			previousValue:    "",
			expectedRow:      1,
			expectedCursor:   1,
			expectedOffset:   0,
		},
		{
			description:      "no buffer but short prefix",
			prefix:           "short",
			previousPrefix:   "",
			currentPosition:  0,
			currentValue:     "",
			previousPosition: 0,
			previousValue:    "",
			expectedRow:      1,
			expectedCursor:   1,
			expectedOffset:   5,
		},
		{
			description:      "short string with no prefix",
			prefix:           "",
			previousPrefix:   "",
			currentPosition:  0,
			currentValue:     "small string",
			previousPosition: 12,
			previousValue:    "",
			expectedRow:      1,
			expectedCursor:   1,
			expectedOffset:   0,
		},
		{
			description:      "cursor in front of long prefix/string",
			prefix:           "small: ",
			previousPrefix:   "",
			currentPosition:  0,
			currentValue:     "012345678901234567890",
			previousPosition: 0,
			previousValue:    "",
			expectedRow:      2,
			expectedCursor:   1,
			expectedOffset:   7,
		},
		{
			description:      "cursor at end of long prefix/string",
			prefix:           "small: ",
			previousPrefix:   "",
			currentPosition:  18,
			currentValue:     "012345678901234567890",
			previousPosition: 0,
			previousValue:    "",
			expectedRow:      2,
			expectedCursor:   2,
			expectedOffset:   5,
		},
		{
			description:      "cursor in the middle of long prefix/string",
			prefix:           "small: ",
			previousPrefix:   "",
			currentPosition:  18,
			currentValue:     "01234567890123456789012345678901234567890",
			previousPosition: 0,
			previousValue:    "",
			expectedRow:      3,
			expectedCursor:   2,
			expectedOffset:   5,
		},
		{
			description:      "cursor end of very long line",
			prefix:           "small: ",
			previousPrefix:   "",
			currentPosition:  40,
			currentValue:     "01234567890123456789012345678901234567890",
			previousPosition: 0,
			previousValue:    "",
			expectedRow:      3,
			expectedCursor:   3,
			expectedOffset:   7,
		},
		{
			description:      "0. several small lines, cursor at start",
			prefix:           "small: ",
			previousPrefix:   "",
			currentPosition:  0,
			currentValue:     "hello\nworld\ncheck\nme",
			previousPosition: 0,
			previousValue:    "",
			expectedRow:      4,
			expectedCursor:   1,
			expectedOffset:   7,
		},
		{
			description:      "1. several small lines, cursor at second word",
			prefix:           "small: ",
			previousPrefix:   "",
			currentPosition:  8,
			currentValue:     "hello\nworld\ncheck\nme",
			previousPosition: 0,
			previousValue:    "",
			expectedRow:      4,
			expectedCursor:   2,
			expectedOffset:   2,
		},
		{
			description:      "2. several small lines, cursor at third word",
			prefix:           "small: ",
			previousPrefix:   "",
			currentPosition:  15,
			currentValue:     "hello\nworld\ncheck\nme",
			previousPosition: 0,
			previousValue:    "",
			expectedRow:      4,
			expectedCursor:   3,
			expectedOffset:   3,
		},
		{
			description:      "3. several small lines, cursor at forth word",
			prefix:           "small: ",
			previousPrefix:   "",
			currentPosition:  18,
			currentValue:     "hello\nworld\ncheck\nme",
			previousPosition: 0,
			previousValue:    "",
			expectedRow:      4,
			expectedCursor:   4,
			expectedOffset:   0,
		},
		{
			description:      "A. Several very long lines, first word",
			prefix:           "small: ",
			previousPrefix:   "",
			currentPosition:  10,
			currentValue:     veryLongLine,
			previousPosition: 0,
			previousValue:    "",
			expectedRow:      12,
			expectedCursor:   1,
			expectedOffset:   17,
		},
		{
			description:      "B. Several very long lines, end of first word",
			prefix:           "small: ",
			previousPrefix:   "",
			currentPosition:  17,
			currentValue:     veryLongLine,
			previousPosition: 0,
			previousValue:    "",
			expectedRow:      12,
			expectedCursor:   2,
			expectedOffset:   4,
		},
		{
			description:      "C. Several very long lines, start of second word",
			prefix:           "small: ",
			previousPrefix:   "",
			currentPosition:  25,
			currentValue:     veryLongLine,
			previousPosition: 0,
			previousValue:    "",
			expectedRow:      12,
			expectedCursor:   3,
			expectedOffset:   5,
		},
		{
			description:      "D. Several very long lines, end of second word",
			prefix:           "small: ",
			previousPrefix:   "",
			currentPosition:  45,
			currentValue:     veryLongLine,
			previousPosition: 0,
			previousValue:    "",
			expectedRow:      12,
			expectedCursor:   4,
			expectedOffset:   5,
		},
		{
			description:      "E. Several very long lines, start of third word",
			prefix:           "small: ",
			previousPrefix:   "",
			currentPosition:  47,
			currentValue:     veryLongLine,
			previousPosition: 0,
			previousValue:    "",
			expectedRow:      12,
			expectedCursor:   5,
			expectedOffset:   0,
		},
		{
			description:      "F. Several very long lines, middle of third word",
			prefix:           "small: ",
			previousPrefix:   "",
			currentPosition:  69,
			currentValue:     veryLongLine,
			previousPosition: 0,
			previousValue:    "",
			expectedRow:      12,
			expectedCursor:   6,
			expectedOffset:   2,
		},
		{
			description:      "G. Several very long lines, end of third word",
			prefix:           "small: ",
			previousPrefix:   "",
			currentPosition:  82,
			currentValue:     veryLongLine,
			previousPosition: 0,
			previousValue:    "",
			expectedRow:      12,
			expectedCursor:   6,
			expectedOffset:   15,
		},
		{
			description:      "H. Several very long lines, middle of fourth word",
			prefix:           "small: ",
			previousPrefix:   "",
			currentPosition:  122,
			currentValue:     veryLongLine,
			previousPosition: 0,
			previousValue:    "",
			expectedRow:      12,
			expectedCursor:   8,
			expectedOffset:   18,
		},
		{
			description:      "I. Several very long lines, end of fifth word",
			prefix:           "small: ",
			previousPrefix:   "",
			currentPosition:  181,
			currentValue:     veryLongLine,
			previousPosition: 0,
			previousValue:    "",
			expectedRow:      12,
			expectedCursor:   12,
			expectedOffset:   10,
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			file := testFile{
				written: []byte{},
			}
			win := window.NewWindow().
				SetWriter(&file).
				SetWindowSize(window.WindowSize{
					Height: 20,
					Width:  20,
				})

			buf := buffer.NewBuffer()
			buf.SetCurrentValues(buffer.BufferValues{
				Prefix:   trial.prefix,
				Value:    trial.currentValue,
				Position: trial.currentPosition,
			})
			buf.SetPreviousValues(buffer.BufferValues{
				Prefix:   trial.previousPrefix,
				Value:    trial.previousValue,
				Position: trial.previousPosition,
			})

			terminalUnderTest := NewTerminal(
				win,
				&buf,
				utils.TestLogger{
					TestPrefix: trial.description[0:15],
					Tester:     tt,
				},
			)
			testingValue, testingCursor := terminalUnderTest.CurrentBuffer().Output()
			actualRow, actualCursor, actualOffset := terminalUnderTest.determineRows(
				testingValue,
				testingCursor,
			)

			assert.Equal(tt, trial.expectedRow, actualRow, "ROWS")
			assert.Equal(tt, trial.expectedCursor, actualCursor, "CUROSOR ROWS")
			assert.Equal(tt, trial.expectedOffset, actualOffset, "CUROSOR OFFSET")
		})
	}
}
