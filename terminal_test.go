package screen_reader_terminal

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
		buffer         Buffer
		expectedOutput string
	}{
		{
			description:    "Draw empty buffer",
			buffer:         NewBuffer(),
			expectedOutput: "",
		},
		{
			description: "Draw new buffer",
			buffer:      NewBufferWithString("Hello world"),
			expectedOutput: fmt.Sprintf("%v%v%v%v",
				fmt.Sprintf("%v%v", CSI, "s"),
				"Hello world",
				fmt.Sprintf("%v%v", CSI, "u"),
				fmt.Sprintf("%v%v%v", CSI, "11", "C"),
			),
		},
		{
			description: "Add character to buffer end",
			buffer: Buffer{
				currentValue:     "Hello world",
				cursorPosition:   11,
				previousValue:    "Hello worl",
				previousPosition: 10,
			},
			expectedOutput: fmt.Sprintf("%v%v%v%v",
				fmt.Sprintf("%v%v", CSI, "s"), // Save Cursor Position
				"d",
				fmt.Sprintf("%v%v", CSI, "u"),      // Load Cursor Position
				fmt.Sprintf("%v%v%v", CSI, 1, "C"), // Move cursor right by 1
			),
		},
		{
			description: "Remove character from buffer end",
			buffer: Buffer{
				currentValue:     "Hello worl",
				cursorPosition:   10,
				previousValue:    "Hello world",
				previousPosition: 11,
			},
			expectedOutput: fmt.Sprintf("%v%v%v%v%v",
				fmt.Sprintf("%v%v", CSI, "s"),                   // Save Cursor Position
				fmt.Sprintf("%v%v%v", CSI, 1, "D"),              // Move cursor left by 1
				fmt.Sprintf("%v%v%v", CSI, CURSOR_FORWARD, "K"), // Remove the text
				fmt.Sprintf("%v%v", CSI, "u"),                   // Save Cursor Position
				fmt.Sprintf("%v%v%v", CSI, 1, "D"),              // Move cursor left by 1
			),
		},
		{
			description: "Move cursor back 3",
			buffer: Buffer{
				currentValue:     "Hello world",
				cursorPosition:   7,
				previousValue:    "Hello world",
				previousPosition: 10,
			},
			expectedOutput: fmt.Sprintf("%v%v%v",
				CSI, 3, "D", // Move cursor left by 3
			),
		},
		{
			description: "Insert text in middle of line",
			buffer: Buffer{
				currentValue:     "Hello, world",
				cursorPosition:   6,
				previousValue:    "Hello world",
				previousPosition: 5,
			},
			expectedOutput: fmt.Sprintf("%v%v%v%v",
				fmt.Sprintf("%v%v", CSI, "s"), // Save Cursor Position
				", world",
				fmt.Sprintf("%v%v", CSI, "u"),      // Load Cursor Position
				fmt.Sprintf("%v%v%v", CSI, 1, "C"), // Move cursor right by 1
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
			terminalUnderTest.history = &newHistory

			terminalUnderTest.Draw()
			assert.Equal(tt, []byte(trial.expectedOutput), file.written)
		})
	}
}
