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
		{
			description:      "Move cursor back over new line",
			currentValue:     "Hello\nworld",
			currentPosition:  4,
			previousValue:    "Hello\nworld",
			previousPosition: 11,
			expectedOutput: fmt.Sprintf("%v%v",
				fmt.Sprintf("%v%v%v", window.CSI, 1, "D"), // Move cursor left by 3
				fmt.Sprintf("%v%v%v", window.CSI, 1, "B"), // Move cursor Up by 3
			),
		},
		{
			description:      "Update last line in multi-line statement",
			currentValue:     "Hello\nworld",
			currentPosition:  11,
			previousValue:    "Hello\nworl",
			previousPosition: 10,
			expectedOutput: fmt.Sprintf("%v",
				"d",
			),
		},
		{
			description:      "Update first line in multi-line statement",
			currentValue:     "Hello,\nworld",
			currentPosition:  6,
			previousValue:    "Hello\nworld",
			previousPosition: 5,
			expectedOutput: fmt.Sprintf("%v",
				",",
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
