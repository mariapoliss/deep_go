package main

import (
	"cmp"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type OrderedMap[K cmp.Ordered, V any] struct {
	size    int
	element *element[K, V]
}

type element[K cmp.Ordered, V any] struct {
	key         K
	value       V
	left, right *element[K, V]
}

// создать упорядоченный словарь
func NewOrderedMap[K cmp.Ordered, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{}
}

// добавить элемент в словарь
func (m *OrderedMap[K, V]) Insert(key K, value V) {
	m.size++
	node := m.element
	for node != nil {
		if node.key < key {
			if node.right == nil {
				node.right = &element[K, V]{
					key:   key,
					value: value,
				}
				return
			} else {
				node = node.right
			}
		} else if node.key > key {
			if node.left == nil {
				node.left = &element[K, V]{
					key:   key,
					value: value,
				}
				return
			} else {
				node = node.left
			}
		}
	}
	m.element = &element[K, V]{
		key:   key,
		value: value,
	}
}

// удалить элемент из словаря
func (m *OrderedMap[K, V]) Erase(key K) {
	m.element = eraseNode(m.element, key, &m.size)
}

// eraseNode returns root of tree/subtree
func eraseNode[K cmp.Ordered, V any](node *element[K, V], key K, size *int) *element[K, V] {
	if node == nil {
		return nil
	}
	if node.key == key {
		*size--
		if node.left == nil {
			return node.right
		}
		if node.right == nil {
			return node.left
		}
		// right and left != nil
		var previous *element[K, V]
		smallest := node.right
		for smallest.left != nil {
			previous = smallest
			smallest = smallest.left
		}
		node.key = smallest.key
		node.value = smallest.value
		node.right = eraseNode(previous, smallest.key, new(int))
	} else if node.key < key {
		node.right = eraseNode(node.right, key, size)
	} else {
		node.left = eraseNode(node.left, key, size)
	}
	return node
}

// проверить существование элемента в словаре
func (m *OrderedMap[K, V]) Contains(key K) bool {
	node := m.element
	for node != nil {
		if node.key == key {
			return true
		} else if node.key < key {
			node = node.right
		} else if node.key > key {
			node = node.left
		}
	}
	return false
}

// получить количество элементов в словаре
func (m *OrderedMap[K, V]) Size() int {
	return m.size
}

// применить функцию к каждому элементу словаря от меньшего к большему
func (m *OrderedMap[K, V]) ForEach(action func(K, V)) {
	if m.element == nil {
		return
	}
	forEachInner(m.element, action)
}

func forEachInner[K cmp.Ordered, V any](node *element[K, V], action func(K, V)) {
	if node.left != nil {
		forEachInner(node.left, action)
	}
	action(node.key, node.value)
	if node.right != nil {
		forEachInner(node.right, action)
	}
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap[int, int]()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(10)
	assert.Equal(t, 3, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(5))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(10))
}
