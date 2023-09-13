package main

import (
	"fmt"
	"log"
	"main/unboundedqueue"
	"os"
	"path"
	"sync"
)

var errorLogger = log.New(os.Stderr, "", 0)

func main() {

	ListDirectoryRecursivelyParallel("/Users/pedjat/Documents")
}

func listDirWorker(dirsToProcess *unboundedqueue.UnboundedQueue, waitGroup *sync.WaitGroup, workerIndex int) {
	for {
		dir := dirsToProcess.Pop()
		ObserveDirRemovedFromQueue()
		dirContents, err := os.ReadDir(dir)
		if err != nil {
			errorLogger.Printf("Failed reading dir %s with error %s\n", dir, err.Error())
			continue
		}

		for _, singleDir := range dirContents {
			singleDirName := path.Join(dir, singleDir.Name())
			fmt.Println(singleDirName)
			if singleDir.IsDir() {
				ObserveAddDirToQueue()
				dirsToProcess.Push(singleDirName)
				waitGroup.Add(1)
			}
		}

		ObserveDirProcessed(workerIndex)
		waitGroup.Done()
	}
}

func ListDirectoryRecursivelyParallel(baseDir string) {
	dirsToProcess := unboundedqueue.New()

	waitGroup := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		go listDirWorker(dirsToProcess, &waitGroup, i)
	}

	dirsToProcess.Push(baseDir)
	waitGroup.Add(1)
	waitGroup.Wait()
}
