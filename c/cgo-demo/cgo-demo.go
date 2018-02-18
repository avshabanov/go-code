package main

/*
#include <math.h>

int sum_with_hundred(int a, int b) {
	return 100 + a + b;
}
*/
import "C"

import (
	"fmt"
)

func main() {
	r := C.sum_with_hundred(20, 1)
	fmt.Printf("C.sum_with_hundred(20, 1)=%d\n", r)

	t := C.cos(3.1415926535 / 3)
	fmt.Printf("cos(pi/3)=%f\n", t)
}
