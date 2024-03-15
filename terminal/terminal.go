package terminal

import (
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

	terminal.logger.Debugf(
		"DETERMINES: previous %v: '%v' current %v: '%v'",
		len(previousDataRow),
		previousDataRow,
		len(currentDataRow),
		currentDataRow,
	)
	// Splicing together
	zippedLines := utils.Zip(
		previousDataRow,
		currentDataRow,
		func() string { return emptyString },
	)

	// Calculating delta
	//width := terminal.window.GetWindowSize().Width
	terminal.logger.Debugf(
		"PC '%v' PR '%v', CC '%v' CR '%v', ",
		previousCursorOffset, previousCursorRow,
		currentCursorOffset, currentCursorRow,
	)
	coords := newCoords(previousCursorOffset, previousCursorRow)
	if previousData != currentData {
		terminal.logger.Debugf(
			"PR %v PC %v CR %v CC %v",
			previousCursorRow, previousCursorOffset,
			currentCursorRow, currentCursorOffset,
		)
		for i, dataPair := range zippedLines {
			terminal.logger.Debugf("================= ROW # %v", 1)
			terminal.logger.Debugf("PD '%v' CD '%v'", dataPair.First, dataPair.Second)
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
				terminal.logger.Debugf("LOCATION Before: %v", coords)
				coords = terminal.drawRow(dataPair.First, dataPair.Second, coords)
				terminal.logger.Debugf("LOCATION after: %v", coords)
			}
		}
	}
	coords = coords.setPendingColumn(currentCursorOffset).
		setPendingRow(currentCursorRow)

	terminal.logger.Debugf("TARGET LOC: %v", coords)
	moveX, moveY := coords.outputDelataToTarget()
	coords = coords.applyPendingDeltas()

	terminal.window.MoveCursor(moveX, moveY)
	terminal.buffer.UpdatePrevious()
	terminal.logger.Debugf("")
	terminal.logger.Debugf("")
}

func (terminal Terminal) NewLine() {
	terminal.window.Write([]byte("\n"))
	terminal.history.AddBuffer(*terminal.buffer)
	terminal.buffer.Clear()
}
