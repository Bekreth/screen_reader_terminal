package terminal

import (
	"github.com/bekreth/screen_reader_terminal/buffer"
	"github.com/bekreth/screen_reader_terminal/history"
	"github.com/bekreth/screen_reader_terminal/utils"
	"github.com/bekreth/screen_reader_terminal/window"
)

type Terminal struct {
	cursorHeight int
	window       window.Window
	buffer       *buffer.Buffer
	history      *history.History
	logger       utils.Logger
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
		cursorHeight: 0,
		window:       win,
		buffer:       buf,
		history:      &history,
		logger:       logger,
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

func (terminal *Terminal) Draw() {
	// Breaking up data from previous render
	previousData, previousCursor := terminal.buffer.PreviousOutput()
	previousDataRow, previousCursorRow, previousCursorOffset := terminal.determineRows(
		previousData,
		previousCursor,
	)

	// Breaking up data from current render
	currentData, currentCursor := terminal.buffer.Output()
	currentDataRow, currentCursorRow, currentCursorOffset := terminal.determineRows(
		currentData,
		currentCursor,
	)

	// Splicing together
	zippedLines := utils.Zip(
		previousDataRow,
		currentDataRow,
		func() string { return emptyString },
	)

	terminal.scrollWindow(len(previousDataRow), len(currentDataRow))

	// Calculating delta
	//width := terminal.window.GetWindowSize().Width
	coords := newCoords(previousCursorOffset, previousCursorRow)
	if previousData != currentData {
		for i, dataPair := range zippedLines {
			rowRequiresUpdate := dataPair.First != dataPair.Second
			if rowRequiresUpdate {
				coords = coords.setPendingRow(i)
				if dataPair.Second == emptyString {
					moveX, moveY := coords.outputDelataToTarget()
					terminal.window.MoveCursor(moveX, moveY)
					coords = coords.applyPendingDeltas()

					terminal.window.ClearLine(window.FULL)
					continue
				}
				coords = terminal.drawRow(dataPair.First, dataPair.Second, coords)
			}
		}
	}
	coords = coords.setPendingColumn(currentCursorOffset).
		setPendingRow(currentCursorRow)

	moveX, moveY := coords.outputDelataToTarget()
	coords = coords.applyPendingDeltas()

	terminal.window.MoveCursor(moveX, moveY)
	terminal.buffer.UpdatePrevious()
}

func (terminal *Terminal) NewLine() {
	terminal.cursorHeight += 1
	terminal.window.Write([]byte("\n"))
	terminal.history.AddBuffer(*terminal.buffer)
	terminal.buffer.Clear()
}
