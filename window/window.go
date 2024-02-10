package window

import "io"

type LineClear int

const (
	CURSOR_FORWARD   LineClear = 0
	CURSOR_BACKWARDS LineClear = 1
	FULL_LINE        LineClear = 2
)

type WindowSize struct {
	Width  int
	Height int
}

type Window interface {
	SetWriter(writer io.Writer) Window
	SetWindowSize(size WindowSize) Window

	GetWindowSize() WindowSize
	ClearLine(LineClear)

	MoveCursor(x int, y int)
	SetCursorColumn(x int)
	SaveCursor()
	RestoreCursor()

	// If int is positive, scrolls the page upwards by the amount shown, opposite
	// for negative
	ScrollPage(int)

	Write([]byte) (int, error)
}
