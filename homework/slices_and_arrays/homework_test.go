package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type QueueElement interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// go test -v homework_test.go
type CircularQueue[T QueueElement] struct {
	values    []T
	len, head int
}

func NewCircularQueue[T QueueElement](size int) CircularQueue[T] {
	return CircularQueue[T]{
		values: make([]T, size),
	}
}

func (q *CircularQueue[T]) Push(value T) bool {
	if q.Full() {
		return false
	}
	tail := (q.head + q.len) % cap(q.values)
	q.values[tail] = value
	q.len++
	return true
}

func (q *CircularQueue[T]) Pop() bool {
	if q.Empty() {
		return false
	}
	q.head = (q.head + 1) % cap(q.values)
	q.len--
	return true
}

func (q *CircularQueue[T]) Front() T {
	if q.Empty() {
		return T(-1)
	}
	return q.values[q.head]
}

func (q *CircularQueue[T]) Back() T {
	if q.Empty() {
		return T(-1)
	}
	tail := (q.head + q.len - 1) % cap(q.values)
	return q.values[tail]
}

func (q *CircularQueue[T]) Empty() bool {
	return q.len == 0
}

func (q *CircularQueue[T]) Full() bool {
	return q.len == cap(q.values)
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue[int](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}
