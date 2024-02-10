package screen_reader_terminal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndiciesOfChar(t *testing.T) {
	trials := []struct {
		description    string
		input          string
		seperator      rune
		expectedOutput []int
	}{
		{
			description:    "String of zero gives zero",
			input:          "",
			seperator:      ' ',
			expectedOutput: make([]int, 0),
		},
		{
			description:    "String without seperator",
			input:          "Somesimplestring",
			seperator:      ' ',
			expectedOutput: make([]int, 0),
		},
		{
			description:    "String with seperators",
			input:          "Some String Seprated",
			seperator:      ' ',
			expectedOutput: []int{4, 11},
		}, {
			description:    "String with seperators at both ends",
			input:          " Some String Seprated ",
			seperator:      ' ',
			expectedOutput: []int{0, 5, 12, 21},
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			actualOutput := indiciesOfChar(trial.input, trial.seperator)
			assert.Equal(tt, trial.expectedOutput, actualOutput)
		})
	}
}

func TestZip(t *testing.T) {
	trials := []struct {
		description    string
		firstList      []int
		secondList     []int
		expectedOutput []ZipResult[int]
	}{
		{
			description:    "empty inputs, no outputs",
			firstList:      []int{},
			secondList:     []int{},
			expectedOutput: []ZipResult[int]{},
		},
		{
			description: "First list is provided",
			firstList:   []int{10, 15, 21},
			secondList:  []int{},
			expectedOutput: []ZipResult[int]{
				{First: 10, Second: 0},
				{First: 15, Second: 0},
				{First: 21, Second: 0},
			},
		},
		{
			description: "Second list is provided",
			firstList:   []int{},
			secondList:  []int{10, 15, 21},
			expectedOutput: []ZipResult[int]{
				{First: 0, Second: 10},
				{First: 0, Second: 15},
				{First: 0, Second: 21},
			},
		},
		{
			description: "Lists unmatched, more in first",
			firstList:   []int{101, 102, 103, 104},
			secondList:  []int{10, 15, 21},
			expectedOutput: []ZipResult[int]{
				{First: 101, Second: 10},
				{First: 102, Second: 15},
				{First: 103, Second: 21},
				{First: 104, Second: 0},
			},
		},
		{
			description: "Lists unmatched, more in second",
			firstList:   []int{10, 15, 21},
			secondList:  []int{101, 102, 103, 104},
			expectedOutput: []ZipResult[int]{
				{First: 10, Second: 101},
				{First: 15, Second: 102},
				{First: 21, Second: 103},
				{First: 0, Second: 104},
			},
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			actualOutput := zip(trial.firstList, trial.secondList, func() int { return 0 })
			assert.Equal(t, trial.expectedOutput, actualOutput)
		})
	}
}
