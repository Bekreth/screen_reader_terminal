package terminal

import "github.com/bekreth/screen_reader_terminal/utils"

// Takes the previous and current lines and determine what characters need to be
// written at the current cursor's location
func lineDiff(
	previousLineData string, previousCursor int,
	currentLineData string, currentCursor int,
) string {
	checkEdge := utils.IntMin(previousCursor, currentCursor)
	newEnd := currentLineData[checkEdge:]
	for i := 0; i < checkEdge; i++ {
		if currentLineData[0:i] != previousLineData[0:i] {
			newEnd = currentLineData[0:i]
			break
		}
	}
	return newEnd
}
