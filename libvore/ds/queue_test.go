package ds

import "testing"

func TestQueuePush(t *testing.T) {
	queue := NewQueue[int]()
	queue.PushBack(1)
	queue.PushBack(2)
	queue.PushBack(3)

	if queue.Size() != 3 {
		t.Errorf("The queue was expected to be 3 but actually was %d :(", queue.Size())
	}

	if *queue.Peek() != 1 {
		t.Errorf("The front of the queue was expected to be 3 but actually was %d :(", *queue.Peek())
	}
}

func TestQueuePop(t *testing.T) {
	queue := NewQueue[int]()
	queue.PushBack(1)
	queue.PushBack(2)
	queue.PushBack(3)

	val := *queue.Pop()
	if val != 1 {
		t.Errorf("The queue was expected to pop 1 but actually was %d :(", val)
	}

	val = *queue.Pop()
	if val != 2 {
		t.Errorf("The queue was expected to pop 2 but actually was %d :(", val)
	}

	val = *queue.Pop()
	if val != 3 {
		t.Errorf("The queue was expected to pop 3 but actually was %d :(", val)
	}

	last := queue.Pop()
	if last != nil {
		t.Errorf("The queue was expected to pop nil but actually was %d :(", last)
	}
}

func TestQueuePushFront(t *testing.T) {
	queue := NewQueue[int]()
	queue.PushFront(1)
	queue.PushFront(2)
	queue.PushFront(3)

	val := *queue.Pop()
	if val != 3 {
		t.Errorf("The queue was expected to pop 3 but actually was %d :(", val)
	}

	val = *queue.Pop()
	if val != 2 {
		t.Errorf("The queue was expected to pop 2 but actually was %d :(", val)
	}

	val = *queue.Pop()
	if val != 1 {
		t.Errorf("The queue was expected to pop 1 but actually was %d :(", val)
	}

	last := queue.Peek()
	if last != nil {
		t.Errorf("The queue was expected to return nil but actually was %d :(", last)
	}
}

func TestQueueIsEmpty(t *testing.T) {
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
