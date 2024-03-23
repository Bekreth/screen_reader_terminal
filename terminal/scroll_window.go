package terminal

func (terminal *Terminal) scrollWindow(previousRowCount int, currentRowCount int) {
	if previousRowCount > 0 {
		previousRowCount -= 1
	}
	if currentRowCount > 0 {
		currentRowCount -= 1
	}

	rowCountDelta := currentRowCount - previousRowCount

	projectedMaxRow := rowCountDelta + terminal.cursorHeight
	height := terminal.window.GetWindowSize().Height
	if projectedMaxRow >= height {
		fitDelta := projectedMaxRow - height + 1
		terminal.window.ScrollPage(fitDelta)
		terminal.window.MoveCursor(0, -1*fitDelta)
		terminal.cursorHeight -= fitDelta
	}

	terminal.cursorHeight += rowCountDelta
}
