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

const emptyString = "[~empty~]"

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
		func() string { return emptyString },
	)

	didUpdate := false

	if previousData != currentData {
		for i, dataPair := range zippedLines {
			if dataPair.First != dataPair.Second {
				thisRow := i + 1
				firstCursor := 0
				secondCursor := 0

				y := currentCursorRow - previousCursorRow
				if previousCursorRow != thisRow {
					// Rewrite Whole line
					existingLine := dataPair.First != emptyString
					if existingLine {
						terminal.window.SaveCursor()
					}
					terminal.window.MoveCursor(0, y)
					terminal.window.SetCursorColumn(0)
					terminal.window.ClearLine(window.FULL)
					terminal.window.Write([]byte(dataPair.Second))
					if existingLine {
						terminal.window.RestoreCursor()
					}
					didUpdate = true
				} else {
					// Update Line
					terminal.logger.Debugf("update line")
					firstCursor = previousCursorOffset
					secondCursor = currentCursorOffset
					terminal.drawLine(dataPair.First, firstCursor, dataPair.Second, secondCursor)
					didUpdate = true
				}
			}
		}
	} else if previousCursor != currentCursor {
		terminal.moveCursor(
			previousCursorRow, previousCursorOffset,
			currentCursorRow, currentCursorOffset,
		)
		didUpdate = true
	}

	if didUpdate {
		terminal.buffer.UpdatePrevious()
	}
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
		terminal.moveCursor(0, previousCursor, 0, currentCursor)
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
func (terminal Terminal) moveCursor(
	previousCursorRow int, previousCursorOffset int,
	currentCursorRow int, currentCursorOffset int,
) {
	windowWidth := terminal.window.GetWindowSize().Width
	_, currentRollover := utils.ModAdd(currentCursorOffset, 0, windowWidth)
	_, previousRollover := utils.ModAdd(previousCursorOffset, 0, windowWidth)
	rolloverValue := currentRollover - previousRollover

	y := rolloverValue + (previousCursorRow - currentCursorRow)
	x := currentCursorOffset - previousCursorOffset + (-1 * rolloverValue * windowWidth)

	terminal.window.MoveCursor(x, y)
}

func (terminal Terminal) NewLine() {
	currentValue, currentPosition := terminal.CurrentBuffer().Output()
	terminal.moveCursor(
		0, currentPosition,
		0, len(currentValue),
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
