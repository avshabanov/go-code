package glidesample

import (
	"fmt"
)

// Foo is a sample exported function
func Foo() string {
	return fmt.Sprintf("Hello from Glidesample library version %s", Version)
}

// Foox is a sample exported function
func Foox(a, b int) int {
	return 100 + a + b
}
