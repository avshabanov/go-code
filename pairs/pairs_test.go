package pairs

import (
	"container/heap"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPairs(t *testing.T) {
	t.Run("pairs heap", func(t *testing.T) {
		pairs := []*IntPair{
			&IntPair{1, 2},
			&IntPair{3, 4},
			&IntPair{2, 4},
			&IntPair{1, 5},
			&IntPair{5, 7},
			&IntPair{6, 7},
			&IntPair{8, 9},
		}

		h := &PairRightMinHeap{}
		heap.Init(h)

		for _, p := range pairs {
			heap.Push(h, p)
		}

		resultPairTo := []int{}
		for h.Len() > 0 {
			p, ok := heap.Pop(h).(*IntPair)
			assert.True(t, ok)

			resultPairTo = append(resultPairTo, p.To)
		}

		assert.Equal(t, []int{2, 4, 4, 5, 7, 7, 9}, resultPairTo)
	})

	t.Run("solution finder", func(t *testing.T) {
		solutions := make(map[string]int)

		var dumper SolutionDumper = func(intersectedPairs []*IntPair) {
			// transfer array set into ad-hoc string to ease assertions
			strs := []string{}
			for _, p := range intersectedPairs {
				strs = append(strs, fmt.Sprintf("%d-%d", p.From, p.To))
			}
			sort.Slice(strs, func(i, j int) bool {
				return strs[i] < strs[j]
			})

			pairStrs := strings.Join(strs, ";")

			//fmt.Printf("dump set: %s\n", pairStrs)
			solutions[pairStrs] = 1
		}
		pairs := []*IntPair{
			&IntPair{1, 2},
			&IntPair{3, 4},
			&IntPair{2, 4},
			&IntPair{1, 5},
			&IntPair{5, 7},
			&IntPair{6, 7},
			&IntPair{8, 9},
		}
		//fmt.Printf("input pairs:\n\t%s\n", pairsToString(pairs))

		FindSolution(pairs, dumper)

		assert.Equal(t, 5, len(solutions))

		assert.Equal(t, 1, solutions["1-2;1-5;2-4"])
		assert.Equal(t, 1, solutions["1-5;2-4;3-4"])
		assert.Equal(t, 1, solutions["1-5;5-7"])
		assert.Equal(t, 1, solutions["5-7;6-7"])
		assert.Equal(t, 1, solutions["8-9"])
	})
}
