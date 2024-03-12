package global

import "sync"

type FileInfo struct {
	Name      string
	FullPath  string
	StorePath string
	Task      string
}

var File = sync.Map{} // Hash FileInfo

var Task = sync.Map{} //taskName Get []Hashes
