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

type testWindow struct {
	window.Window
	io.Writer
}

func fmtLine(input ...any) string {
	format := ""
	for i := 0; i < len(input); i++ {
		format += "%v"
	}
	return fmt.Sprintf(format, input...)
}

func column(input int) string {
	return fmtLine(window.CSI, input, "G")
}

func up(input int) string {
	return fmtLine(window.CSI, input, "A")
}

func down(input int) string {
	return fmtLine(window.CSI, input, "B")
}

func left(input int) string {
	return fmtLine(window.CSI, input, "D")
}

func right(input int) string {
	return fmtLine(window.CSI, input, "C")
}

// tests

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
			expectedOutput: fmtLine(
				"Hello world",
			),
		},
		{
			description:      "Add character to buffer end",
			currentValue:     "Hello world",
			currentPosition:  11,
			previousValue:    "Hello worl",
			previousPosition: 10,
			expectedOutput: fmtLine(
				"d",
			),
		},
		{
			description:      "Remove character from buffer end",
			currentValue:     "Hello worl",
			currentPosition:  10,
			previousValue:    "Hello world",
			previousPosition: 11,
			expectedOutput: fmtLine(
				left(1),
				fmtLine(window.CSI, window.CURSOR_FORWARD, "K"), // Remove the text
			),
		},
		{
			description:      "Insert text in middle of line",
			currentValue:     "Hello, world",
			currentPosition:  6,
			previousValue:    "Hello world",
			previousPosition: 5,
			expectedOutput: fmtLine(
				", world",
				left(6),
			),
		},
		{
			description:      "Move cursor back 3",
			currentValue:     "Hello world",
			currentPosition:  7,
			previousValue:    "Hello world",
			previousPosition: 10,
			expectedOutput: fmtLine(
				left(3),
			),
		},
		{
			description:      "Move cursor forward over new line",
			currentValue:     "Hello\nworld",
			currentPosition:  11,
			previousValue:    "Hello\nworld",
			previousPosition: 4,
			expectedOutput: fmtLine(
				right(1),
				down(1),
			),
		},
		{
			description:      "Move cursor back over new line",
			currentValue:     "Hello\nworld",
			currentPosition:  4,
			previousValue:    "Hello\nworld",
			previousPosition: 11,
			expectedOutput: fmtLine(
				left(1),
				up(1),
			),
		},
		{
			description:      "Update short line to have new line",
			currentValue:     "Hello\n",
			currentPosition:  6,
			previousValue:    "Hello",
			previousPosition: 5,
			expectedOutput: fmtLine(
				down(1),
				column(0),
				fmtLine(window.CSI, window.FULL, "K"), // Remove the text
			),
		},
		{
			description:      "Remove new line from short line",
			currentValue:     "Hello",
			currentPosition:  5,
			previousValue:    "Hello\n",
			previousPosition: 6,
			expectedOutput: fmtLine(
				right(5),
				up(1),
			),
		},
		{
			description:      "Update short line to have new line and short line",
			currentValue:     "Hello\nWorld",
			currentPosition:  11,
			previousValue:    "Hello",
			previousPosition: 5,
			expectedOutput: fmtLine(
				down(1),
				column(0),
				fmtLine(window.CSI, window.FULL, "K"), // Remove the text
				"World",
			),
		},
		{
			description:      "Update last line in multi-line statement",
			currentValue:     "Hello\nworld",
			currentPosition:  11,
			previousValue:    "Hello\nworl",
			previousPosition: 10,
			expectedOutput: fmtLine(
				"d",
			),
		},
		{
			description:      "Update first line in multi-line statement",
			currentValue:     "Hello,\nworld",
			currentPosition:  6,
			previousValue:    "Hello\nworld",
			previousPosition: 5,
			expectedOutput: fmtLine(
				",",
			),
		},
		{
			description:      "Paste multiple lines",
			currentValue:     "Hello,\nworld\nmultiple",
			currentPosition:  20,
			previousValue:    "",
			previousPosition: 0,
			expectedOutput: fmtLine(
				"Hello,",
				down(1),
				column(0),
				fmtLine(window.CSI, window.FULL, "K"), // Remove the text
				"world",
				down(1),
				column(0),
				fmtLine(window.CSI, window.FULL, "K"), // Remove the text
				"multiple",
			),
		},
		{
			description:      "Long line before rollover",
			currentValue:     "12345678901234567890",
			currentPosition:  20,
			previousValue:    "1234567890123456789",
			previousPosition: 19,
			expectedOutput: fmtLine(
				"0",
			),
		},
		{
			description:      "Long line during rollover",
			currentValue:     "123456789012345678901",
			currentPosition:  21,
			previousValue:    "12345678901234567890",
			previousPosition: 20,
			expectedOutput: fmtLine(
				"1",
			),
		},
		{
			description:      "Long line after rollover",
			currentValue:     "1234567890123456789012",
			currentPosition:  22,
			previousValue:    "123456789012345678901",
			previousPosition: 21,
			expectedOutput: fmtLine(
				"2",
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
			his := history.NewBufferHistory()

			terminalUnderTest := Terminal{
				window:  win,
				buffer:  &buf,
				history: &his,
				logger: utils.TestLogger{
					TestPrefix: trial.description[0:10],
					Tester:     tt,
				},
			}

			terminalUnderTest.Draw()
			assert.Equal(tt, []byte(trial.expectedOutput), file.written)
		})
	}
}
