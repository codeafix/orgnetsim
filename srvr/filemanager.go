package srvr

import "sync"

//UpdaterRepo contains all FileUpdaters in use in this instance
type UpdaterRepo struct {
	Rootpath string
	Repo     map[string]FileUpdater
	Lock     sync.Mutex
	DirLock  sync.RWMutex
}

//FileManager is a repository for all instances of FileUpdater
type FileManager interface {
	Get(path string) FileUpdater
}

//NewFileManager returns a new instance of FileManager
func NewFileManager(rootpath string) FileManager {
	ur := UpdaterRepo{
		Rootpath: rootpath,
		Repo:     make(map[string]FileUpdater),
	}
	return &ur
}

//Get a FileUpdater for the specifed path
func (ur *UpdaterRepo) Get(path string) FileUpdater {
	ur.Lock.Lock()
	defer ur.Lock.Unlock()
	fu, exists := ur.Repo[path]
	if exists {
		return fu
	}
	fu = &FileDetails{
		Rootpath: ur.Rootpath,
		Filepath: path,
		DirLock:  &ur.DirLock,
	}
	ur.Repo[path] = fu
	return fu
}
