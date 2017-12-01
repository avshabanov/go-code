package protobufdemo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProtodemo(t *testing.T) {
	t.Run("foo", func(t *testing.T) {
		result := Foo("a")
		assert.Equal(t, "a1", result)
	})

	t.Run("new bob profile", func(t *testing.T) {
		result := NewBobProfile()
		assert.Equal(t, "bob", result.GetName().GetFirst())
	})
}
