package terminal

import (
	"testing"

	"github.com/bekreth/screen_reader_terminal/buffer"
	"github.com/bekreth/screen_reader_terminal/history"
	"github.com/bekreth/screen_reader_terminal/utils"
	"github.com/bekreth/screen_reader_terminal/window"
	"github.com/stretchr/testify/assert"
)

// tests

type drawTrial struct {
	description      string
	previousPrefix   string
	prefix           string
	previousValue    string
	previousPosition int
	currentValue     string
	currentPosition  int
	expectedOutput   string
}

func TestDraw(t *testing.T) {
	trialsCollections := [][]drawTrial{
		singleLineManipulations,
		updatesWithNewLine,
		longLineInsertRollover,
		longLineDeleteRollover,
		moveAcrossWrap,
		removeFromMultiline,
		insertInMultiline,
	}

	trials := []drawTrial{}
	for _, trialCollection := range trialsCollections {
		trials = append(trials, trialCollection...)
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

var singleLineManipulations = []drawTrial{
	{
		description:      "Draw empty buffer",
		previousValue:    "",
		previousPosition: 0,
		currentValue:     "",
		currentPosition:  0,
		expectedOutput:   "",
	},
	{
		description:      "Draw new buffer",
		previousValue:    "",
		previousPosition: 0,
		currentValue:     "Hello world",
		currentPosition:  11,
		expectedOutput: fmtLine(
			"Hello world",
		),
	},
	{
		description:      "Add character to buffer end",
		previousValue:    "Hello worl",
		previousPosition: 10,
		currentValue:     "Hello world",
		currentPosition:  11,
		expectedOutput: fmtLine(
			"d",
		),
	},
	{
		description:      "Remove character from buffer end",
		previousValue:    "Hello world",
		previousPosition: 11,
		currentValue:     "Hello worl",
		currentPosition:  10,
		expectedOutput: fmtLine(
			left(1),
			clearCursorForward(),
		),
	},
	{
		description:      "Insert text in middle of line",
		previousValue:    "Hello world",
		previousPosition: 5,
		currentValue:     "Hello, world",
		currentPosition:  6,
		expectedOutput: fmtLine(
			", world",
			left(6),
		),
	},
	{
		description:      "Move cursor back 3",
		previousValue:    "Hello world",
		previousPosition: 10,
		currentValue:     "Hello world",
		currentPosition:  7,
		expectedOutput: fmtLine(
			left(3),
		),
	},
}

var updatesWithNewLine = []drawTrial{
	{
		description:      "Move cursor forward over new line",
		previousValue:    "Hello\nworld",
		previousPosition: 4,
		currentValue:     "Hello\nworld",
		currentPosition:  11,
		expectedOutput: fmtLine(
			right(1),
			down(1),
		),
	},
	{
		description:      "Move cursor back over new line",
		previousValue:    "Hello\nworld",
		previousPosition: 11,
		currentValue:     "Hello\nworld",
		currentPosition:  4,
		expectedOutput: fmtLine(
			left(1),
			up(1),
		),
	},
	{
		description:      "Update short line to have new line",
		previousValue:    "Hello",
		previousPosition: 5,
		currentValue:     "Hello\n",
		currentPosition:  6,
		expectedOutput: fmtLine(
			left(5),
			down(1),
		),
	},
	{
		description:      "Remove new line from short line",
		previousValue:    "Hello\n",
		previousPosition: 6,
		currentValue:     "Hello",
		currentPosition:  5,
		expectedOutput: fmtLine(
			clearCursorFullLine(),
			right(5),
			up(1),
		),
	},
	{
		description:      "Update short line to have new line and short line",
		previousValue:    "Hello",
		previousPosition: 5,
		currentValue:     "Hello\nWorld",
		currentPosition:  11,
		expectedOutput: fmtLine(
			left(5),
			down(1),
			"World",
		),
	},
	{
		description:      "Update last line in multi-line statement",
		previousValue:    "Hello\nworl",
		previousPosition: 10,
		currentValue:     "Hello\nworld",
		currentPosition:  11,
		expectedOutput: fmtLine(
			"d",
		),
	},
	{
		description:      "Update first line in multi-line statement",
		previousValue:    "Hello\nworld",
		previousPosition: 5,
		currentValue:     "Hello,\nworld",
		currentPosition:  6,
		expectedOutput: fmtLine(
			",",
		),
	},
	{
		description:      "Paste multiple lines",
		previousValue:    "",
		previousPosition: 0,
		currentValue:     "Hello,\nworld\nmultiple",
		currentPosition:  21,
		expectedOutput: fmtLine(
			"Hello,",
			left(6),
			down(1),
			"world",
			left(5),
			down(1),
			"multiple",
		),
	},
}

var longLineInsertRollover = []drawTrial{
	{
		description:      "Long line well before rollover",
		previousValue:    "123456789012345678",
		previousPosition: 18,
		currentValue:     "1234567890123456789",
		currentPosition:  19,
		expectedOutput: fmtLine(
			"9",
		),
	},
	{
		description:      "Long line before rollover",
		previousValue:    "1234567890123456789",
		previousPosition: 19,
		currentValue:     "12345678901234567890",
		currentPosition:  20,
		expectedOutput: fmtLine(
			"0",
			left(20),
			down(1),
		),
	},
	{
		description:      "Long line during rollover",
		previousValue:    "12345678901234567890",
		previousPosition: 20,
		currentValue:     "123456789012345678901",
		currentPosition:  21,
		expectedOutput: fmtLine(
			"1",
		),
	},
	{
		description:      "Long line after rollover",
		previousValue:    "123456789012345678901",
		previousPosition: 21,
		currentValue:     "1234567890123456789012",
		currentPosition:  22,
		expectedOutput: fmtLine(
			"2",
		),
	},
}

var longLineDeleteRollover = []drawTrial{
	{
		description:      "Long line delete before rollover",
		previousValue:    "1234567890123456789012",
		previousPosition: 22,
		currentValue:     "123456789012345678901",
		currentPosition:  21,
		expectedOutput: fmtLine(
			left(1),
			clearCursorForward(),
		),
	},
	{
		description:      "Long line delete on rollover",
		previousValue:    "123456789012345678901",
		previousPosition: 21,
		currentValue:     "12345678901234567890",
		currentPosition:  20,
		expectedOutput: fmtLine(
			left(1),
			clearCursorFullLine(),
		),
	},
	{
		description:      "Long line delete across rollover",
		previousValue:    "12345678901234567890",
		previousPosition: 20,
		currentValue:     "1234567890123456789",
		currentPosition:  19,
		expectedOutput: fmtLine(
			right(19),
			up(1),
			clearCursorForward(),
		),
	},
	{
		description: "Delete charater from row before row crossover",
		//                                    *                  |*
		//                 123456789012345678901234567890123456789012345678901234567890
		previousValue:    "a simple but long line of text that will allow for editing",
		previousPosition: 39,
		currentValue:     "a simple but long line of text that wil allow for editing",
		currentPosition:  38,
		expectedOutput: fmtLine(
			" ",

			left(20),
			down(1),
			clearCursorForward(),
			"allow for editing",
			right(1),
			up(1),
		),
	},
	{
		description: "Delete charater from row one of multi-line",
		//                                    *                   *
		//                 123456789012345678901234567890123456789012345678901234567890
		previousValue:    "a simple but long line of text that will allow for editing",
		previousPosition: 30,
		currentValue:     "a simple but long line of tex that will allow for editing",
		currentPosition:  29,
		expectedOutput: fmtLine(
			left(1),
			" that will ",

			left(20),
			down(1),
			clearCursorForward(),
			"allow for editing",
			left(8),
			up(1),
		),
	},
}

var moveAcrossWrap = []drawTrial{
	{
		description:      "Move backwards before wrapped line",
		currentValue:     "1234567890123456789012",
		currentPosition:  21,
		previousValue:    "1234567890123456789012",
		previousPosition: 22,
		expectedOutput: fmtLine(
			left(1),
		),
	},
	{
		description:      "Move backwards over wrapped line",
		currentValue:     "1234567890123456789012",
		currentPosition:  20,
		previousValue:    "1234567890123456789012",
		previousPosition: 21,
		expectedOutput: fmtLine(
			left(1),
		),
	},
	{
		description:      "Move backwards after wrapped line",
		currentValue:     "1234567890123456789012",
		currentPosition:  19,
		previousValue:    "1234567890123456789012",
		previousPosition: 20,
		expectedOutput: fmtLine(
			right(19),
			up(1),
		),
	},
}

var insertInMultiline = []drawTrial{
	{
		description: "Insert in last line in multi line phrase",
		//                                    *                   *
		//                 123456789012345678901234567890123456789012345678901234567890
		previousValue:    "a simple but long line of text that will allow for diting",
		previousPosition: 51,
		currentValue:     "a simple but long line of text that will allow for editing",
		currentPosition:  52,
		expectedOutput: fmtLine(
			"editing",
			left(6),
		),
	},
	{
		description: "Insert in middle line in multi line phrase",
		//                                    *                   *
		//                 123456789012345678901234567890123456789012345678901234567890
		previousValue:    "a simple but long line o text that will allow for editing",
		previousPosition: 24,
		currentValue:     "a simple but long line of text that will allow for editing",
		currentPosition:  25,
		expectedOutput: fmtLine(
			"f text that will",

			left(20),
			down(1),
			" allow for editing",

			left(13),
			up(1),
		),
	},
	{
		description: "Insert in first line in multi line phrase",
		//                                    *                   *
		//                 123456789012345678901234567890123456789012345678901234567890
		previousValue:    "a simple ut long line of text that will allow for editing",
		previousPosition: 9,
		currentValue:     "a simple but long line of text that will allow for editing",
		currentPosition:  10,
		expectedOutput: fmtLine(
			"but long li",

			left(20),
			down(1),
			"ne of text that will",

			left(20),
			down(1),
			" allow for editing",

			left(8),
			up(2),
		),
	},
}

var removeFromMultiline = []drawTrial{
	{
		description: "Remove character from end of multiple line wrap",
		//                                    *                   *
		//                 123456789012345678901234567890123456789012345678901234567890
		previousValue:    "a simple but long line of text that will allow for editing",
		previousPosition: 48,
		currentValue:     "a simple but long line of text that will allow or editing",
		currentPosition:  47,
		expectedOutput: fmtLine(
			left(1),
			clearCursorForward(),
			"or editing",
			left(10),
		),
	},
	{
		description: "Remove character from middle of multiple line wrap",
		//                                    *                   *
		//                 123456789012345678901234567890123456789012345678901234567890
		previousValue:    "a simple but long line of text that will allow for editing",
		previousPosition: 27,
		currentValue:     "a simple but long line of ext that will allow for editing",
		currentPosition:  26,
		expectedOutput: fmtLine(
			left(1),
			"ext that will ",

			left(20),
			down(1),
			clearCursorForward(),
			"allow for editing",
			left(11),
			up(1),
		),
	},
	{
		description: "Remove character from first word of multiple line wrap",
		//                                    *                   *
		//                 123456789012345678901234567890123456789012345678901234567890
		previousValue:    "a simple but long line of text that will allow for editing",
		previousPosition: 12,
		currentValue:     "a simple bu long line of text that will allow for editing",
		currentPosition:  11,
		expectedOutput: fmtLine(
			left(1),
			" long lin",

			left(20),
			down(1),
			"e of text that will ",

			left(20),
			down(1),
			clearCursorForward(),
			"allow for editing",

			left(6),
			up(2),
		),
	},
}
