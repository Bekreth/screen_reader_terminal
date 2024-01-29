package main

type Buffer struct {
	cursorPosition int
	currentValue   string
}

func NewBuffer() Buffer {
	return Buffer{
		cursorPosition: 0,
		currentValue:   "",
	}
}

// Adds a character to the current cursor position, advancing the cursor by 1
func (buffer *Buffer) AddCharacter(character rune) {
	buffer.currentValue = buffer.currentValue[0:buffer.cursorPosition] +
		string(character) +
		buffer.currentValue[buffer.cursorPosition:]
	buffer.cursorPosition += 1
}

// Removes the character before the current cursor position if a character exists and
// retreats the cursor by 1
func (buffer *Buffer) RemoveCharacter() {
	if buffer.cursorPosition != 0 {
		buffer.cursorPosition -= 1
		buffer.currentValue = buffer.currentValue[0:buffer.cursorPosition] +
			buffer.currentValue[buffer.cursorPosition+1:]
	}
}

func (buffer *Buffer) AdvanceCursor(amount int) {
	buffer.cursorPosition += amount
}

func (buffer *Buffer) RetreatCursor(amount int) {
	buffer.cursorPosition -= amount
}
