//go:build windows

package main

import (
	"os"

	tsize "github.com/kopoli/go-terminal-size"
)

type windowsWindow struct {
	size WindowSize
	file *os.File
}

func NewWindow() Window {
	// TODO handle error
	terminalSize, _ := tsize.GetSize()
	return &unixWindow{
		size: WindowSize{
			width:  terminalSize.Width,
			height: terminalSize.Height,
		},
		file: os.Stdout,
	}
}

func (window unixWindow) GetWindowSize() WindowSize {
	return window.size
}

func (window *unixWindow) SetCursorPosition(x int, y int) {
	panic("TODO")
}

func (window *unixWindow) Write(input []byte) (int, error) {
	return window.file.Write(input)
}
