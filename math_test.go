package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModAdd(t *testing.T) {
	trials := []struct {
		description      string
		inputValue       int
		inputAccumulator int
		inputBase        int
		expectedOutput   int
		expectedRollover int
	}{
		{
			description:      "Add value less than mod",
			inputValue:       5,
			inputAccumulator: 6,
			inputBase:        25,
			expectedOutput:   11,
			expectedRollover: 0,
		},
		{
			description:      "Add value exceeds one revolution than mod",
			inputValue:       20,
			inputAccumulator: 6,
			inputBase:        25,
			expectedOutput:   1,
			expectedRollover: 1,
		},
		{
			description:      "Add value exceeds multiple revolution than mod",
			inputValue:       26,
			inputAccumulator: 6,
			inputBase:        25,
			expectedOutput:   7,
			expectedRollover: 1,
		},
		{
			description:      "Add one negative value without rollover",
			inputValue:       -10,
			inputAccumulator: 16,
			inputBase:        25,
			expectedOutput:   6,
			expectedRollover: 0,
		},
		{
			description:      "Add one negative value with rollover",
			inputValue:       -10,
			inputAccumulator: 0,
			inputBase:        25,
			expectedOutput:   15,
			expectedRollover: -1,
		},
		{
			description:      "Add two negative values",
			inputValue:       -10,
			inputAccumulator: -6,
			inputBase:        25,
			expectedOutput:   9,
			expectedRollover: -1,
		},
		{
			description:      "Add large negative numbers",
			inputValue:       -101,
			inputAccumulator: -63,
			inputBase:        25,
			expectedOutput:   11,
			expectedRollover: -7,
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			actualOutput, actualRollover := ModAdd(
				trial.inputValue,
				trial.inputAccumulator,
				trial.inputBase,
			)
			assert.Equal(tt, trial.expectedOutput, actualOutput)
			assert.Equal(tt, trial.expectedRollover, actualRollover)
		})
	}
}
