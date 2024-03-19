package terminal

import "github.com/bekreth/screen_reader_terminal/utils"

// moveCursor calculates the difference between cursor and lastCursor and writes the
// appropriate ANSII control characters to make the terminal match the difference
func (terminal Terminal) moveCursor(
	previousCursorRow int, previousCursorOffset int,
	currentCursorRow int, currentCursorOffset int,
) {
	terminal.logger.Debugf(
		"move cursor: PR: %v:%v, CR: %v:%v",
		previousCursorRow, previousCursorOffset,
		currentCursorRow, currentCursorOffset,
	)

	width := terminal.window.GetWindowSize().Width
	_, previousRollover := utils.ModAdd(previousCursorOffset, 0, width)
	_, currentRollover := utils.ModAdd(currentCursorOffset, 0, width)

	rolloverValue := currentRollover - previousRollover
	rowDiff := currentCursorRow - previousCursorRow

	y := rolloverValue + rowDiff
	x := currentCursorOffset - previousCursorOffset + (-1 * rolloverValue * width)

	terminal.window.MoveCursor(x, y)
}
