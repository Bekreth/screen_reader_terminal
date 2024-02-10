package terminal

import (
	"strings"

	"github.com/bekreth/screen_reader_terminal/buffer"
	"github.com/bekreth/screen_reader_terminal/history"
	"github.com/bekreth/screen_reader_terminal/utils"
	"github.com/bekreth/screen_reader_terminal/window"
)

type Terminal struct {
	window  window.Window
	buffer  *buffer.Buffer
	history *history.History
	logger  utils.Logger
}

func NewTerminal(
	win window.Window,
	buf *buffer.Buffer,
	logger utils.Logger,
) Terminal {
	history := history.NewBufferHistory()
	win.ClearWindow(window.FULL)
	win.SetCursorPosition(0, 0)
	return Terminal{
		window:  win,
		buffer:  buf,
		history: &history,
		logger:  logger,
	}
}

func (terminal *Terminal) AddBuffer(buffer *buffer.Buffer) {
	terminal.history.AddBuffer(*terminal.buffer)
	terminal.buffer = buffer
}

func (terminal Terminal) CurrentBuffer() *buffer.Buffer {
	return terminal.buffer
}

func (terminal Terminal) Draw() {
	previousData, previousCursor := terminal.buffer.PreviousOutput()
	splitPreviousData := strings.Split(previousData, "\n")
	_, previousCursorRow, previousCursorOffset := terminal.determineRows(
		previousData,
		previousCursor,
	)

	currentData, currentCursor := terminal.buffer.Output()
	splitCurrentData := strings.Split(currentData, "\n")
	_, currentCursorRow, currentCursorOffset := terminal.determineRows(
		currentData,
		currentCursor,
	)

	zippedLines := utils.Zip(
		splitPreviousData,
		splitCurrentData,
		func() string { return "" },
	)

	didUpdate := false

	if previousData != currentData {
		for i, dataPair := range zippedLines {
			if dataPair.First != dataPair.Second {
				firstCursor := 0
				secondCursor := 0
				restoreCursor := false

				if previousCursorRow != i+1 {
					//TODO:
					terminal.logger.Debugf("Not last line")
					terminal.window.SaveCursor()
					x := currentCursorOffset - previousCursorOffset
					y := currentCursorRow - previousCursorRow
					terminal.window.MoveCursor(x, y)
					restoreCursor = true
				} else {
					firstCursor = previousCursorOffset
					secondCursor = currentCursorOffset
				}
				terminal.drawLine(dataPair.First, firstCursor, dataPair.Second, secondCursor)
				didUpdate = true
				if restoreCursor {
					terminal.window.RestoreCursor()
				}
			}
		}
	} else if previousCursor != currentCursor {
		terminal.moveCursor(previousCursor, currentCursor)
		didUpdate = true
	}

	if didUpdate {
		terminal.buffer.UpdatePrevious()
	}
}

func lineDiff(
	lastLineData string, lastCursor int,
	lineData string, cursor int,
) string {
	checkEdge := utils.IntMin(lastCursor, cursor)
	newEnd := lineData[checkEdge:]
	for i := 0; i <= checkEdge; i++ {
		if lineData[0:i] != lastLineData[0:i] {
			newEnd = lineData[i:]
			break
		}
	}
	return newEnd
}

// drawLine
func (terminal Terminal) drawLine(
	previousLineData string, previousCursor int,
	currentLineData string, currentCursor int,
) {
	width := terminal.window.GetWindowSize().Width
	newEnd := lineDiff(
		previousLineData, previousCursor,
		currentLineData, currentCursor,
	)

	if currentCursor < previousCursor {
		terminal.moveCursor(previousCursor, currentCursor)
		terminal.window.ClearLine(window.CURSOR_FORWARD)
		terminal.window.Write([]byte(newEnd))
		terminal.window.MoveCursor(-1*len(newEnd), 0)
	} else {
		terminal.window.Write([]byte(newEnd))
		if len(newEnd) != currentCursor-previousCursor {
			backshift := len(newEnd) - 1
			delta := backshift % width
			rollback := backshift / width
			terminal.window.MoveCursor(-1*delta, rollback)
		}
	}
}

// moveCursor calculates the difference between cursor and lastCursor and writes the
// appropriate ANSII control characters to make the terminal match the difference
func (terminal Terminal) moveCursor(previousCursor int, currentCursor int) {
	// TODO make this new line aware
	windowWidth := terminal.window.GetWindowSize().Width
	_, rollover := utils.ModAdd(currentCursor, 0, windowWidth)
	_, lastRollover := utils.ModAdd(previousCursor, 0, windowWidth)
	y := rollover - lastRollover
	x := currentCursor - previousCursor + (-1 * y * windowWidth)
	terminal.window.MoveCursor(x, y)
}

func (terminal Terminal) NewLine() {
	currentValue, currentPosition := terminal.CurrentBuffer().Output()
	terminal.moveCursor(
		currentPosition,
		len(currentValue),
	)
	lineCount := terminal.CurrentBuffer().NewLineCount()
	newLines := "\n"
	for i := 0; i < lineCount-1; i++ {
		newLines += "\n"
	}
	terminal.window.Write([]byte(newLines))
	terminal.history.AddBuffer(*terminal.buffer)
	terminal.buffer.Clear()
}

// Calculates how many rows the current value crosses and on which line
// the cursor is currently positioned
func (terminal Terminal) determineRows(currentValue string, cursor int) (int, int, int) {
	if currentValue == "" {
		return 1, 1, 0
	}
	width := terminal.window.GetWindowSize().Width

	newLineIndicies := append([]int{0}, utils.IndiciesOfChar(currentValue, '\n')...)
	splitValues := strings.Split(currentValue, "\n")

	totalRowCount := 0
	cursorRow := 0
	cursorOffset := 0
	for i := range newLineIndicies {
		lineLength := len(splitValues[i]) + 1
		rowIncrementor := int(lineLength/width) + 1
		totalRowCount += rowIncrementor

		if cursor >= 0 {
			if lineLength > cursor {
				cursorOffset = cursor % width
				if cursor == 1 {
					cursorRow += 1
				} else {
					cursorRow += (cursor / width) + 1
				}
			} else {
				cursorRow += rowIncrementor
			}
			cursor -= lineLength
		}
	}
	return totalRowCount, cursorRow, cursorOffset
}
