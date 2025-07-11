package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Number interface {
	int | int8 | int16 | int32 | int64
}

// Предположим, что эта очередь будет оперировать только положительными
// числами (отрицательные числа ей никогда не поступят на вход)
type CircularQueue[T Number] struct {
	values []T
	front  int
	rear   int
}

// создать очередь с определенным размером буффера
func NewCircularQueue[T Number](size int) CircularQueue[T] {
	return CircularQueue[T]{
		values: make([]T, size),
		front:  -1,
		rear:   -1,
	}
}

// добавить значение в конец очереди (false, если очередь заполнена)
func (q *CircularQueue[T]) Push(value T) bool {
	if q.Full() {
		return false
	}
	if q.Empty() {
		q.front = 0
		q.rear = 0
	} else {
		q.rear++
		if q.rear > cap(q.values)-1 {
			q.rear = 0
		}
	}
	q.values[q.rear] = value
	return true
}

// удалить значение из начала очереди (false, если очередь пустая)
func (q *CircularQueue[T]) Pop() bool {
	if q.Empty() {
		return false
	}
	if q.front == q.rear {
		q.front = -1
		q.rear = -1
	} else {
		q.front++
		if q.front > cap(q.values)-1 {
			q.front = 0
		}
	}
	return true
}

// получить значение из начала очереди (-1, если очередь пустая)
func (q *CircularQueue[T]) Front() T {
	if q.Empty() {
		return -1
	}

	return q.values[q.front]
}

// получить значение из конца очереди (-1, если очередь пустая)
func (q *CircularQueue[T]) Back() T {
	if q.Empty() {
		return -1
	}
	return q.values[q.rear]
}

// проверить пустая ли очередь
func (q *CircularQueue[T]) Empty() bool {
	return q.front == -1 && q.rear == -1
}

// проверить заполнена ли очередь
func (q *CircularQueue[T]) Full() bool {
	return q.front == 0 && q.rear == cap(q.values)-1 || q.front == q.rear+1
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue[int8](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, int8(-1), queue.Front())
	assert.Equal(t, int8(-1), queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int8{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, int8(1), queue.Front())
	assert.Equal(t, int8(3), queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int8{4, 2, 3}, queue.values))

	assert.Equal(t, int8(2), queue.Front())
	assert.Equal(t, int8(4), queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}
