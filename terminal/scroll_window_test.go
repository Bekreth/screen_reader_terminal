package terminal

import (
	"testing"

	"github.com/bekreth/screen_reader_terminal/buffer"
	"github.com/bekreth/screen_reader_terminal/utils"
	"github.com/bekreth/screen_reader_terminal/window"
	"github.com/stretchr/testify/assert"
)

type scrollWindowTrial struct {
	description string

	cursorHeight     int
	previousRowCount int
	currentRowCount  int

	expectedCursorHeight int
	expectedWrite        string
}

func TestScrollWindow(t *testing.T) {
	trials := []scrollWindowTrial{
		{
			description:          "Write an empty row",
			cursorHeight:         0,
			previousRowCount:     0,
			currentRowCount:      0,
			expectedCursorHeight: 0,
			expectedWrite:        "",
		},
		{
			description:          "Write 1 line nowhere near end",
			cursorHeight:         0,
			previousRowCount:     0,
			currentRowCount:      1,
			expectedCursorHeight: 0,
			expectedWrite:        "",
		},
		{
			description:          "Write new 2 lines nowhere near end",
			cursorHeight:         0,
			previousRowCount:     0,
			currentRowCount:      2,
			expectedCursorHeight: 1,
			expectedWrite:        "",
		},
		{
			description:          "Rewrite same 2 lines",
			cursorHeight:         5,
			previousRowCount:     2,
			currentRowCount:      2,
			expectedCursorHeight: 5,
			expectedWrite:        "",
		},
		{
			description:          "Add 1 row at height limit",
			cursorHeight:         19,
			previousRowCount:     1,
			currentRowCount:      2,
			expectedCursorHeight: 19,
			expectedWrite:        fmtLine(scrollup(1), up(1)),
		},
		{
			description:          "Add 2 rows at height limit",
			cursorHeight:         19,
			previousRowCount:     1,
			currentRowCount:      3,
			expectedCursorHeight: 19,
			expectedWrite:        fmtLine(scrollup(2), up(2)),
		},
		{
			description:          "Add 10 rows near height limit",
			cursorHeight:         15,
			previousRowCount:     2,
			currentRowCount:      12,
			expectedCursorHeight: 19,
			expectedWrite:        fmtLine(scrollup(6), up(6)),
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

			logger := utils.TestLogger{
				TestPrefix: trial.description[0:utils.IntMin(len(trial.description), 15)],
				Tester:     tt,
			}

			terminalUnderTest := NewTerminal(
				win,
				&buf,
				logger,
			)
			file.written = []byte{}
			terminalUnderTest.cursorHeight = trial.cursorHeight

			// -----------------

			terminalUnderTest.scrollWindow(trial.previousRowCount, trial.currentRowCount)

			assert.Equal(tt, trial.expectedCursorHeight, terminalUnderTest.cursorHeight)
			assert.Equal(tt, []byte(trial.expectedWrite), file.written)
		})
	}
}
