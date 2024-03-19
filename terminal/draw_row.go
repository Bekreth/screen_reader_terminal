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
	var newEnd string
	var column int
	shouldClearFromCursor := false
	if previousRowData == "" {
		newEnd = currentRowData
	} else {
		newEnd, column = rowDiff(previousRowData, currentRowData)
		coords = coords.setPendingColumn(column)
		shouldClearFromCursor = len(previousRowData) > len(currentRowData) &&
			previousRowData != emptyString
	}

	xMove, yMove := coords.outputDelataToTarget()
	terminal.window.MoveCursor(xMove, yMove)
	coords = coords.applyPendingDeltas()

	if shouldClearFromCursor {
		terminal.window.ClearLine(window.CURSOR_FORWARD)
	}

	terminal.window.Write([]byte(newEnd))

	coords = coords.addColumnDelta(len(newEnd))
	coords = coords.applyPendingDeltas()
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
