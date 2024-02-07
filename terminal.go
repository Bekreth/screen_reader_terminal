package main

type Terminal struct {
	window  Window
	buffer  *Buffer
	history *History
	logger  Logger
}

func NewTerminal(window Window, buffer *Buffer, logger Logger) Terminal {
	history := NewBufferHistory()
	return Terminal{
		window:  window,
		buffer:  buffer,
		history: &history,
		logger:  logger,
	}
}

func (terminal Terminal) Draw() {
	lastData, lastCursor := terminal.buffer.PreviousOutput()
	data, cursor := terminal.buffer.Output()
	windowWidth := terminal.window.GetWindowSize().width

	didUpdate := false

	if lastData != data {
		checkEdge := IntFloor(lastCursor, cursor)
		newEnd := data[checkEdge:]
		for i := 0; i <= checkEdge; i++ {
			if data[0:i] != lastData[0:i] {
				newEnd = data[i:]
				break
			}
		}

		if cursor < lastCursor {
			terminal.window.SaveCursor()

			terminal.moveCursor(cursor, lastCursor)
			terminal.window.ClearLine(CURSOR_FORWARD)

			currentLine := cursor / windowWidth
			totalLines := len(data) / windowWidth
			if totalLines != currentLine {
				terminal.window.MoveCursor(0, totalLines-currentLine)
				terminal.window.ClearLine(FULL_LINE)
				terminal.window.MoveCursor(0, currentLine-totalLines)
			}

			terminal.window.Write([]byte(newEnd))
			terminal.window.RestoreCursor()
		} else {
			terminal.window.SaveCursor()
			terminal.window.Write([]byte(newEnd))
			terminal.window.RestoreCursor()
		}
		didUpdate = true
	}

	if lastCursor != cursor {
		terminal.moveCursor(cursor, lastCursor)
		didUpdate = true
	}

	if didUpdate {
		terminal.buffer.UpdatePrevious()
	}
}

func (terminal Terminal) moveCursor(cursor int, lastCursor int) {
	windowWidth := terminal.window.GetWindowSize().width
	_, rollover := ModAdd(cursor, 0, windowWidth)
	_, lastRollover := ModAdd(lastCursor, 0, windowWidth)
	y := rollover - lastRollover
	x := cursor - lastCursor + (-1 * y * windowWidth)
	terminal.window.MoveCursor(x, y)
}

func (terminal Terminal) NewLine() {
	terminal.moveCursor(len(terminal.buffer.currentValue), terminal.buffer.cursorPosition)
	terminal.window.Write([]byte{'\n'})
	terminal.history.AddBuffer(*terminal.buffer)
	terminal.buffer.Clear()
}
