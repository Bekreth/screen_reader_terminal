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
	trialsCollections := [][]determineRowTrial{
		basicTests,
		serveralSmallLines,
		severalVeryLongLines,
	}

	trials := []determineRowTrial{}
	for _, trialCollection := range trialsCollections {
		trials = append(trials, trialCollection...)
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
				Prefix:   "",
				Value:    "",
				Position: 0,
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

			assert.Equal(tt, trial.expectedRows, actualRow, "ROWS")
			assert.Equal(tt, trial.expectedCursorRow, actualCursor, "CUROSOR ROWS")
			assert.Equal(tt, trial.expectedOffset, actualOffset, "CUROSOR OFFSET")
		})
	}
}

type determineRowTrial struct {
	description string

	prefix          string
	currentValue    string
	currentPosition int

	expectedRows      []string
	expectedCursorRow int
	expectedOffset    int
}

var basicTests = []determineRowTrial{
	{
		description:       "Empty buffer, should see 1,1",
		prefix:            "",
		currentPosition:   0,
		currentValue:      "",
		expectedRows:      []string{},
		expectedCursorRow: 0,
		expectedOffset:    0,
	},
	{
		description:       "no buffer but short prefix",
		prefix:            "short",
		currentPosition:   0,
		currentValue:      "",
		expectedRows:      []string{"short"},
		expectedCursorRow: 0,
		expectedOffset:    5,
	},
	{
		description:       "short string with no prefix",
		prefix:            "",
		currentPosition:   12,
		currentValue:      "small string",
		expectedRows:      []string{"small string"},
		expectedCursorRow: 0,
		expectedOffset:    12,
	},
	{
		description:       "cursor in front of long prefix/string",
		prefix:            "small: ",
		currentPosition:   0,
		currentValue:      "012345678901234567890",
		expectedRows:      []string{"small: 0123456789012", "34567890"},
		expectedCursorRow: 0,
		expectedOffset:    7,
	},
	{
		description:     "cursor at end of long prefix/string",
		prefix:          "small: ",
		currentPosition: 18,
		currentValue:    "012345678901234567",
		expectedRows: []string{
			"small: 0123456789012",
			"34567",
		},
		expectedCursorRow: 1,
		expectedOffset:    5,
	},
	{
		description:     "cursor in the middle of long prefix/string",
		prefix:          "small: ",
		currentPosition: 18,
		currentValue:    "01234567890123456789012345678901234567890",
		expectedRows: []string{
			"small: 0123456789012",
			"34567890123456789012",
			"34567890",
		},
		expectedCursorRow: 1,
		expectedOffset:    5,
	},
	{
		description:     "cursor end of very long line",
		prefix:          "small: ",
		currentPosition: 41,
		currentValue:    "01234567890123456789012345678901234567890",
		expectedRows: []string{
			"small: 0123456789012",
			"34567890123456789012",
			"34567890",
		},
		expectedCursorRow: 2,
		expectedOffset:    8,
	},
	{
		description:       "cursor end of short line with new line character",
		prefix:            "small: ",
		currentPosition:   12,
		currentValue:      "hello world\n",
		expectedRows:      []string{"small: hello world", ""},
		expectedCursorRow: 1,
		expectedOffset:    0,
	},
	/*
		{
			description:       "Rollover match at max",
			prefix:            "",
			currentPosition:   20,
			currentValue:      "12345678901234567890",
			expectedRows:      []string{"12345678901234567890"},
			expectedCursorRow: 1,
			expectedOffset:    0,
		},
	*/
}

