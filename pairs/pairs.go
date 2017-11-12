package pairs

import (
	"bytes"
	"container/heap"
	"fmt"
	"sort"
)

// IntPair defines a range starting at From and ending at To
type IntPair struct {
	From int
	To   int
}

func pairsToString(pairs []*IntPair) string {
	var buffer bytes.Buffer
	for index, pair := range pairs {
		if index > 0 {
			buffer.WriteString(", ")
		}
		fmt.Fprintf(&buffer, "(%d - %d)", pair.From, pair.To)
	}

	return buffer.String()
}

// PairRightMinHeap is a helper type representing min heap for To-field of pairs
type PairRightMinHeap []*IntPair

func (h PairRightMinHeap) Len() int           { return len(h) }
func (h PairRightMinHeap) Less(i, j int) bool { return h[i].To < h[j].To }
func (h PairRightMinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// Push is a function, that pushes element from a given array
func (h *PairRightMinHeap) Push(x interface{}) {
	pair, ok := x.(*IntPair)
	if !ok {
		panic("Given argument is not of type *intPair")
	}

	*h = append(*h, pair)
}

// Pop is a function, that pops element from a given array
func (h *PairRightMinHeap) Pop() interface{} {
	if len(*h) == 0 {
		panic("Heap is empty")
	}

	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// end of pairMinHeap definition

// SolutionDumper is a function, that dumps found set of intersecting pairs
type SolutionDumper func(intersectedPairs []*IntPair)

// FindSolution finds all intersecting pair sets
func FindSolution(pairs []*IntPair, dumpFn SolutionDumper) {
	// copy pairs and sort
	sortedPairs := make([]*IntPair, len(pairs))
	copy(sortedPairs, pairs)
	sort.Slice(sortedPairs, func(i, j int) bool {
		return sortedPairs[i].From < sortedPairs[j].From
	})

	pairHeap := &PairRightMinHeap{}
	heap.Init(pairHeap)

	for _, pair := range sortedPairs {
		//fmt.Printf("\n--> inspecting %s\n", pairsToString([]*IntPair{pair}))

		inspected := false

		for pairHeap.Len() > 0 {
			candidatePair, ok := heap.Pop(pairHeap).(*IntPair)
			if !ok {
				panic("Can't convert popped element to int pair")
			}
			//fmt.Printf("comparing with pair: %s\n", pairsToString([]*IntPair{candidatePair}))

			// does this pair belong to the pair heap?
			if pair.From > candidatePair.To {
				// we need to dump pairHeap array and continue to the next candidatePair
				if !inspected {
					dumpFn(append(*pairHeap, candidatePair))
					inspected = true
				}

				//fmt.Println("\tpop")
			} else {
				// put back previously extracted candidate pair
				heap.Push(pairHeap, candidatePair)
				break
			}
		}

		//fmt.Printf("adding pair: %s\n", pairsToString([]*IntPair{pair}))
		heap.Push(pairHeap, pair)
	}

	if pairHeap.Len() > 0 {
		dumpFn(*pairHeap)
	}
}
