package main

import (
	"sync"
	"time"
)

type Master struct {
	numMap    int
	NumReduce int

	inputs []string

	mapStatus    []TaskState
	reduceStatus []TaskState

	intermediate map[int]map[string][]string
	midMu        sync.Mutex

	Outputs map[int]string
	outMu   sync.Mutex

	mu sync.Mutex
}

func NewMaster(inputs []string, numReduce int) *Master {
	m := &Master{
		numMap:       len(inputs),
		NumReduce:    numReduce,
		inputs:       inputs,
		mapStatus:    make([]TaskState, len(inputs)),
		reduceStatus: make([]TaskState, numReduce),
		intermediate: make(map[int]map[string][]string),
		Outputs:      make(map[int]string),
	}

	// 1: Initialize states to "idle"
	for i := range m.mapStatus {
		m.mapStatus[i].State = "idle"
	}
	for i := range m.reduceStatus {
		m.reduceStatus[i].State = "idle"
	}

	for i := 0; i < numReduce; i++ {
		m.intermediate[i] = make(map[string][]string)
	}
	return m
}

// RequestTask assigns a task to a worker.
//
// It returns a Task indicating the type of task and its details.
// The function locks the Master's mutex to ensure exclusive access.
//
// Steps:
//  1. Assign idle map tasks:
//     If a map task is idle, it is assigned to the worker and its state is
//     changed to "in-progress". The assigned time is also recorded.
//  2. Include actual data:
//     The function returns a Task with the type "MapTask" and the task details
//     which include the task ID, the associated data chunk index, and the
//     actual data.
func (m *Master) RequestTask() Task {
	m.mu.Lock()

	// Lock the mutex to ensure exclusive access to the Master's state
	defer m.mu.Unlock()

	// 1. Assign idle map tasks
	for i := 0; i < m.numMap; i++ {
		if m.mapStatus[i].State == "idle" {
			m.mapStatus[i].State = "in-progress"
			m.mapStatus[i].AssignedAt = time.Now()
			// 2: Include actual data
			return Task{
				Type:       MapTask,
				TaskID:     i,
				ChunkIndex: i,
				Data:       m.inputs[i],
			}
		}
	}

	// 2. Reassign timed-out map tasks
	for i := 0; i < m.numMap; i++ {
		if m.mapStatus[i].State == "in-progress" {
			if time.Since(m.mapStatus[i].AssignedAt) > 5*time.Second {
				m.mapStatus[i].AssignedAt = time.Now()
				return Task{
					Type:       MapTask,
					TaskID:     i,
					ChunkIndex: i,
					Data:       m.inputs[i],
				}
			}
		}
	}

	// 3. Check if all map tasks are done
	allMapDone := true
	for i := 0; i < m.numMap; i++ {
		if m.mapStatus[i].State != "done" {
			allMapDone = false
			break
		}
	}
	if !allMapDone {
		return Task{Type: NoTask}
	}

	// 4. Assign idle reduce tasks
	for i := 0; i < m.NumReduce; i++ {
		if m.reduceStatus[i].State == "idle" {
			m.reduceStatus[i].State = "in-progress"
			m.reduceStatus[i].AssignedAt = time.Now()
			return Task{Type: ReduceTask, TaskID: i, ChunkIndex: i}
		}
	}

	// 5. Reassign timed-out reduce tasks
	for i := 0; i < m.NumReduce; i++ {
		if m.reduceStatus[i].State == "in-progress" {
			if time.Since(m.reduceStatus[i].AssignedAt) > 10*time.Second {
				m.reduceStatus[i].AssignedAt = time.Now()
				return Task{Type: ReduceTask, TaskID: i, ChunkIndex: i}
			}
		}
	}

	return Task{Type: NoTask}
}

func (m *Master) ReportMapDone(mapID int, partitions map[int][]KeyValue) {
	m.midMu.Lock()
	for r, kvs := range partitions {
		for _, kv := range kvs {
			m.intermediate[r][kv.Key] = append(m.intermediate[r][kv.Key], kv.Value)
		}
	}
	m.midMu.Unlock()

	m.mu.Lock()
	m.mapStatus[mapID].State = "done"
	m.mu.Unlock()
}

func (m *Master) ReportReduceDone(reduceID int, out string) {
	m.outMu.Lock()
	m.Outputs[reduceID] = out
	m.outMu.Unlock()

	m.mu.Lock()
	m.reduceStatus[reduceID].State = "done"
	m.mu.Unlock()
}

func (m *Master) GetReducePartition(reduceIdx int) map[string][]string {
	m.midMu.Lock()
	defer m.midMu.Unlock()

	copyMap := make(map[string][]string)
	for k, v := range m.intermediate[reduceIdx] {
		copyv := make([]string, len(v))
		copy(copyv, v)
		copyMap[k] = copyv
	}
	return copyMap
}

func (m *Master) Done() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := 0; i < m.numMap; i++ {
		if m.mapStatus[i].State != "done" {
			return false
		}
	}
	for i := 0; i < m.NumReduce; i++ {
		if m.reduceStatus[i].State != "done" {
			return false
		}
	}
	return true
}
