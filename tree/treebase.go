package tree

import (
	"fmt"
	"io"
)

// Node is a base class for trees
type Node struct {
	Value interface{}
	Left  *Node
	Right *Node
}

// ValueWriter ad hoc value-to-string conversion function
type ValueWriter func(value interface{}, writer io.Writer) (int, error)

// ValueComparator is a dedicated comparator function
type ValueComparator func(left interface{}, right interface{}) int

// Tree represents binary tree
type Tree struct {
	Root              *Node
	ValueWriterFn     ValueWriter
	ValueComparatorFn ValueComparator
}

// NewTree creates new instance of a tree
func NewTree(valueWriter ValueWriter, valueComparator ValueComparator) *Tree {
	return &Tree{
		Root:              nil,
		ValueComparatorFn: valueComparator,
		ValueWriterFn:     valueWriter,
	}
}

// NewTreeFromValues creates new instance of a tree from given values
func NewTreeFromValues(
	valueWriter ValueWriter,
	valueComparator ValueComparator,
	values []interface{}) *Tree {
	tree := NewTree(valueWriter, valueComparator)

	for _, value := range values {
		newNode := &Node{Value: value}
		if tree.Root == nil {
			tree.Root = newNode
			continue
		}

		node := tree.Root
		for {
			cmp := valueComparator(value, node.Value)
			if cmp < 0 {
				if node.Left == nil {
					node.Left = newNode
					break
				} else {
					node = node.Left
				}
			} else if cmp > 0 {
				if node.Right == nil {
					node.Right = newNode
					break
				} else {
					node = node.Right
				}
			} else {
				panic("Duplicate value in a given array")
			}
		}
	}

	return tree
}

// NewIntTreeFromValues creates a tree from integer slice
func NewIntTreeFromValues(values []int) *Tree {
	uv := []interface{}{}
	for _, v := range values {
		uv = append(uv, v)
	}

	valueWriter := func(value interface{}, writer io.Writer) (int, error) {
		return io.WriteString(writer, fmt.Sprintf("%d", value))
	}

	valueComparator := func(left interface{}, right interface{}) int {
		lv, ok := left.(int)
		if ok {
			rv, ok := right.(int)
			if ok {
				return lv - rv
			}
		}
		panic("Can't convert value to int")
	}

	return NewTreeFromValues(valueWriter, valueComparator, uv)
}

// INDENT is a single indentation unit
var INDENT = []byte{' ', ' '}

// NEWLINE is a line break between tree value entries
var NEWLINE = []byte{'\n'}

// helper function for dumping tree nodes
func writeNodeTo(indent int, node *Node, valueWriter ValueWriter, w io.Writer) (int, error) {
	if node == nil {
		return 0, nil
	}

	var result int
	var n int
	var err error

	// write left subtree
	n, err = writeNodeTo(indent+1, node.Left, valueWriter, w)
	result += n
	if err != nil {
		return n, err
	}

	// write indent
	for i := 0; i < indent; i++ {
		n, err = w.Write(INDENT)
		result += n
		if err != nil {
			return result, err
		}
	}
	// write value itself
	n, err = valueWriter(node.Value, w)
	result += n
	if err != nil {
		return n, err
	}
	// write newline after value
	n, err = w.Write(NEWLINE)
	result += n
	if err != nil {
		return result, err
	}

	// write right subtree
	n, err = writeNodeTo(indent+1, node.Right, valueWriter, w)
	result += n
	return result, err
}

// WriteAsStringTo writes tree to the output writer in a form of a string
func (t *Tree) WriteAsStringTo(w io.Writer) (int, error) {
	n, err := writeNodeTo(0, t.Root, t.ValueWriterFn, w)
	return n, err
}
