package terminal

import (
	"strings"

	"github.com/bekreth/screen_reader_terminal/utils"
)

// Calculates how many rows the current value crosses and on which line
// the cursor is currently positioned
func (terminal Terminal) determineRows(currentValue string, cursor int) ([]string, int, int) {
	if currentValue == "" {
		return []string{}, 1, 0
	}
	width := terminal.window.GetWindowSize().Width

	newLineIndicies := append([]int{0}, utils.IndiciesOfChar(currentValue, '\n')...)
	splitValues := strings.Split(currentValue, "\n")

	rows := []string{}
	cursorRow := 0
	cursorOffset := 0
	for i := range newLineIndicies {
		lineLength := len(splitValues[i]) + 1
		rowIncrementor := int(lineLength/width) + 1
		previousIndex := 0
		rowLength := len(splitValues[i])
		for j := 0; j < rowIncrementor; j++ {
			row := splitValues[i]
			start := previousIndex
			end := start + width
			if start+width > rowLength {
				end = rowLength
			}
			rows = append(rows, row[start:end])
			previousIndex = end
		}

		if cursor >= 0 {
			if lineLength > cursor {
				cursorOffset = cursor % width
				if cursorOffset == 0 && cursor != 0 {
					cursorOffset = width
				}
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
	return rows, cursorRow, cursorOffset
}
