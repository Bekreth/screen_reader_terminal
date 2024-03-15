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
	terminal.logger.Debugf("CURSOR %v\tVALUE '%v'", cursor, currentValue)
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
		terminal.logger.Debugf("  LINE '%v', LEN '%v'", line, len(line))
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
			row := line[j:lineEnd]
			terminal.logger.Debugf("    ROW '%v', LINEEND '%v'", row, lineEnd)
			rows = append(rows, row)
			if advancingCursor {
				advanceCursorStep := utils.IntMin(len(row), cursor-cursorCounter)
				cursorCounter += advanceCursorStep
				cursorOffset = advanceCursorStep % width
				terminal.logger.Debugf("      ROWLEN '%v', WIDTH '%v'", len(row), width)
				terminal.logger.Debugf("      AC '%v', CC '%v', CO '%v'", advanceCursorStep, cursorCounter, cursorOffset)
				if cursorCounter >= cursor {
					if cursorCounter%width == 0 {
						cursorRows += 1
					}
					terminal.logger.Debugf("    === Closing cursor ===")
					advancingCursor = false
				}
			}
			if advancingCursor {
				cursorRows += 1
			}
			terminal.logger.Debugf("      CROW: %v, OFFSET: %v\n", cursorRows, cursorOffset)
		}
	}
	return rows, cursorRows, cursorOffset
}
