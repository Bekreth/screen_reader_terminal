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

func NewBufferWithString(input string) Buffer {
	return Buffer{
		cursorPosition: len(input),
		currentValue:   input,
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

// Move the cursor forward by a word count, delineated by white space
func (buffer *Buffer) AdvanceCursorByWord(wordCount int) {
	indicies := append(indiciesOfChar(buffer.currentValue, ' '), len(buffer.currentValue))
	for _, i := range indicies {
		if i > buffer.cursorPosition {
			buffer.cursorPosition = i
			return
		}
	}
}

func (buffer *Buffer) RetreatCursor(amount int) {
	buffer.cursorPosition -= amount
}

// Move the cursor backwards by a word count, delineated by white space
func (buffer *Buffer) RetreatCursorByWord(wordCount int) {
	indicies := indiciesOfChar(buffer.currentValue, ' ')
	possibleIndex := 0
	for _, i := range indicies {
		if i > possibleIndex && i < buffer.cursorPosition {
			possibleIndex = i
		}
	}
	buffer.cursorPosition = possibleIndex
}

func (buffer Buffer) Output() (string, int) {
	return buffer.currentValue, buffer.cursorPosition
}

func (buffer *Buffer) Clear() {
	buffer.cursorPosition = 0
	buffer.currentValue = ""
}

func indiciesOfChar(input string, char rune) []int {
	output := make([]int, 0)
	for i, c := range input {
		if c == char {
			output = append(output, i)
		}
	}
	return output
}
