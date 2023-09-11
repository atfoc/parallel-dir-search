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
	lock sync.Mutex
}

func (q *UnboundedQueue) Pop() (data string, ok bool) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if !q.hasMore() {
		return "", false
	}

	result := q.data[0]
	q.data = q.data[1:]
	return result, true
}

func (q *UnboundedQueue) hasMore() bool {
	return len(q.data) > 0
}

func (q *UnboundedQueue) Push(name string) {
	q.lock.Lock()
	q.data = append(q.data, name)
	q.lock.Unlock()
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
		dir, ok := dirsToProcess.Pop()
		if !ok {
			continue
		}

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
