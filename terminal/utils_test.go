package terminal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDrawLine(t *testing.T) {
	trials := []struct {
		description      string
		previousData     string
		previousPosition int
		currentValue     string
		currentPosition  int
		expectedOutput   string
	}{
		{
			description:      "Two empty strings",
			previousData:     "",
			previousPosition: 0,
			currentValue:     "",
			currentPosition:  0,
			expectedOutput:   "",
		},
		{
			description:      "Current string is appended",
			previousData:     "hello worl",
			previousPosition: 10,
			currentValue:     "hello world",
			currentPosition:  11,
			expectedOutput:   "d",
		},
		{
			description:      "Current string is appended, cursor not at end",
			previousData:     "hello worl",
			previousPosition: 6,
			currentValue:     "hello world",
			currentPosition:  6,
			expectedOutput:   "world",
		},
		{
			description:      "Current string has a character insert",
			previousData:     "hello world",
			previousPosition: 11,
			currentValue:     "hello HERE world",
			currentPosition:  6,
			expectedOutput:   "HERE world",
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			actualOutput := lineDiff(
				trial.previousData, trial.previousPosition,
				trial.currentValue, trial.currentPosition,
			)
			assert.Equal(tt, trial.expectedOutput, actualOutput)
		})
	}
}
