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

	if stack.Size() != 3 {
		t.Errorf("The stack was expected to be 3 but actually was %d :(", stack.Size())
	}

	if *stack.Peek() != 3 {
		t.Errorf("The top of the stack was expected to be 3 but actually was %d :(", *stack.Peek())
	}
}

func TestPop(t *testing.T) {
	stack := NewStack[int]()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	val := *stack.Pop()
	if val != 3 {
		t.Errorf("The stack was expected to pop 3 but actually was %d :(", val)
	}

	val = *stack.Pop()
	if val != 2 {
		t.Errorf("The stack was expected to pop 2 but actually was %d :(", val)
	}

	val = *stack.Pop()
	if val != 1 {
		t.Errorf("The stack was expected to pop 1 but actually was %d :(", val)
	}

	last := stack.Pop()
	if last != nil {
		t.Errorf("The stack was expected to pop nil but actually was %d :(", last)
	}
}

func TestIsEmpty(t *testing.T) {
	stack := NewStack[int]()
	if !stack.IsEmpty() {
		t.Errorf("The stack was expected to be empty after creating new stack :(")
	}

	stack.Push(1)
	if stack.IsEmpty() {
		t.Errorf("The stack was supposed to have an element in it after adding an element :(")
	}

	stack.Pop()
	if !stack.IsEmpty() {
		t.Errorf("The stack was expected to be empty after removing an element :(")
	}
}

func TestIndex(t *testing.T) {
	stack := NewStack[int]()
	if !stack.IsEmpty() {
		t.Errorf("The stack was expected to be empty after creating new stack :(")
	}

	if stack.Index(0) != nil {
		t.Errorf("Index was suppose to return nil when the stack is empty")
	}

	stack.Push(1)
	stack.Push(2)
	stack.Index(0)
	if stack.Size() != 2 {
		t.Errorf("Index is not supposed to change the size of the stack")
	}

	if *stack.Index(1) != 2 {
		t.Errorf("Expected 2 but got %d", *stack.Index(1))
	}
}

func TestCopy(t *testing.T) {
	stack := NewStack[int]()
	stack.Push(1)
	stack.Push(2)

	stackCopy := stack.Copy()
	testutils.AssertEqual(t, 2, stackCopy.Size())
	testutils.AssertEqual(t, 2, stackCopy.Pop())
	testutils.AssertEqual(t, 1, stackCopy.Pop())
	testutils.AssertEqual(t, 2, stack.Pop())
	testutils.AssertEqual(t, 1, stack.Pop())
}

func TestPeekEmpty(t *testing.T) {
	stack := NewStack[int]()

	var expectedValue *int = nil

	testutils.AssertEqual(t, expectedValue, stack.Peek())
}
