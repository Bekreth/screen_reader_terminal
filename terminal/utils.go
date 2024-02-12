package terminal

import "github.com/bekreth/screen_reader_terminal/utils"

// Takes the previous and current lines and determine what characters need to be
// written at the current cursor's location
func lineDiff(
	lastLineData string, lastCursor int,
	lineData string, cursor int,
) string {
	checkEdge := utils.IntMin(lastCursor, cursor)
	newEnd := lineData[checkEdge:]
	for i := 0; i <= checkEdge; i++ {
		if lineData[0:i] != lastLineData[0:i] {
			newEnd = lineData[i:]
			break
		}
	}
	return newEnd
}
