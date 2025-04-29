package main

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type RWMutex struct {
	writingCond  *sync.Cond
	readingCond  *sync.Cond
	mutex        *sync.Mutex
	readersCount int
	writing      bool
}

func NewRWMutex() *RWMutex {
	mutex := &sync.Mutex{}
	return &RWMutex{
		writingCond: sync.NewCond(mutex),
		readingCond: sync.NewCond(mutex),
		mutex:       mutex,
	}
}

func (m *RWMutex) Lock() {
	m.mutex.Lock()
	m.writing = true

	for m.readersCount > 0 {
		m.writingCond.Wait()
	}
}

func (m *RWMutex) Unlock() {
	// do not lock here, because it is already locked
	m.writing = false
	m.readingCond.Broadcast()
	m.mutex.Unlock()
}

func (m *RWMutex) RLock() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for m.writing {
		m.readingCond.Wait()
	}
	m.readersCount++
}

func (m *RWMutex) RUnlock() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.readersCount--
	if m.readersCount == 0 {
		m.writingCond.Broadcast()
	}
	if m.readersCount < 0 {
		panic("RUnlock called before RLock")
	}
}

func TestRWMutexWithWriter(t *testing.T) {
	mutex := NewRWMutex()
	mutex.Lock() // writer

	var mutualExlusionWithWriter atomic.Bool
	mutualExlusionWithWriter.Store(true)
	var mutualExlusionWithReader atomic.Bool
	mutualExlusionWithReader.Store(true)

	go func() {
		mutex.Lock() // another writer
		mutualExlusionWithWriter.Store(false)
	}()

	go func() {
		mutex.RLock() // another reader
		mutualExlusionWithReader.Store(false)
	}()

	time.Sleep(time.Second)
	assert.True(t, mutualExlusionWithWriter.Load())
	assert.True(t, mutualExlusionWithReader.Load())
}

func TestRWMutexWithReaders(t *testing.T) {
	mutex := NewRWMutex()
	mutex.RLock() // reader

	var mutualExlusionWithWriter atomic.Bool
	mutualExlusionWithWriter.Store(true)

	go func() {
		mutex.Lock() // another writer
		mutualExlusionWithWriter.Store(false)
	}()

	time.Sleep(time.Second)
	assert.True(t, mutualExlusionWithWriter.Load())
}

func TestRWMutexMultipleReaders(t *testing.T) {
	mutex := NewRWMutex()
	mutex.RLock() // reader

	var readersCount atomic.Int32
	readersCount.Add(1)

	go func() {
		mutex.RLock() // another reader
		readersCount.Add(1)
	}()

	go func() {
		mutex.RLock() // another reader
		readersCount.Add(1)
	}()

	time.Sleep(time.Second)
	assert.Equal(t, int32(3), readersCount.Load())
}

func TestRWMutexWithWriterPriority(t *testing.T) {
	mutex := NewRWMutex()
	mutex.RLock() // reader

	var mutualExlusionWithWriter atomic.Bool
	mutualExlusionWithWriter.Store(true)
	var readersCount atomic.Int32
	readersCount.Add(1)

	go func() {
		mutex.Lock() // another writer is waiting for reader
		mutualExlusionWithWriter.Store(false)
	}()

	time.Sleep(time.Second)

	go func() {
		mutex.RLock() // another reader is waiting for a higher priority writer
		readersCount.Add(1)
	}()

	go func() {
		mutex.RLock() // another reader is waiting for a higher priority writer
		readersCount.Add(1)
	}()

	time.Sleep(time.Second)

	assert.True(t, mutualExlusionWithWriter.Load())
	assert.Equal(t, int32(1), readersCount.Load())
}

func TestRWMutexWithMultiplyLocking(t *testing.T) {
	mutex := NewRWMutex()
	mutex.RLock()
	mutex.RLock()
	mutex.RLock()

	mutex.RUnlock()
	mutex.RUnlock()
	mutex.RUnlock()

	mutex.Lock()
	mutex.Unlock()

	assert.True(t, true, "Reached this line without a deadlock")
}

func TestRWMutexWithUnlockingAllWriters(t *testing.T) {
	mutex := NewRWMutex()
	mutex.RLock()

	wg := &sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		mutex.Lock()
		mutex.Unlock()
	}()

	go func() {
		defer wg.Done()
		mutex.Lock()
		mutex.Unlock()
	}()

	go func() {
		defer wg.Done()
		mutex.Lock()
		mutex.Unlock()
	}()

	time.Sleep(time.Second)
	mutex.RUnlock()

	wg.Wait()
	assert.True(t, true, "Reached this line without a deadlock")
}

func TestRWMutexWithUnlockingAllReaders(t *testing.T) {
	mutex := NewRWMutex()
	mutex.Lock()

	wg := &sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		mutex.RLock()
		mutex.RUnlock()
	}()

	go func() {
		defer wg.Done()
		mutex.RLock()
		mutex.RUnlock()
	}()

	go func() {
		defer wg.Done()
		mutex.RLock()
		mutex.RUnlock()
	}()

	time.Sleep(time.Second)
	mutex.Unlock()

	wg.Wait()
	assert.True(t, true, "Reached this line without a deadlock")
}
