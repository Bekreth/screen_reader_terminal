package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddCharacter(t *testing.T) {
	input := []rune{'a', 'b', 'c', 'd'}

	trials := []struct {
		description    string
		startingBuffer Buffer
		expectedOutput Buffer
	}{
		{
			description: "Empty string successful additions",
			startingBuffer: Buffer{
				cursorPosition: 0,
				currentValue:   "",
			},
			expectedOutput: Buffer{
				cursorPosition: 4,
				currentValue:   "abcd",
			},
		},
		{
			description: "Existing buffer successful additions",
			startingBuffer: Buffer{
				cursorPosition: 4,
				currentValue:   "1234",
			},
			expectedOutput: Buffer{
				cursorPosition: 8,
				currentValue:   "1234abcd",
			},
		},
		{
			description: "Existing buffer successful insertions",
			startingBuffer: Buffer{
				cursorPosition: 4,
				currentValue:   "1234zxcv",
			},
			expectedOutput: Buffer{
				cursorPosition: 8,
				currentValue:   "1234abcdzxcv",
			},
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			actualOutput := trial.startingBuffer
			for _, character := range input {
				actualOutput.AddCharacter(character)
			}
			assert.Equal(tt, trial.expectedOutput, actualOutput)
		})
	}
}

func TestRemoveCharacter(t *testing.T) {
	trials := []struct {
		description    string
		startingBuffer Buffer
		expectedOutput Buffer
	}{
		{
			description: "Empty string, no change",
			startingBuffer: Buffer{
				cursorPosition: 0,
				currentValue:   "",
			},
			expectedOutput: Buffer{
				cursorPosition: 0,
				currentValue:   "",
			},
		},
		{
			description: "Cursor at zero, no change",
			startingBuffer: Buffer{
				cursorPosition: 0,
				currentValue:   "1234",
			},
			expectedOutput: Buffer{
				cursorPosition: 0,
				currentValue:   "1234",
			},
		},
		{
			description: "Cursor in middle, delete 4 character",
			startingBuffer: Buffer{
				cursorPosition: 8,
				currentValue:   "1234abcdzxcv",
			},
			expectedOutput: Buffer{
				cursorPosition: 4,
				currentValue:   "1234zxcv",
			},
		},
		{
			description: "Cursor at end, delete 4 characters",
			startingBuffer: Buffer{
				cursorPosition: 12,
				currentValue:   "1234abcdzxcv",
			},
			expectedOutput: Buffer{
				cursorPosition: 8,
				currentValue:   "1234abcd",
			},
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			actualOutput := trial.startingBuffer
			for i := 0; i < 4; i += 1 {
				actualOutput.RemoveCharacter()
			}
			assert.Equal(tt, trial.expectedOutput, actualOutput)
		})
	}
}

func TestAdvanceCursorByWord(t *testing.T) {
	trials := []struct {
		description    string
		advanceCount   int
		startingBuffer Buffer
		expectedOutput Buffer
	}{
		{
			description:    "Empty string, do nothing",
			advanceCount:   1,
			startingBuffer: NewBufferWithString(""),
			expectedOutput: NewBufferWithString(""),
		},
		{
			description:    "At end of string, do nothing",
			advanceCount:   1,
			startingBuffer: NewBufferWithString("Hello world"),
			expectedOutput: NewBufferWithString("Hello world"),
		},
		{
			description:  "Starting in the middle, go to end",
			advanceCount: 1,
			startingBuffer: Buffer{
				currentValue:   "Hello world",
				cursorPosition: 5,
			},
			expectedOutput: NewBufferWithString("Hello world"),
		},
		{
			description:  "Starting in the middle, go to end next word",
			advanceCount: 1,
			startingBuffer: Buffer{
				currentValue:   "Hello world, this is a test",
				cursorPosition: 5,
			},
			expectedOutput: Buffer{
				currentValue:   "Hello world, this is a test",
				cursorPosition: 12,
			},
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			actualOutput := trial.startingBuffer
			actualOutput.AdvanceCursorByWord(trial.advanceCount)
			assert.Equal(tt, trial.expectedOutput, actualOutput)
		})
	}

}

func TestRetreatCursorByWord(t *testing.T) {
	trials := []struct {
		description    string
		advanceCount   int
		startingBuffer Buffer
		expectedOutput Buffer
	}{
		{
			description:    "Empty string, do nothing",
			advanceCount:   1,
			startingBuffer: NewBufferWithString(""),
			expectedOutput: NewBufferWithString(""),
		},
		{
			description:  "At beginning of string, do nothing",
			advanceCount: 1,
			startingBuffer: Buffer{
				currentValue:   "Hello world",
				cursorPosition: 0,
			},
			expectedOutput: Buffer{
				currentValue:   "Hello world",
				cursorPosition: 0,
			},
		},
		{
			description:  "Starting in the middle, beginning",
			advanceCount: 1,
			startingBuffer: Buffer{
				currentValue:   "Hello world",
				cursorPosition: 5,
			},
			expectedOutput: Buffer{
				currentValue:   "Hello world",
				cursorPosition: 0,
			},
		},
		{
			description:  "Starting in the middle, go to end next word",
			advanceCount: 1,
			startingBuffer: Buffer{
				currentValue:   "Hello world, this is a test",
				cursorPosition: 10,
			},
			expectedOutput: Buffer{
				currentValue:   "Hello world, this is a test",
				cursorPosition: 5,
			},
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			actualOutput := trial.startingBuffer
			actualOutput.RetreatCursorByWord(trial.advanceCount)
			assert.Equal(tt, trial.expectedOutput, actualOutput)
		})
	}
}
