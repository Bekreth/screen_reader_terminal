package terminal

import (
	"fmt"
	"testing"

	"github.com/bekreth/screen_reader_terminal/buffer"
	"github.com/bekreth/screen_reader_terminal/utils"
	"github.com/bekreth/screen_reader_terminal/window"
	"github.com/stretchr/testify/assert"
)

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
			description:      "cursor end of short line with new line character",
			prefix:           "small: ",
			previousPrefix:   "",
			currentPosition:  12,
			currentValue:     "hello world\n",
			previousPosition: 11,
			previousValue:    "hello world",
			expectedRow:      2,
			expectedCursor:   2,
			expectedOffset:   0,
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
