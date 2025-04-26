package main

import (
	"container/heap"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Tasks struct {
	data    []Task
	indexes map[int]int
}

func NewTasks() *Tasks {
	return &Tasks{
		data:    make([]Task, 0),
		indexes: make(map[int]int),
	}
}

func (t *Tasks) Len() int {
	return len(t.data)
}

func (t *Tasks) Less(i int, j int) bool {
	return t.data[i].Priority > t.data[j].Priority
}

func (t *Tasks) Swap(i int, j int) {
	t.data[i], t.data[j] = t.data[j], t.data[i]
	t.indexes[t.data[i].Identifier] = i
	t.indexes[t.data[j].Identifier] = j
}

func (t *Tasks) Push(x any) {
	task := x.(Task)
	t.data = append(t.data, task)
	t.indexes[task.Identifier] = len(t.data) - 1
}

func (t *Tasks) Pop() any {
	el := t.data[len(t.data)-1]
	t.data = t.data[:len(t.data)-1]
	delete(t.indexes, el.Identifier)
	return el
}

func (t *Tasks) Empty() bool {
	return len(t.data) == 0
}

func (t *Tasks) Exists(taskID int) bool {
	_, ok := t.indexes[taskID]
	return ok
}

func (t *Tasks) ChangePriority(taskID, priority int) (int, bool) {
	i, ok := t.indexes[taskID]
	if ok {
		t.data[i].Priority = priority
	}
	return i, ok
}

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	tasks *Tasks
}

func NewScheduler() Scheduler {
	return Scheduler{tasks: NewTasks()}
}

func (s *Scheduler) AddTask(task Task) {
	if s.tasks.Exists(task.Identifier) {
		// if task already exists - change its priority instead of adding new
		s.ChangeTaskPriority(task.Identifier, task.Priority)
	} else {
		heap.Push(s.tasks, task)
	}
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	i, ok := s.tasks.ChangePriority(taskID, newPriority)
	if ok {
		heap.Fix(s.tasks, i)
	}
}

func (s *Scheduler) GetTask() Task {
	if s.tasks.Empty() {
		return Task{}
	}
	return heap.Pop(s.tasks).(Task)
}

func TestScheduler(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)

	task = scheduler.GetTask()
	// looks like a bug we need to have priority like 100 here, because we changed it upper
	assert.Equal(t, Task{Identifier: 1, Priority: 100}, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}

func TestSchedulerAddHighPriorities(t *testing.T) {
	scheduler := NewScheduler()
	for i := 100; i > 0; i-- {
		scheduler.AddTask(Task{Identifier: i, Priority: i * 10})
	}
	for i := 100; i > 0; i-- {
		task := scheduler.GetTask()
		assert.Equal(t, Task{Identifier: i, Priority: i * 10}, task)
	}
}

func TestSchedulerAddLowPriorities(t *testing.T) {
	scheduler := NewScheduler()
	for i := 1; i <= 100; i++ {
		scheduler.AddTask(Task{Identifier: i, Priority: i * 10})
	}
	for i := 100; i > 0; i-- {
		task := scheduler.GetTask()
		assert.Equal(t, Task{Identifier: i, Priority: i * 10}, task)
	}
}

func TestChangeTaskPriority(t *testing.T) {
	scheduler := NewScheduler()
	scheduler.AddTask(Task{Identifier: 1, Priority: 10})
	scheduler.AddTask(Task{Identifier: 2, Priority: 20})
	scheduler.AddTask(Task{Identifier: 3, Priority: 30})

	scheduler.ChangeTaskPriority(2, 50)
	scheduler.ChangeTaskPriority(1, 40)

	task := scheduler.GetTask()
	assert.Equal(t, Task{Identifier: 2, Priority: 50}, task)

	task = scheduler.GetTask()
	assert.Equal(t, Task{Identifier: 1, Priority: 40}, task)

	task = scheduler.GetTask()
	assert.Equal(t, Task{Identifier: 3, Priority: 30}, task)
}

func TestAddDuplicate(t *testing.T) {
	scheduler := NewScheduler()
	scheduler.AddTask(Task{Identifier: 1, Priority: 10})
	scheduler.AddTask(Task{Identifier: 1, Priority: 15})
	scheduler.AddTask(Task{Identifier: 1, Priority: 5})

	task := scheduler.GetTask()
	assert.Equal(t, Task{Identifier: 1, Priority: 5}, task)
}
