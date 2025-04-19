package ds

import (
	"reflect"
	"testing"
)

func TestQueuePush(t *testing.T) {
	queue := NewQueue[int]()
	queue.Push(1)
	queue.Push(2)
	queue.Push(3)

	if queue.Size() != 3 {
		t.Errorf("The queue was expected to be 3 but actually was %d :(", queue.Size())
	}

	if *queue.Peek() != 1 {
		t.Errorf("The front of the queue was expected to be 3 but actually was %d :(", *queue.Peek())
	}
}

func TestQueuePop(t *testing.T) {
	queue := NewQueue[int]()
	queue.Push(1)
	queue.Push(2)
	queue.Push(3)

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
	queue := NewQueue[int]()
	if !queue.IsEmpty() {
		t.Errorf("The queue was expected to be empty after creating new stack :(")
	}

	queue.Push(1)
	if queue.IsEmpty() {
		t.Errorf("The queue was supposed to have an element in it after adding an element :(")
	}

	queue.Pop()
	if !queue.IsEmpty() {
		t.Errorf("The queue was expected to be empty after removing an element :(")
	}
}

func TestQueueLimit(t *testing.T) {
	queue := NewQueue[int]()
	queue.Push(1)
	queue.Push(2)
	queue.Push(3)

	queue.Limit(2)

	if queue.Size() != 2 {
		t.Errorf("The queue was expected to be limited to 2 elements but was %d", queue.Size())
	}

	value := queue.Pop()
	if value == nil || *value != 2 {
		t.Errorf("The first value in the queue was expected to be 2 but was %d", value)
	}

	value = queue.Pop()
	if value == nil || *value != 3 {
		t.Errorf("The second value in the queue was expected to be 3 but was %d", value)
	}
}

func TestQueueContents(t *testing.T) {
	queue := NewQueue[int]()
	queue.Push(1)
	queue.Push(2)
	queue.Push(3)

	data := queue.Contents()
	if !reflect.DeepEqual(data, []int{1, 2, 3}) {
		t.Errorf("Expected data to be [1, 2, 3] but got %+v", data)
	}
}
