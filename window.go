package main

type WindowSize struct {
	width  int
	height int
}

type Window interface {
	GetWindowSize() WindowSize
	SetCursorPosition(x int, y int)
	Write([]byte) (int, error)
}
