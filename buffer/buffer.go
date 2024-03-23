package buffer

import (
	"github.com/bekreth/screen_reader_terminal/utils"
)

type BufferValues struct {
	Prefix   string
	Value    string
	Position int
}

type Buffer struct {
	prefix           string
	previousPrefix   string
	currentPosition  int
	currentValue     string
	previousPosition int
	previousValue    string
}

func NewBuffer() Buffer {
	return Buffer{
		currentPosition: 0,
		currentValue:    "",
	}
}

func NewBufferWithString(input string) Buffer {
	return Buffer{
		currentPosition: len(input),
		currentValue:    input,
	}
}

func (buffer Buffer) IsEmpty() bool {
	return buffer.currentPosition == 0 && buffer.currentValue == ""
}

func (buffer *Buffer) GetPrefix() string {
	return buffer.prefix
}

func (buffer *Buffer) SetPrefix(input string) *Buffer {
	buffer.prefix = input
	return buffer
}

func (buffer *Buffer) SetString(input string) *Buffer {
	buffer.currentValue = input
	buffer.currentPosition = len(input)
	return buffer
}

func (buffer *Buffer) SetCurrentValues(input BufferValues) *Buffer {
	buffer.prefix = input.Prefix
	buffer.currentValue = input.Value
	buffer.currentPosition = input.Position
	return buffer
}

func (buffer *Buffer) SetPreviousValues(input BufferValues) *Buffer {
	buffer.previousPrefix = input.Prefix
	buffer.previousValue = input.Value
	buffer.previousPosition = input.Position
	return buffer
}

// Adds a string to the cursor position
func (buffer *Buffer) AddString(input string) {
	buffer.currentValue = buffer.currentValue[0:buffer.currentPosition] +
		input +
		buffer.currentValue[buffer.currentPosition:]
	buffer.currentPosition += len(input)
}

// Adds a character to the current cursor position, advancing the cursor by 1
func (buffer *Buffer) AddCharacter(character rune) {
	buffer.currentValue = buffer.currentValue[0:buffer.currentPosition] +
		string(character) +
		buffer.currentValue[buffer.currentPosition:]
	buffer.currentPosition += 1
}

// Removes the character before the current cursor position if a character exists and
// retreats the cursor by 1
func (buffer *Buffer) RemoveCharacter() {
	if buffer.currentPosition != 0 {
		buffer.currentPosition -= 1
		buffer.currentValue = buffer.currentValue[0:buffer.currentPosition] +
			buffer.currentValue[buffer.currentPosition+1:]
	}
}

func (buffer *Buffer) SetCursor(position int) *Buffer {
	buffer.currentPosition = position
	return buffer
}

func (buffer *Buffer) AdvanceCursor(amount int) {
	if buffer.currentPosition < len(buffer.currentValue) {
		buffer.currentPosition += amount
	}
}

// Move the cursor forward by a word count, delineated by white space
func (buffer *Buffer) AdvanceCursorByWord(wordCount int) {
	indicies := append(
		utils.IndiciesOfChar(buffer.currentValue, ' '),
		len(buffer.currentValue)-1,
	)
	for _, i := range indicies {
		if i > buffer.currentPosition {
			if i == len(buffer.currentValue) {
				buffer.currentPosition = i
			} else {
				buffer.currentPosition = i + 1
			}
			return
		}
	}
}

func (buffer *Buffer) RetreatCursor(amount int) {
	if buffer.currentPosition > 0 {
		buffer.currentPosition -= amount
	}
}

// Move the cursor backwards by a word count, delineated by white space
func (buffer *Buffer) RetreatCursorByWord(wordCount int) {
	indicies := utils.IndiciesOfChar(buffer.currentValue, ' ')
	reversedIndicies := make([]int, len(indicies))
	for i, value := range indicies {
		j := len(indicies) - i - 1
		reversedIndicies[j] = value
	}
	reversedIndicies = append(reversedIndicies, 0)
	for _, r := range reversedIndicies {
		if r+1 < buffer.currentPosition {
			if r == 0 {
				buffer.currentPosition = r
			} else {
				buffer.currentPosition = r + 1
			}
			return
		}
	}
}

func (buffer Buffer) OutputWithoutPrefix() (string, int) {
	return buffer.currentValue, buffer.currentPosition
}

func (buffer Buffer) Output() (string, int) {
	return buffer.prefix + buffer.currentValue,
		buffer.currentPosition + len(buffer.prefix)
}

func (buffer Buffer) PreviousOutput() (string, int) {
	return buffer.previousPrefix + buffer.previousValue,
		buffer.previousPosition + len(buffer.previousPrefix)
}

func (buffer *Buffer) UpdatePrevious() {
	buffer.previousPrefix = buffer.prefix
	buffer.previousValue = buffer.currentValue
	buffer.previousPosition = buffer.currentPosition
}

func (buffer *Buffer) NewLineCount() int {
	count := 0
	for _, char := range buffer.currentValue {
		if char == '\n' {
			count += 1
		}
	}
	return count
}

func (buffer *Buffer) ClearPrevious() {
	buffer.previousValue = ""
	buffer.previousPosition = 0

	buffer.previousPrefix = ""
}

func (buffer *Buffer) Clear() {
	buffer.currentValue = ""
	buffer.currentPosition = 0
	buffer.ClearPrevious()
}
