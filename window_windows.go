//go:build windows

package main

import (
	"fmt"
	"io"
	"os"

	tsize "github.com/kopoli/go-terminal-size"
)

const CSI = "\x1B["

type windowsWindow struct {
	size WindowSize
	file io.Writer
}

func NewWindow() Window {
	// TODO handle error
	terminalSize, _ := tsize.GetSize()
	return &windowsWindow{
		size: WindowSize{
			width:  terminalSize.Width,
			height: terminalSize.Height,
		},
		file: os.Stdout,
	}
}

func (window windowsWindow) GetWindowSize() WindowSize {
	return window.size
}

func (window windowsWindow) ClearLine(lineClear LineClear) {
	window.file.Write([]byte(fmt.Sprintf("%v%v%v", CSI, lineClear, "K")))
}

func (window windowsWindow) MoveCursor(x int) {
	if x < 0 {
		window.file.Write([]byte(fmt.Sprintf("%v%v%v", CSI, -1*x, "D")))
	} else if x > 0 {
		window.file.Write([]byte(fmt.Sprintf("%v%v%v", CSI, x, "C")))
	}
	if y < 0 {
		window.file.Write([]byte(fmt.Sprintf("%v%v%v", CSI, -1*y, "A")))
	} else if y > 0 {
		window.file.Write([]byte(fmt.Sprintf("%v%v%v", CSI, y, "B")))
	}
}

func (window windowsWindow) SetCursorColumn(x int) {
	window.file.Write([]byte(fmt.Sprintf("%v%v%v", CSI, x, "G")))
}

func (window windowsWindow) SaveCursor() {
	window.file.Write([]byte(fmt.Sprintf("%v%v", CSI, "s")))
}

func (window windowsWindow) RestoreCursor() {
	window.file.Write([]byte(fmt.Sprintf("%v%v", CSI, "u")))
}

func (window windowsWindow) Write(input []byte) (int, error) {
	return window.file.Write(input)
}
