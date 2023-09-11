package main

import (
	"fmt"
	"os"
	"path"
	"sync"
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

func (q *UnboundedQueue) HasMore() bool {
	q.init()
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	return q.hasMore()
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

func (q *UnboundedQueue) PushAll(names []string) {
	q.init()
	q.cond.L.Lock()
	for _, el := range names {
		q.data = append(q.data, el)
	}
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

	waitGroup := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		go listDirWorker(&dirsToProcess, &waitGroup)
	}

	dirsToProcess.Push(baseDir)
	waitGroup.Add(1)
	waitGroup.Wait()
}

func listDirWorker(dirsToProcess *UnboundedQueue, waitGroup *sync.WaitGroup) {
	for {
		dir := dirsToProcess.Pop()

		dirContents, err := os.ReadDir(dir)
		if err != nil {
			panic(err)
		}

		var dirs []string
		for _, singleDir := range dirContents {
			singleDirName := path.Join(dir, singleDir.Name())
			fmt.Println(singleDirName)
			if singleDir.IsDir() {
				dirs = append(dirs, singleDirName)
			}
		}
		dirsToProcess.PushAll(dirs)
		waitGroup.Add(len(dirs))
		waitGroup.Done()
	}
}
