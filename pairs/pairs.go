package main

import (
  "fmt"
  "encoding/json"
)

type IntPair struct {
  From  int   `json:"from"`
  To    int   `json:"to"`
}

func NewIntPair(from int, to int) *IntPair {
  p := new(IntPair)
  p.From = from
  p.To = to
  return p
}

func PairsToStr(pairs []*IntPair) string {
  j, err := json.Marshal(pairs)
  if err != nil {
    return "<error>"
  }
  return string(j)
}

func main() {
  p := NewIntPair(1, 2)
  var pairs = []*IntPair {
    &IntPair{From: 1, To: 2},
    &IntPair{1, 5},
  }
  fmt.Printf("Pairs demo, p=(%d, %d)\n", p.From, p.To)
  fmt.Printf("pairs=%s\n", PairsToStr(pairs))
}

