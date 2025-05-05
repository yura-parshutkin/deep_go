package main

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Walker map[uintptr]struct{}

func (w Walker) Add(ptr uintptr) {
	// it's an empty pointer
	if ptr == 0 {
		return
	}
	// it's already visited
	if _, ok := w[ptr]; ok {
		return
	}
	w[ptr] = struct{}{}
	var nextPtr = *(*uintptr)(unsafe.Pointer(ptr))
	w.Add(nextPtr)
}

func (w Walker) Slice() []uintptr {
	res := make([]uintptr, 0, len(w))
	for el := range w {
		res = append(res, el)
	}
	return res
}

func Trace(stacks [][]uintptr) []uintptr {
	walker := Walker{}
	for _, stack := range stacks {
		for _, ptr := range stack {
			walker.Add(ptr)
		}
	}
	return walker.Slice()
}

func TestTrace(t *testing.T) {
	var heapObjects = []int{
		0x00, 0x00, 0x00, 0x00, 0x00,
	}

	var heapPointer1 *int = &heapObjects[1]
	var heapPointer2 *int = &heapObjects[2]
	var heapPointer3 *int = nil
	var heapPointer4 **int = &heapPointer3

	var stacks = [][]uintptr{
		{
			uintptr(unsafe.Pointer(&heapPointer1)), 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[0])),
			0x00, 0x00, 0x00, 0x00,
		},
		{
			uintptr(unsafe.Pointer(&heapPointer2)), 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[1])),
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[2])),
			uintptr(unsafe.Pointer(&heapPointer4)), 0x00, 0x00, 0x00,
		},
		{
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[3])),
		},
	}

	pointers := Trace(stacks)
	expectedPointers := []uintptr{
		uintptr(unsafe.Pointer(&heapPointer1)),
		uintptr(unsafe.Pointer(&heapObjects[0])),
		uintptr(unsafe.Pointer(&heapPointer2)),
		uintptr(unsafe.Pointer(&heapObjects[1])),
		uintptr(unsafe.Pointer(&heapObjects[2])),
		uintptr(unsafe.Pointer(&heapPointer4)),
		uintptr(unsafe.Pointer(&heapPointer3)),
		uintptr(unsafe.Pointer(&heapObjects[3])),
	}
	assert.ElementsMatch(t, expectedPointers, pointers)
}
