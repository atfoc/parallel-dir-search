package main

import (
	"fmt"
	"main/unboundedqueue"
	"os"
	"path"
	"sync"
)

func main() {
	ListDirectoryRecursivelyParallel("/Users/pedjat/Documents")
}

func listDirWorker(dirsToProcess *unboundedqueue.UnboundedQueue, waitGroup *sync.WaitGroup) {
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
				waitGroup.Add(1)
			}
		}
		waitGroup.Done()
	}
}

func ListDirectoryRecursivelyParallel(baseDir string) {
	dirsToProcess := unboundedqueue.New()

	waitGroup := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		go listDirWorker(dirsToProcess, &waitGroup)
	}

	dirsToProcess.Push(baseDir)
	waitGroup.Add(1)
	waitGroup.Wait()
}
