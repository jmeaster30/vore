package ds

import (
	"testing"

	"github.com/jmeaster30/vore/libvore/testutils"
)

func TestSubsliceStack(t *testing.T) {
	stack := NewStack[int]()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	stack.Push(4)
	stack.Push(5)

	subslice := Subslice[int](stack, 1, 3)
	testutils.AssertLength(t, 3, subslice)
	testutils.AssertEqual(t, []int{4, 3, 2}, subslice)
}

func TestSubsliceStack_Overflow(t *testing.T) {
	stack := NewStack[int]()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	stack.Push(4)
	stack.Push(5)

	subslice := Subslice[int](stack, 3, 10)
	testutils.AssertLength(t, 2, subslice)
	testutils.AssertEqual(t, []int{2, 1}, subslice)
}
