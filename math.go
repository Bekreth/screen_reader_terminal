package screen_reader_terminal

func IntFloor(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

// Adds 2 numbers (either positive or negative) and returns the mod and rollover count.
func ModAdd(value int, accumulator int, base int) (int, int) {
	rollover := 0

	if value < 0 {
		// base compliment
		value *= -1
		rollover -= 1
		rollover -= value / base
		value = base - (value % base)
	} else {
		rollover += value / base
		value = value % base
	}

	if accumulator < 0 {
		// base compliment
		accumulator *= -1
		rollover -= 1
		rollover -= accumulator / base
		accumulator = base - (accumulator % base)
	} else {
		rollover += accumulator / base
		accumulator = accumulator % base
	}
	sum := value + accumulator
	rollover += sum / base
	return sum % base, rollover
}
