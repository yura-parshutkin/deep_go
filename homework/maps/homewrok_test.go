package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type node struct {
	key         int
	value       int
	right, left *node
}

// go test -v homework_test.go

type OrderedMap struct {
	root *node
	size int
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{}
}

func (m *OrderedMap) Insert(key, value int) {
	if m.root == nil {
		m.root = &node{key: key, value: value}
		m.size = 1
	} else {
		m.size = m.size + insert(m.root, key, value)
	}
}

func (m *OrderedMap) Erase(key int) {
	var ok bool
	m.root, ok = erase(m.root, key)
	if ok {
		m.size--
	}
}

func (m *OrderedMap) Contains(key int) bool {
	return find(m.root, key) != nil
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	each(m.root, action)
}

func insert(cur *node, key, value int) int {
	if key == cur.key {
		cur.value = value
		return 0
	}
	if key < cur.key {
		if cur.left == nil {
			cur.left = &node{key: key, value: value}
			return 1
		}
		return insert(cur.left, key, value)
	} else {
		if cur.right == nil {
			cur.right = &node{key: key, value: value}
			return 1
		}
		return insert(cur.right, key, value)
	}
}

func erase(cur *node, key int) (*node, bool) {
	// use ok to indicate if the node was found and deleted
	var ok bool
	if cur == nil {
		return nil, false
	}
	if key == cur.key {
		if cur.right == nil && cur.left == nil {
			return nil, true
		}
		if cur.right == nil {
			return cur.left, true
		}
		if cur.left == nil {
			return cur.right, true
		}
		// find the minimum value in the right subtree to replace the current node
		mv := findMin(cur.right)
		cur.key = mv.key
		cur.value = mv.value
		cur.right, ok = erase(cur.right, mv.key)
	} else if key < cur.key {
		cur.left, ok = erase(cur.left, key)
	} else {
		cur.right, ok = erase(cur.right, key)
	}
	return cur, ok
}

func each(node *node, action func(int, int)) {
	if node == nil {
		return
	}
	each(node.left, action)
	action(node.key, node.value)
	each(node.right, action)
}

func findMin(cur *node) *node {
	if cur.left == nil {
		return cur
	}
	return findMin(cur.left)
}

func find(cur *node, key int) *node {
	if cur == nil || key == cur.key {
		return cur
	}
	if key < cur.key {
		return find(cur.left, key)
	} else {
		return find(cur.right, key)
	}
}

func TestOrderedMap(t *testing.T) {
	data := NewOrderedMap()
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
		m := NewOrderedMap()
		assert.Zero(t, m.Size())
	})

	t.Run("add an element", func(t *testing.T) {
		m := NewOrderedMap()
		m.Insert(10, 10)
		assert.Equal(t, 1, m.Size())
	})

	t.Run("erase last element", func(t *testing.T) {
		m := NewOrderedMap()
		m.Insert(10, 10)
		m.Erase(10)
		assert.Zero(t, m.Size())
	})

	t.Run("erase element with 1 child", func(t *testing.T) {
		m := NewOrderedMap()
		m.Insert(10, 10)
		m.Insert(5, 10)
		m.Erase(10)
		assert.Equal(t, 1, m.Size())
	})

	t.Run("erase element with 2 child", func(t *testing.T) {
		m := NewOrderedMap()
		m.Insert(5, 0)
		m.Insert(2, 0)
		m.Insert(7, 0)
		m.Insert(6, 0)

		m.Erase(5)
		assert.Equal(t, 3, m.Size())
	})
}
