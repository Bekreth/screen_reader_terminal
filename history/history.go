package history

import (
	"github.com/bekreth/screen_reader_terminal/buffer"
)

type History struct {
	buffers []buffer.Buffer
	index   int
}

func NewBufferHistory() History {
	return History{
		buffers: []buffer.Buffer{buffer.NewBuffer()},
		index:   0,
	}
}

func (history History) GetPrevious() buffer.Buffer {
	if len(history.buffers) == 0 {
		return buffer.NewBuffer()
	}
	return history.buffers[history.index]
}

func (history *History) AddBuffer(buffer buffer.Buffer) {
	history.buffers = append(history.buffers, buffer)
	history.index += 1
}

func (history *History) ReplaceLastBuffer(buffer buffer.Buffer) {
	history.buffers[history.index] = buffer
}
