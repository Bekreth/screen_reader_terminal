package main

import (
	"strings"
)

type Logger interface {
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

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
	lastData, lastCursor := terminal.history.GetPrevious().Output()
	data, cursor := terminal.buffer.Output()

	if lastData == data {
		// Move Cursor
		cursorDiff := cursor - lastCursor
		terminal.window.MoveCursor(cursorDiff)
		terminal.history.ReplaceLastBuffer(*terminal.buffer)
	} else if strings.HasPrefix(data, lastData) {
		// Insert Character
		newString := strings.Replace(data, lastData, "", 1)
		cursorDiff := cursor - lastCursor
		if cursorDiff != 1 {
			terminal.window.SaveCursor()
		}
		terminal.window.Write([]byte(newString))
		if cursorDiff != 1 {
			terminal.window.RestoreCursor()
			terminal.window.MoveCursor(cursorDiff)
		}
		terminal.history.ReplaceLastBuffer(*terminal.buffer)
	} else if strings.HasPrefix(lastData, data) {
		// Remove Character from End
		cursorDiff := cursor - lastCursor
		terminal.window.MoveCursor(cursorDiff)
		terminal.window.ClearLine(CURSOR_FORWARD)
		terminal.history.ReplaceLastBuffer(*terminal.buffer)
	} else {
		// Full Write
		terminal.window.ClearLine(FULL_LINE)
		terminal.window.SetCursorColumn(0)
		terminal.window.Write([]byte(data))
		terminal.window.SetCursorColumn(cursor + 1)
		terminal.history.ReplaceLastBuffer(*terminal.buffer)
	}
}

func (terminal Terminal) NewLine() {
	terminal.window.Write([]byte{'\n'})
	terminal.buffer.Clear()
}
