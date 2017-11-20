package computation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type outputPair struct {
	Index  int
	Output string
}

func TestFizzBuzz(t *testing.T) {
	t.Run("fizz buzz 35 1..10", func(t *testing.T) {
		expectedPairs := []outputPair{
			{3, "Fizz"},
			{5, "Buzz"},
			{6, "Fizz"},
			{9, "Fizz"},
			{10, "Buzz"},
			{12, "Fizz"},
			{15, "FizzBuzz"},
			{18, "Fizz"},
			{20, "Buzz"},
			{21, "Fizz"},
			{24, "Fizz"},
			{25, "Buzz"},
			{27, "Fizz"},
			{30, "FizzBuzz"},
		}
		actualPairs := []outputPair{}

		FizzBuzz(1, 31, func(index int, output string) {
			//fmt.Printf("{%d, \"%s\"}\n", index, output)
			actualPairs = append(actualPairs, outputPair{index, output})
		})

		assert.Equal(t, len(actualPairs), len(expectedPairs), "pair slice length mismatch")
		for index, actualPair := range actualPairs {
			expectedPair := expectedPairs[index]
			assert.Equal(t, expectedPair.Index, actualPair.Index, "index mismatch")
			assert.Equal(t, expectedPair.Output, actualPair.Output, "output mismatch")
		}
	})
}
