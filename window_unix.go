//go:build linux

package screen_reader_terminal

import (
	"fmt"
	"io"
	"os"

	tsize "github.com/kopoli/go-terminal-size"
)

const CSI = "\x1B["

type unixWindow struct {
	size WindowSize
	file io.Writer
}

func NewWindow() Window {
	// TODO handle error
	terminalSize, _ := tsize.GetSize()
	terminalSize.Width = 20
	return unixWindow{
		size: WindowSize{
			Width:  terminalSize.Width,
			Height: terminalSize.Height,
		},
		file: os.Stdout,
	}
}

func (window unixWindow) GetWindowSize() WindowSize {
	return window.size
}

func (window unixWindow) ClearLine(lineClear LineClear) {
	window.file.Write([]byte(fmt.Sprintf("%v%v%v", CSI, lineClear, "K")))
}

func (window unixWindow) MoveCursor(x int, y int) {
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

func (window unixWindow) SetCursorColumn(x int) {
	window.file.Write([]byte(fmt.Sprintf("%v%v%v", CSI, x, "G")))
}

func (window unixWindow) SaveCursor() {
	window.file.Write([]byte(fmt.Sprintf("%v%v", CSI, "s")))
}

func (window unixWindow) RestoreCursor() {
	window.file.Write([]byte(fmt.Sprintf("%v%v", CSI, "u")))
}

func (window unixWindow) Write(input []byte) (int, error) {
	return window.file.Write(input)
}

func (window unixWindow) ScrollPage(input int) {
	if input > 0 {
		window.file.Write([]byte(fmt.Sprintf("%v%v%v", CSI, input, "S")))
	} else if input < 0 {
		window.file.Write([]byte(fmt.Sprintf("%v%v%v", CSI, input, "T")))
	}
}
