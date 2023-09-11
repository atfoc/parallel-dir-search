package main

func main() {
	dirLister := NewDirLister()
	dirLister.List(".")
	dirLister.WaitToFinish()
}

type DirLister struct {
}

func (o *DirLister) List(basePath string) {

}

func (o *DirLister) WaitToFinish() {

}

func NewDirLister() DirLister {
	return DirLister{}
}
