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

func ListDirectoryRecursivelyParallel(baseDir string) {
	dirsToProcess := unboundedqueue.UnboundedQueue{}

	waitGroup := sync.WaitGroup{}
	for i := 0; i < 20; i++ {
		go listDirWorker(&dirsToProcess, &waitGroup)
	}

	dirsToProcess.Push(baseDir)
	waitGroup.Add(1)
	waitGroup.Wait()
}
