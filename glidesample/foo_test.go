package glidesample

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFoo(t *testing.T) {
	t.Run("foo", func(t *testing.T) {
		result := Foo()
		assert.True(t, strings.Contains(result, "Glidesample"))
	})

	t.Run("foox", func(t *testing.T) {
		// Given:
		a := 1
		b := 10

		// When:
		result := Foox(a, b)

		// Then:
		assert.Equal(t, 111, result)
	})
}
