package main

import (
	"strings"
	"fmt"
)

var TEXT = `package main

import (
	"strings"
	"fmt"
)

K

func main() {
	kIndex := strings.Index(TEXT, "K")
	fmt.Print(TEXT[:kIndex])
	fmt.Print("var TEXT = \x60")
	fmt.Print(TEXT)
	fmt.Print("\x60")
	fmt.Println(TEXT[kIndex + 1:])
}`

func main() {
	kIndex := strings.Index(TEXT, "K")
	fmt.Print(TEXT[:kIndex])
	fmt.Print("var TEXT = \x60")
	fmt.Print(TEXT)
	fmt.Print("\x60")
	fmt.Println(TEXT[kIndex + 1:])
}
