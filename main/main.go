package main

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"
)

func main() {
	ListDirectoryRecursivelyParallel("/Users/pedjat/Documents")
}

type UnboundedQueue struct {
	data []string
	cond *sync.Cond
}

func (q *UnboundedQueue) Pop() string {
	q.init()
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	for !q.hasMore() {
		q.cond.Wait()
	}
	result := q.data[0]
	q.data = q.data[1:]

	return result
}

func (q *UnboundedQueue) hasMore() bool {
	return len(q.data) > 0
}

func (q *UnboundedQueue) Push(name string) {
	q.init()
	q.cond.L.Lock()
	q.data = append(q.data, name)
	q.cond.Broadcast()
	q.cond.L.Unlock()
}

func (q *UnboundedQueue) init() {
	if q.cond != nil {
		return
	}

	q.cond = sync.NewCond(&sync.Mutex{})
}

func ListDirectoryRecursivelyParallel(baseDir string) {
	dirsToProcess := UnboundedQueue{}

	for i := 0; i < 1000; i++ {
		go listDirWorker(&dirsToProcess)
	}

	dirsToProcess.Push(baseDir)
	time.Sleep(time.Second * 10)
}

func listDirWorker(dirsToProcess *UnboundedQueue) {
	for {
		dir := dirsToProcess.Pop()

		dirContents, err := os.ReadDir(dir)
		if err != nil {
			panic(err)
		}

		for _, singleDir := range dirContents {
			singleDirName := path.Join(dir, singleDir.Name())
			fmt.Println(singleDirName)
			if singleDir.IsDir() {
				dirsToProcess.Push(singleDirName)
			}
		}
	}
}
