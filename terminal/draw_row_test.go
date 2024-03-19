package terminal

import (
	"testing"

	"github.com/bekreth/screen_reader_terminal/utils"
	"github.com/bekreth/screen_reader_terminal/window"
	"github.com/stretchr/testify/assert"
)

type drawRowTrial struct {
	description         string
	previousLineData    string
	currentLineData     string
	cursorCoordinates   coordinates
	expectedWrite       string
	expectedCoordinates coordinates
}

func TestDrawRow(t *testing.T) {
	trials := []drawRowTrial{
		{
			description:      "Simple append to row",
			previousLineData: "",
			currentLineData:  "h",
			cursorCoordinates: coordinates{
				currentX: 0,
				currentY: 0,
			},
			expectedWrite: "h",
			expectedCoordinates: coordinates{
				currentX:      1,
				currentY:      0,
				pendingDeltaX: 1,
				pendingDeltaY: 0,
			},
		},
		{
			description:      "Simple delete from row",
			previousLineData: "h",
			currentLineData:  "",
			cursorCoordinates: coordinates{
				currentX: 1,
				currentY: 0,
			},
			expectedWrite: fmtLine(
				left(1),
				clearCursorForward(),
			),
			expectedCoordinates: coordinates{
				currentX:      0,
				currentY:      0,
				pendingDeltaX: 0,
				pendingDeltaY: 0,
			},
		},
		{
			description:      "appending to existing string",
			previousLineData: "hello ",
			currentLineData:  "hello w",
			cursorCoordinates: coordinates{
				currentX: 6,
				currentY: 0,
			},
			expectedWrite: fmtLine(
				"w",
			),
			expectedCoordinates: coordinates{
				currentX:      7,
				currentY:      0,
				pendingDeltaX: 7,
				pendingDeltaY: 0,
			},
		},
		{
			description:      "deleting from end of existing string",
			previousLineData: "hello w",
			currentLineData:  "hello ",
			cursorCoordinates: coordinates{
				currentX: 7,
				currentY: 0,
			},
			expectedWrite: fmtLine(
				left(1),
				clearCursorForward(),
			),
			expectedCoordinates: coordinates{
				currentX:      6,
				currentY:      0,
				pendingDeltaX: 6,
				pendingDeltaY: 0,
			},
		},
		{
			description:      "inserting into string",
			previousLineData: "hello orld",
			currentLineData:  "hello world",
			cursorCoordinates: coordinates{
				currentX: 6,
				currentY: 0,
			},
			expectedWrite: fmtLine(
				"world",
			),
			expectedCoordinates: coordinates{
				currentX:      11,
				currentY:      0,
				pendingDeltaX: 11,
				pendingDeltaY: 0,
			},
		},
		{
			description:      "removing from  string",
			previousLineData: "hello world",
			currentLineData:  "hello orld",
			cursorCoordinates: coordinates{
				currentX: 6,
				currentY: 0,
			},
			expectedWrite: fmtLine(
				clearCursorForward(),
				"orld",
			),
			expectedCoordinates: coordinates{
				currentX:      10,
				currentY:      0,
				pendingDeltaX: 10,
				pendingDeltaY: 0,
			},
		},
		{
			description:      "Move to location then update",
			previousLineData: "hello orld",
			currentLineData:  "hello world",
			cursorCoordinates: coordinates{
				currentX: 0,
				currentY: 0,
			},
			expectedWrite: fmtLine(
				right(6),
				"world",
			),
			expectedCoordinates: coordinates{
				currentX:      11,
				currentY:      0,
				pendingDeltaX: 11,
				pendingDeltaY: 0,
			},
		},
		{
			description:      "Move to location from other row then update",
			previousLineData: "hello orld",
			currentLineData:  "hello world",
			cursorCoordinates: coordinates{
				currentX:      16,
				currentY:      1,
				pendingDeltaY: 0,
			},
			expectedWrite: fmtLine(
				left(10),
				up(1),
				"world",
			),
			expectedCoordinates: coordinates{
				currentX:      11,
				currentY:      0,
				pendingDeltaX: 11,
				pendingDeltaY: 0,
			},
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

			terminalUnderTest := Terminal{
				window: win,
				logger: utils.TestLogger{
					TestPrefix: trial.description[0:10],
					Tester:     tt,
				},
			}

			actualCoordinates := terminalUnderTest.drawRow(
				trial.previousLineData,
				trial.currentLineData,
				trial.cursorCoordinates,
			)
			assert.Equal(tt, []byte(trial.expectedWrite), file.written)
			assert.Equal(tt, trial.expectedCoordinates, actualCoordinates)
		})
	}
}

func TestRowDiff(t *testing.T) {
	trials := []struct {
		description    string
		previousRow    string
		currentRow     string
		expectedEnd    string
		expectedColumn int
	}{
		{
			description:    "empty strings going in",
			previousRow:    "",
			currentRow:     "",
			expectedEnd:    "",
			expectedColumn: 0,
		},
		{
			description:    "current data has simple appends",
			previousRow:    "",
			currentRow:     "a",
			expectedEnd:    "a",
			expectedColumn: 0,
		},
		{
			description:    "current data has simple backspace",
			previousRow:    "a",
			currentRow:     "",
			expectedEnd:    "",
			expectedColumn: 0,
		},
		{
			description:    "append after string",
			previousRow:    "hell",
			currentRow:     "hello",
			expectedEnd:    "o",
			expectedColumn: 4,
		},
		{
			description:    "remove from end of string",
			previousRow:    "hello",
			currentRow:     "hell",
			expectedEnd:    "",
			expectedColumn: 4,
		},
		{
			description:    "insert in the middle",
			previousRow:    "hello orld",
			currentRow:     "hello world",
			expectedEnd:    "world",
			expectedColumn: 6,
		},
		{
			description:    "remove from middle",
			previousRow:    "hello wworld",
			currentRow:     "hello world",
			expectedEnd:    "orld",
			expectedColumn: 7,
		},
		{
			description:    "insert string at end",
			previousRow:    "hello",
			currentRow:     "hello world",
			expectedEnd:    " world",
			expectedColumn: 5,
		},
		{
			description:    "remove string from end",
			previousRow:    "hello world",
			currentRow:     "hello",
			expectedEnd:    "",
			expectedColumn: 5,
		},
		{
			description:    "insert string in middle",
			previousRow:    "hello world",
			currentRow:     "hello this world",
			expectedEnd:    "this world",
			expectedColumn: 6,
		},
		{
			description:    "remove string from middle",
			previousRow:    "hello this world",
			currentRow:     "hello world",
			expectedEnd:    "world",
			expectedColumn: 6,
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			actualEnd, actualColumn := rowDiff(trial.previousRow, trial.currentRow)

			assert.Equal(tt, trial.expectedEnd, actualEnd)
			assert.Equal(tt, trial.expectedColumn, actualColumn)
		})
	}
}
