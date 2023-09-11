package main

import "os"

func main() {
	dirLister := NewDirLister()
	dirLister.List(".")
	dirLister.WaitToFinish()
}

type DirLister struct {
	workers              []Worker
	writeOutWorkerInput  chan os.DirEntry
	writeOutWorkerOutput chan bool
}

type Worker struct {
	inputChanel  chan string
	outputChanel chan os.DirEntry
}

func (w *Worker) GetInputChanel() chan string {
	w.init()
	return w.inputChanel
}

func (w *Worker) init() {
	if w.inputChanel != nil {
		return
	}

	w.inputChanel = make(chan string)
	w.outputChanel = make(chan os.DirEntry)
}

func (w *Worker) GetOutputChanel() chan os.DirEntry {
	w.init()
	return w.outputChanel
}

func (o *DirLister) List(basePath string) {
	for i := 0; i < 5; i++ {
		worker := Worker{}
		o.workers = append(o.workers, worker)
		go dirListerWorker(worker.GetInputChanel(), worker.GetOutputChanel())
	}

	o.writeOutWorkerInput = make(chan os.DirEntry)
	o.writeOutWorkerOutput = make(chan bool)
	go dirWriteOutWorker(o.writeOutWorkerInput, o.writeOutWorkerOutput)

	go o.listImpl(basePath)
}

func dirWriteOutWorker(input chan os.DirEntry, output chan bool) {
}

func (o *DirLister) WaitToFinish() {

}

func (o *DirLister) listImpl(path string) {

}

func NewDirLister() DirLister {
	return DirLister{}
}

func dirListerWorker(chan string, chan os.DirEntry) {

}
