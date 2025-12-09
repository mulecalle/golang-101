package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

type Worker struct {
	ID       int
	Master   *Master
	mapFn    MapFunc
	reduceFn ReduceFunc
}

func NewWorker(id int, m *Master, mf MapFunc, rf ReduceFunc) *Worker {
	return &Worker{ID: id, Master: m, mapFn: mf, reduceFn: rf}
}

func (w *Worker) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		task := w.Master.RequestTask()
		switch task.Type {
		case MapTask:
			fmt.Printf("worker %d: got MAP task %d\n", w.ID, task.TaskID)
			w.doMap(task)
		case ReduceTask:
			fmt.Printf("worker %d: got REDUCE task %d\n", w.ID, task.TaskID)
			w.doReduce(task)
		case NoTask:
			if w.Master.Done() {
				return
			}
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func (w *Worker) doMap(task Task) {
	// 3: Use task.Data instead of re-indexing
	kvs := w.mapFn(fmt.Sprintf("doc-%d", task.ChunkIndex), task.Data)

	// Create a map to hold the partitions
	partitions := make(map[int][]KeyValue)

	// Iterate over each key-value pair
	for _, kv := range kvs {
		// Calculate the partition index using the hash of the key
		r := ihash(kv.Key) % w.Master.NumReduce

		// Append the key-value pair to the corresponding partition
		partitions[r] = append(partitions[r], kv)
	}

	// Report the completion of the map task with the collected partitions
	w.Master.ReportMapDone(task.TaskID, partitions)
}

func (w *Worker) doReduce(task Task) {
	partition := w.Master.GetReducePartition(task.ChunkIndex)

	var keys []string
	for k := range partition {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var outLines []string
	for _, k := range keys {
		res := w.reduceFn(k, partition[k])
		outLines = append(outLines, fmt.Sprintf("%s %s", k, res))
	}
	final := strings.Join(outLines, "\n")
	w.Master.ReportReduceDone(task.TaskID, final)
}
