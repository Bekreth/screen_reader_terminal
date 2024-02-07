package main

type LineClear int

const (
	CURSOR_FORWARD   LineClear = 0
	CURSOR_BACKWARDS LineClear = 1
	FULL_LINE        LineClear = 2
)

type WindowSize struct {
	width  int
	height int
}

type Window interface {
	GetWindowSize() WindowSize
	ClearLine(LineClear)
	MoveCursor(x int, y int)
	SetCursorColumn(x int)
	SaveCursor()
	RestoreCursor()
	Write([]byte) (int, error)
}
