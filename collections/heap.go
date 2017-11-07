
package main

import (
  "container/heap"
  "fmt"
)

type MinFloatHeap []float64

func (h MinFloatHeap) Len() int           { return len(h) }
func (h MinFloatHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h MinFloatHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinFloatHeap) Push(x interface{}) {
  // use pointers to modify slice length
  *h = append(*h, x.(float64))
}

func (h *MinFloatHeap) Pop() interface{} {
  old := *h
  n := len(old)
  x := old[n - 1]
  *h = old[0 : n - 1]
  return x
}



// Use of float min-heap:

func main() {
  h := &MinFloatHeap{3.0, 1.0, 5.0}
  heap.Init(h)

  heap.Push(h, 2.0)
  fmt.Printf("Min=%f\n---\n", (*h)[0])

  for h.Len() > 0 {
    fmt.Printf("Pop element=%f\n", heap.Pop(h));
  }
}

