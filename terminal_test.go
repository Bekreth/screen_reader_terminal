package main

import (
	"fmt"
	"testing"

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

type testLogger struct {
	*testing.T
}

func (l testLogger) Infof(format string, args ...interface{}) {
	l.Logf("INFO: %v", fmt.Sprintf(format, args...))
}
func (l testLogger) Debugf(format string, args ...interface{}) {
	l.Logf("DEBUG: %v", fmt.Sprintf(format, args...))
}

// tests

func TestDraw(t *testing.T) {
	trials := []struct {
		description    string
		previousBuffer Buffer
		buffer         Buffer
		expectedOutput string
	}{
		{
			description:    "Draw empty buffer",
			previousBuffer: NewBuffer(),
			buffer:         NewBuffer(),
			expectedOutput: "",
		},
		{
			description:    "Draw new buffer",
			previousBuffer: NewBuffer(),
			buffer:         NewBufferWithString("Hello world"),
			expectedOutput: fmt.Sprintf("%v%v%v%v",
				fmt.Sprintf("%v%v", CSI, "s"),
				"Hello world",
				fmt.Sprintf("%v%v", CSI, "u"),
				fmt.Sprintf("%v%v%v", CSI, "11", "C"),
			),
		},
		{
			description:    "Add character to buffer end",
			previousBuffer: NewBufferWithString("Hello worl"),
			buffer:         NewBufferWithString("Hello world"),
			expectedOutput: "d",
		},
		{
			description:    "Remove character from buffer end",
			previousBuffer: NewBufferWithString("Hello world"),
			buffer:         NewBufferWithString("Hello worl"),
			expectedOutput: fmt.Sprintf("%v%v",
				fmt.Sprintf("%v%v%v", CSI, 1, "D"),              // Move cursor left by 1
				fmt.Sprintf("%v%v%v", CSI, CURSOR_FORWARD, "K"), // Remove the text
			),
		},
		{
			description: "Move cursor back 3",
			previousBuffer: Buffer{
				currentValue:   "Hello world",
				cursorPosition: 10,
			},
			buffer: Buffer{
				currentValue:   "Hello world",
				cursorPosition: 7,
			},
			expectedOutput: fmt.Sprintf("%v%v%v",
				CSI, 3, "D", // Move cursor left by 3
			),
		},
		{
			description: "Insert text in middle of line",
			previousBuffer: Buffer{
				currentValue:   "Hello world",
				cursorPosition: 5,
			},
			buffer: Buffer{
				currentValue:   "Hello, world",
				cursorPosition: 6,
			},
			expectedOutput: fmt.Sprintf("%v%v%v%v",
				fmt.Sprintf("%v%v%v", CSI, FULL_LINE, "K"), // Remove the text
				fmt.Sprintf("%v%v%v", CSI, "0", "G"),       // Set Cursor position
				"Hello, world",
				fmt.Sprintf("%v%v%v", CSI, "7", "G"), // Set Cursor position
			),
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			file := testFile{
				written: []byte{},
			}
			window := unixWindow{
				size: WindowSize{
					height: 20,
					width:  80,
				},
				file: &file,
			}
			terminalUnderTest := NewTerminal(&window, &trial.buffer, testLogger{tt})
			newHistory := NewBufferHistory()
			newHistory.AddBuffer(trial.previousBuffer)
			terminalUnderTest.history = &newHistory

			terminalUnderTest.Draw()
			assert.Equal(tt, []byte(trial.expectedOutput), file.written)
		})
	}
}
