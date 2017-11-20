package tree

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTree(t *testing.T) {
	t.Run("test int tree", func(t *testing.T) {
		// Given:
		tr := NewIntTreeFromValues([]int{4, 2, 6, 1, 3, 5, 7})
		expectedTreeStr := "    1\n" +
			"  2\n" +
			"    3\n" +
			"4\n" +
			"    5\n" +
			"  6\n" +
			"    7\n"

		// When:
		var buffer bytes.Buffer
		count, err := tr.WriteAsStringTo(&buffer)
		//fmt.Println(buffer.String())

		// Then:
		assert.Nil(t, err)
		assert.True(t, count > 0)
		assert.Equal(t, expectedTreeStr, buffer.String(), "Unexpected tree output")
		assert.Equal(t, len(expectedTreeStr), count, "Unexpected count of written bytes")
	})

	t.Run("test generic tree", func(t *testing.T) {
		// Given:
		valueWriter := func(value interface{}, writer io.Writer) (int, error) {
			valueStr, ok := value.(string)
			if !ok {
				panic("Can't convert value to string")
			}
			return io.WriteString(writer, valueStr)
		}

		valueComparator := func(left interface{}, right interface{}) int {
			lv, ok := left.(string)
			if ok {
				rv, ok := right.(string)
				if ok {
					return strings.Compare(rv, lv)
				}
			}
			panic("Can't convert value to string")
		}

		tr := NewTreeFromValues(valueWriter, valueComparator, []interface{}{
			"Dd", "Bb", "Aa", "Cc", "Ff", "Ee", "Gg",
		})
		expectedTreeStr := "    Gg\n" +
			"  Ff\n" +
			"    Ee\n" +
			"Dd\n" +
			"    Cc\n" +
			"  Bb\n" +
			"    Aa\n"

		// When:
		var buffer bytes.Buffer
		count, err := tr.WriteAsStringTo(&buffer)
		//fmt.Println(buffer.String())

		// Then:
		assert.Nil(t, err)
		assert.True(t, count > 0)
		assert.Equal(t, expectedTreeStr, buffer.String(), "Unexpected tree output")
		assert.Equal(t, len(expectedTreeStr), count, "Unexpected count of written bytes")
	})
}
