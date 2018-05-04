package fizzbuzz

// FizzBuzzOutput is a callback function type for FizzBuzz function
type FizzBuzzOutput func(index int, output string)

// FizzBuzz function
func FizzBuzz(from, to int, outputFunc FizzBuzzOutput) {
	m3 := from % 3
	m5 := from % 5
	for i := from; i < to; i++ {
		output := ""

		if m3 == 0 {
			output = "Fizz"
		}

		if m5 == 0 {
			if len(output) > 0 {
				output = "FizzBuzz"
			} else {
				output = "Buzz"
			}
		}

		if len(output) > 0 {
			outputFunc(i, output)
		}

		m3 = (m3 + 1) % 3
		m5 = (m5 + 1) % 5
	}
}
