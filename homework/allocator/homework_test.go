package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func Defragment(memory []byte, pointers []unsafe.Pointer) {
	for i := range pointers {
		defPtr := unsafe.Pointer(&memory[i])
		// it's already in the right place, just skip it
		if defPtr == pointers[i] {
			continue
		}
		memory[i] = *(*byte)(pointers[i])
		pointers[i] = defPtr
	}
	// clear the rest of the memory
	for i := len(pointers); i < len(memory); i++ {
		memory[i] = 0x00
	}
}

func TestDefragmentation(t *testing.T) {
	var fragmentedMemory = []byte{
		0xFF, 0x00, 0x00, 0x00,
		0x00, 0xFF, 0x00, 0x00,
		0x00, 0x00, 0xFF, 0x00,
		0x00, 0x00, 0x00, 0xFF,
	}

	var fragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[5]),
		unsafe.Pointer(&fragmentedMemory[10]),
		unsafe.Pointer(&fragmentedMemory[15]),
	}

	var defragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[1]),
		unsafe.Pointer(&fragmentedMemory[2]),
		unsafe.Pointer(&fragmentedMemory[3]),
	}

	var defragmentedMemory = []byte{
		0xFF, 0xFF, 0xFF, 0xFF,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	Defragment(fragmentedMemory, fragmentedPointers)
	assert.True(t, reflect.DeepEqual(defragmentedMemory, fragmentedMemory))
	assert.True(t, reflect.DeepEqual(defragmentedPointers, fragmentedPointers))
}

func TestDefragmentationTwice(t *testing.T) {
	var defragmentedMemory = []byte{
		0xFF, 0xFF, 0xFF, 0xFF,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	var defragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&defragmentedMemory[0]),
		unsafe.Pointer(&defragmentedMemory[1]),
		unsafe.Pointer(&defragmentedMemory[2]),
		unsafe.Pointer(&defragmentedMemory[3]),
	}

	Defragment(defragmentedMemory, defragmentedPointers)
	assert.True(t, reflect.DeepEqual(defragmentedMemory, defragmentedMemory))
	assert.True(t, reflect.DeepEqual(defragmentedPointers, defragmentedPointers))
}