var serveralSmallLines = []determineRowTrial{
	{
		description:       "0. several small lines, cursor at start",
		prefix:            "small: ",
		currentPosition:   0,
		currentValue:      "hello\nworld\ncheck\nme",
		expectedRows:      []string{"small: hello", "world", "check", "me"},
		expectedCursorRow: 0,
		expectedOffset:    7,
	},
	{
		description:       "1. several small lines, cursor at second word",
		prefix:            "small: ",
		currentPosition:   6,
		currentValue:      "hello\nworld\ncheck\nme",
		expectedRows:      []string{"small: hello", "world", "check", "me"},
		expectedCursorRow: 1,
		expectedOffset:    0,
	},
	{
		description:       "2. several small lines, cursor at third word",
		prefix:            "small: ",
		currentPosition:   12,
		currentValue:      "hello\nworld\ncheck\nme",
		expectedRows:      []string{"small: hello", "world", "check", "me"},
		expectedCursorRow: 2,
		expectedOffset:    0,
	},
	{
		description:       "3. several small lines, cursor at forth word",
		prefix:            "small: ",
		currentPosition:   18,
		currentValue:      "hello\nworld\ncheck\nme",
		expectedRows:      []string{"small: hello", "world", "check", "me"},
		expectedCursorRow: 3,
		expectedOffset:    0,
	},
}

var veryLongLine = fmt.Sprintf(
	"18. %v\n25. %v\n35. %v\n45. %v\n55. %v",
	"456789012345678",
	"4567890123456789012345",
	"45678901234567890123456789012345",
	"456789012345678901234567890123456789012345",
	"45678901234567890123456789012345678901234567890",
)
var expectedLongLine = []string{
	"small: 18. 456789012",
	"345678",
	"25. 4567890123456789",
	"012345",
	"35. 4567890123456789",
	"0123456789012345",
	"45. 4567890123456789",
	"01234567890123456789",
	"012345",
	"55. 4567890123456789",
	"01234567890123456789",
	"01234567890",
}

var severalVeryLongLines = []determineRowTrial{
	{
		description:       "A. Several very long lines, first word",
		prefix:            "small: ",
		currentPosition:   0,
		currentValue:      veryLongLine,
		expectedRows:      expectedLongLine,
		expectedCursorRow: 0,
		expectedOffset:    7,
	},
	{
		description:       "B. Several very long lines, end of first word",
		prefix:            "small: ",
		currentPosition:   19,
		currentValue:      veryLongLine,
		expectedRows:      expectedLongLine,
		expectedCursorRow: 1,
		expectedOffset:    6,
	},
	{
		description:       "C. Several very long lines, start of second word",
		prefix:            "small: ",
		currentPosition:   20,
		currentValue:      veryLongLine,
		expectedRows:      expectedLongLine,
		expectedCursorRow: 2,
		expectedOffset:    0,
	},
	{
		description:       "D. Several very long lines, end of second word",
		prefix:            "small: ",
		currentPosition:   46,
		currentValue:      veryLongLine,
		expectedRows:      expectedLongLine,
		expectedCursorRow: 3,
		expectedOffset:    6,
	},
	{
		description:       "E. Several very long lines, start of third word",
		prefix:            "small: ",
		currentPosition:   47,
		currentValue:      veryLongLine,
		expectedRows:      expectedLongLine,
		expectedCursorRow: 4,
		expectedOffset:    0,
	},
	{
		description:       "G. Several very long lines, end of third word",
		prefix:            "small: ",
		currentPosition:   83,
		currentValue:      veryLongLine,
		expectedRows:      expectedLongLine,
		expectedCursorRow: 5,
		expectedOffset:    16,
	},
	{
		description:       "H. Several very long lines, Start of fourth word",
		prefix:            "small: ",
		currentPosition:   84,
		currentValue:      veryLongLine,
		expectedRows:      expectedLongLine,
		expectedCursorRow: 6,
		expectedOffset:    0,
	},
	{
		description:       "I. Several very long lines, end of fifth word",
		prefix:            "small: ",
		currentPosition:   181,
		currentValue:      veryLongLine,
		expectedRows:      expectedLongLine,
		expectedCursorRow: 11,
		expectedOffset:    10,
	},
}
