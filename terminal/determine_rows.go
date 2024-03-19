package terminal

import (
	"strings"

	"github.com/bekreth/screen_reader_terminal/utils"
)

// Calculates how many rows the current value crosses and on which line
// the cursor is currently positioned
func (terminal Terminal) determineRows(
	currentValue string,
	cursor int,
) ([]string, int, int) {
	if currentValue == "" {
		return []string{}, 0, 0
	}
	width := terminal.window.GetWindowSize().Width
	splitValues := strings.Split(currentValue, "\n")

	rows := []string{}
	cursorRows := 0
	cursorOffset := 0

	cursorCounter := 0
	advancingCursor := true
	for i, line := range splitValues {
		lineLength := len(line)
		if lineLength == 0 {
			rows = append(rows, "")
		}
		if advancingCursor {
			cursorOffset = 0
			if lineLength == 0 {
				continue
			}
			if i != 0 {
				cursorCounter += 1
			}
		}

		for j := 0; j < lineLength; j += width {
			lineEnd := utils.IntMin(j+width, lineLength)
			row := strings.ReplaceAll(line[j:lineEnd], "\t", " ")
			rows = append(rows, row)
			if advancingCursor {
				advanceCursorStep := utils.IntMin(len(row), cursor-cursorCounter)
				cursorCounter += advanceCursorStep
				cursorOffset = advanceCursorStep % width
				if cursorCounter >= cursor {
					if cursorCounter%width == 0 {
						cursorRows += 1
					}
					advancingCursor = false
				}
			}
			if advancingCursor {
				cursorRows += 1
			}
		}
	}
	return rows, cursorRows, cursorOffset
}
