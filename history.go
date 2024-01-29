package main

type History struct {
	buffers []Buffer
	index   int
}

func NewBufferHistory() History {
	return History{
		buffers: make([]Buffer, 1000),
	}
}

func (history History) GetPrevious() Buffer {
	if len(history.buffers) == 0 {
		return NewBuffer()
	}
	return history.buffers[history.index]
}

func (history *History) AddBuffer(buffer Buffer) {
	history.buffers = append(history.buffers, buffer)
	history.index += 1
}

func (history *History) ReplaceLastBuffer(buffer Buffer) {
	history.buffers[history.index] = buffer
}
