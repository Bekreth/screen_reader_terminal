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
	lastData, _ := terminal.history.GetPrevious().Output()
	data, _ := terminal.buffer.Output()

	if strings.HasPrefix(data, lastData) {
		newString := strings.Replace(data, lastData, "", 1)
		terminal.window.Write([]byte(newString))
		terminal.history.ReplaceLastBuffer(*terminal.buffer)
	} else {
		terminal.window.Write([]byte(data))
		terminal.history.AddBuffer(*terminal.buffer)
	}

	/*
		if lastCursor != cursor {
			terminal.window.Write([]byte(fmt.Sprintf("%v%v", "\x1B", "3D")))
		}
	*/
}

func (terminal Terminal) NewLine() {
	terminal.window.Write([]byte{'\n'})
	terminal.buffer.Clear()
}
