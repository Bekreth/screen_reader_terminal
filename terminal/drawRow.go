package terminal

import (
	"github.com/bekreth/screen_reader_terminal/utils"
	"github.com/bekreth/screen_reader_terminal/window"
)

func (terminal Terminal) drawRow(
	previousRowData string,
	currentRowData string,
	coords coordinates,
) coordinates {
	terminal.logger.Debugf(
		"  DRAW ROW: \nPRD %v\nCRD %v",
		previousRowData,
		currentRowData,
	)

	var newEnd string
	var column int
	shouldClearFromCursor := false
	if previousRowData == "" {
		newEnd = currentRowData
	} else {
		newEnd, column = rowDiff(previousRowData, currentRowData)
		coords = coords.setPendingColumn(column)
		terminal.logger.Debugf("  COORDS 1: %v", coords)
		shouldClearFromCursor = len(previousRowData) > len(currentRowData) &&
			previousRowData != emptyString
	}

	xMove, yMove := coords.outputDelataToTarget()
	terminal.logger.Debugf("  MOVE 1: %v %v", xMove, yMove)
	terminal.window.MoveCursor(xMove, yMove)
	coords = coords.applyPendingDeltas()

	if shouldClearFromCursor {
		terminal.window.ClearLine(window.CURSOR_FORWARD)
	}

	terminal.logger.Debugf("  NEW END: %v", newEnd)
	terminal.window.Write([]byte(newEnd))

	coords = coords.addColumnDelta(len(newEnd))
	coords = coords.applyPendingDeltas()
	terminal.logger.Debugf("  COORDS 2: %v", coords)
	return coords
}

func rowDiff(previousRowData string, currentRowData string) (string, int) {
	checkEdge := utils.IntMin(len(previousRowData), len(currentRowData))
	newEnd := currentRowData[checkEdge:]
	column := 0
	for i := 0; i < checkEdge; i++ {
		if column == 0 {
			column = checkEdge
		}
		if currentRowData[i] != previousRowData[i] {
			newEnd = currentRowData[i:]
			column = i
			break
		}
	}
	return newEnd, column
}

// drawRow draws the row using known cursor position and returns how much the cursor
// moved while writing
/*
func (terminal Terminal) drawRow(
	shouldBackshift bool,
	previousLineData string, previousCursor int,
	currentLineData string, currentCursor int,
) int {
	terminal.logger.Debugf(
		"DRAW ROW: \nPC %v - PD '%v'\n CC %v - CD '%v'",
		previousCursor, previousLineData,
		currentCursor, currentLineData,
	)

	width := terminal.window.GetWindowSize().Width
	var newEnd string
	if previousLineData == "" {
		newEnd = currentLineData
	} else {
		adjustedCursor := currentCursor
		if len(currentLineData)-len(previousLineData) == 1 &&
			previousCursor == width-1 &&
			currentCursor == 0 {
			adjustedCursor = width
		}
		newEnd = lineDiff(
			previousLineData, previousCursor,
			currentLineData, adjustedCursor,
		)
	}

	cursorDelta := currentCursor - previousCursor
	edgeRollback := cursorDelta == -1*width+1
	if currentLineData == "" && previousLineData == emptyString {
		terminal.window.MoveCursor(0, 1)
		terminal.window.SetCursorColumn(0)
	} else if currentCursor < previousCursor && !edgeRollback {
		// Update line in place
		terminal.logger.Debugf("Update line in place")
		terminal.moveCursor(0, previousCursor, 0, currentCursor)
		terminal.window.ClearLine(window.CURSOR_FORWARD)
		terminal.window.Write([]byte(newEnd))
		terminal.window.MoveCursor(-1*len(newEnd), 0)
	} else {
		terminal.window.Write([]byte(newEnd))
		if shouldBackshift && len(newEnd) != cursorDelta {
			backshift := len(newEnd) - 1
			delta := backshift % width
			rollback := backshift / width
			terminal.window.MoveCursor(-1*delta, rollback)
		}
	}
}
*/
