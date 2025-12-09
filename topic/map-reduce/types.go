package main

import (
	"time"
)

type KeyValue struct {
	Key   string
	Value string
}

type (
	MapFunc    func(docID string, contents string) []KeyValue
	ReduceFunc func(key string, values []string) string
)

type TaskType int

const (
	MapTask TaskType = iota
	ReduceTask
	NoTask
)

type Task struct {
	Type       TaskType
	TaskID     int
	ChunkIndex int
	Data       string
}

type TaskState struct {
	State      string
	AssignedAt time.Time
}
