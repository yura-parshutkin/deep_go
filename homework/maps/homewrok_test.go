package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/constraints"
)

type node[K constraints.Ordered, V any] struct {
	key         K
	value       V
	right, left *node[K, V]
}

type OrderedMap[K constraints.Ordered, V any] struct {
	root *node[K, V]
	size int
}

func NewOrderedMap[K constraints.Ordered, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{}
}

func (m *OrderedMap[K, V]) Insert(key K, value V) {
	if m.root == nil {
		m.root = &node[K, V]{key: key, value: value}
		m.size++
	} else {
		m.size += insert(m.root, key, value)
	}
}

func (m *OrderedMap[K, V]) Erase(key K) {
	var count int
	m.root, count = erase(m.root, key)
	m.size -= count
}

func (m *OrderedMap[K, V]) Contains(key K) bool {
	return find(m.root, key) != nil
}

func (m *OrderedMap[K, V]) Size() int {
	return m.size
}

func (m *OrderedMap[K, V]) ForEach(action func(K, V)) {
	each(m.root, action)
}

// insert returns 1 if the key was inserted, 0 if the key was updated
func insert[K constraints.Ordered, V any](cur *node[K, V], key K, value V) int {
	if key == cur.key {
		cur.value = value
		return 0
	}
	if key < cur.key {
		if cur.left == nil {
			cur.left = &node[K, V]{key: key, value: value}
			return 1
		}
		return insert(cur.left, key, value)
	} else {
		if cur.right == nil {
			cur.right = &node[K, V]{key: key, value: value}
			return 1
		}
		return insert(cur.right, key, value)
	}
}

func erase[K constraints.Ordered, V any](cur *node[K, V], key K) (*node[K, V], int) {
	// use count to track if the node was found and deleted
	var count int
	if cur == nil {
		return nil, 0
	}
	if key == cur.key {
		if cur.right == nil && cur.left == nil {
			return nil, 1
		}
		if cur.right == nil {
			return cur.left, 1
		}
		if cur.left == nil {
			return cur.right, 1
		}
		mv := findMin(cur.right)
		cur.key = mv.key
		cur.value = mv.value
		cur.right, count = erase(cur.right, mv.key)
	} else if key < cur.key {
		cur.left, count = erase(cur.left, key)
	} else {
		cur.right, count = erase(cur.right, key)
	}
	return cur, count
}

func each[K constraints.Ordered, V any](n *node[K, V], action func(K, V)) {
	if n == nil {
		return
	}
	each(n.left, action)
	action(n.key, n.value)
	each(n.right, action)
}

func findMin[K constraints.Ordered, V any](cur *node[K, V]) *node[K, V] {
	if cur.left == nil {
		return cur
	}
	return findMin(cur.left)
}

func find[K constraints.Ordered, V any](cur *node[K, V], key K) *node[K, V] {
	if cur == nil || key == cur.key {
		return cur
	}
	if key < cur.key {
		return find(cur.left, key)
	}
	return find(cur.right, key)
}

func TestOrderedMap(t *testing.T) {
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
}

func TestOrderedMap_HappyPaths(t *testing.T) {
	t.Run("map is empty", func(t *testing.T) {
		m := NewOrderedMap[int, int]()
		assert.Zero(t, m.Size())
	})

	t.Run("add an element", func(t *testing.T) {
		m := NewOrderedMap[string, float32]()
		m.Insert("10", 10.111)
		assert.Equal(t, 1, m.Size())
	})

	t.Run("erase last element", func(t *testing.T) {
		m := NewOrderedMap[int, string]()
		m.Insert(10, "10")
		m.Erase(10)
		assert.Zero(t, m.Size())
	})

	t.Run("erase element with 1 child", func(t *testing.T) {
		m := NewOrderedMap[int, bool]()
		m.Insert(10, true)
		m.Insert(5, false)
		m.Erase(10)
		assert.Equal(t, 1, m.Size())
	})

	t.Run("erase element with 2 child", func(t *testing.T) {
		m := NewOrderedMap[int, int]()
		m.Insert(5, 0)
		m.Insert(2, 0)
		m.Insert(7, 0)
		m.Insert(6, 0)

		m.Erase(5)
		assert.Equal(t, 3, m.Size())
	})
}
