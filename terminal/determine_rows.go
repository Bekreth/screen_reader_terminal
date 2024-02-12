package terminal

import (
	"strings"

	"github.com/bekreth/screen_reader_terminal/utils"
)

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
