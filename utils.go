package screen_reader_terminal

func indiciesOfChar(input string, char rune) []int {
	output := make([]int, 0)
	for i, c := range input {
		if c == char {
			output = append(output, i)
		}
	}
	return output
}

type ZipResult[R any] struct {
	First  R
	Second R
}

func zip[R any](input1 []R, input2 []R, emptyMaker func() R) []ZipResult[R] {
	i1Len := len(input1)
	i2Len := len(input2)
	max := IntMax(i1Len, i2Len)

	output := make([]ZipResult[R], max)

	for i := 0; i < max; i++ {
		first := emptyMaker()
		second := emptyMaker()

		if i1Len > i {
			first = input1[i]
		}

		if i2Len > i {
			second = input2[i]
		}

		output[i] = ZipResult[R]{
			First:  first,
			Second: second,
		}
	}

	return output
}
