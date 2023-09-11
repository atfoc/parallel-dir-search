package main

import (
	"fmt"
	"os"
	"path"
)

func main() {
	ListDirectoryRecursivelyParallel("../..")
}

type UnboundedQueue struct {
	data []string
}

func (q *UnboundedQueue) Pop() (data string, ok bool) {
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
	q.data = append(q.data, name)
}

func ListDirectoryRecursivelyParallel(baseDir string) {
	dirsToProcess := UnboundedQueue{}
	dirsToProcess.Push(baseDir)

	for dir, ok := dirsToProcess.Pop(); ok; dir, ok = dirsToProcess.Pop() {
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
