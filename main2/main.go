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

func listDirWorker(input chan []string, waitGroup *sync.WaitGroup) {
	for {
		dirs := <-input

		var dirsToAdd []string
		for _, dir := range dirs {
			dirContents, err := os.ReadDir(dir)
			if err != nil {
				panic(err)
			}

			for _, singleDir := range dirContents {
				singleDirName := path.Join(dir, singleDir.Name())
				fmt.Println(singleDirName)
				if singleDir.IsDir() {
					dirsToAdd = append(dirsToAdd, singleDirName)
					waitGroup.Add(1)
				}
			}
			waitGroup.Done()
		}
		input <- dirsToAdd
	}
}

func ListDirectoryRecursivelyParallel(baseDir string) {
	waitGroup := sync.WaitGroup{}
	numberOfWorkers := 10
	input := make(chan []string, numberOfWorkers)
	for i := 0; i < numberOfWorkers; i++ {
		go listDirWorker(input, &waitGroup)
	}

	input <- []string{baseDir}
	waitGroup.Add(1)
	waitGroup.Wait()
}
