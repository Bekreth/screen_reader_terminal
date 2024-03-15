package terminal

import (
	"fmt"
	"io"

	"github.com/bekreth/screen_reader_terminal/window"
)

// test interface implmenetations

type testFile struct {
	written []byte
}

func (window *testFile) Write(input []byte) (int, error) {
	window.written = append(window.written, input...)
	return len(input), nil
}

type testWindow struct {
	window.Window
	io.Writer
}

func fmtLine(input ...any) string {
	format := ""
	for i := 0; i < len(input); i++ {
		format += "%v"
	}
	return fmt.Sprintf(format, input...)
}

func column(input int) string {
	return fmtLine(window.CSI, input, "G")
}

func up(input int) string {
	return fmtLine(window.CSI, input, "A")
}

func down(input int) string {
	return fmtLine(window.CSI, input, "B")
}

func left(input int) string {
	return fmtLine(window.CSI, input, "D")
}

func right(input int) string {
	return fmtLine(window.CSI, input, "C")
}

func save() string {
	return fmtLine(window.CSI, "s")
}

func restore() string {
	return fmtLine(window.CSI, "u")
}

func clearCursorForward() string {
	return fmtLine(window.CSI, window.CURSOR_FORWARD, "K")
}

func clearCursorFullLine() string {
	return fmtLine(window.CSI, window.FULL, "K")
}
