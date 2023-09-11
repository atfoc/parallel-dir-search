package main

import (
	"fmt"
	"os"
	"path"
)

func main() {
	ListDirectoryRecursivelyParallel("../..")
}

func ListDirectoryRecursivelyParallel(baseDir string) {
	dirsToProcess := []string{baseDir}

	for len(dirsToProcess) > 0 {
		dir := dirsToProcess[0]
		dirsToProcess = dirsToProcess[1:]

		dirContents, err := os.ReadDir(dir)
		if err != nil {
			panic(err)
		}

		for _, singleDir := range dirContents {
			singleDirName := path.Join(dir, singleDir.Name())
			fmt.Println(singleDirName)
			if singleDir.IsDir() {
				dirsToProcess = append(dirsToProcess, singleDirName)
			}
		}
	}
}
