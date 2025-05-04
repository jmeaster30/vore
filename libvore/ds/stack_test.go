package ds

import (
	"testing"

	"github.com/jmeaster30/vore/libvore/testutils"
)

func TestPush(t *testing.T) {
	stack := NewStack[int]()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	testutils.AssertEqual(t, 3, stack.Size())
	testutils.AssertEqual(t, Some(3), stack.Peek())
}

func TestPop(t *testing.T) {
	stack := NewStack[int]()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	val := stack.Pop()
	testutils.AssertEqual(t, Some(3), val)

	val = stack.Pop()
	testutils.AssertEqual(t, Some(2), val)

	val = stack.Pop()
	testutils.AssertEqual(t, Some(1), val)

	last := stack.Pop()
	testutils.AssertEqual(t, None[int](), last)
}

func TestIsEmpty(t *testing.T) {
	stack := NewStack[int]()
	testutils.AssertTrue(t, stack.IsEmpty())

	stack.Push(1)
	testutils.AssertFalse(t, stack.IsEmpty())

	stack.Pop()
	testutils.AssertTrue(t, stack.IsEmpty())
}

func TestIndex(t *testing.T) {
	stack := NewStack[int]()
	testutils.AssertTrue(t, stack.IsEmpty())
	testutils.AssertEqual(t, None[int](), stack.Index(0))

	stack.Push(1)
	stack.Push(2)
	stack.Index(0)
	testutils.AssertEqual(t, 2, stack.Size())
	testutils.AssertEqual(t, Some(2), stack.Index(0))
	testutils.AssertEqual(t, Some(1), stack.Index(1))
}

func TestCopy(t *testing.T) {
	stack := NewStack[int]()
	stack.Push(1)
	stack.Push(2)

	stackCopy := stack.Copy()
	testutils.AssertEqual(t, 2, stackCopy.Size())
	testutils.AssertEqual(t, Some(2), stackCopy.Pop())
	testutils.AssertEqual(t, Some(1), stackCopy.Pop())
	testutils.AssertEqual(t, Some(2), stack.Pop())
	testutils.AssertEqual(t, Some(1), stack.Pop())
}

func TestPeekEmpty(t *testing.T) {
	stack := NewStack[int]()

	testutils.AssertEqual(t, None[int](), stack.Peek())
}
