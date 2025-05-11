package srvr

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// FileDetails contains information about a file that will be managed by the FileUpdater
type FileDetails struct {
	Rootpath string
	Filepath string
	Lock     sync.RWMutex
	DirLock  *sync.RWMutex
}

// FileUpdater manages a file and provides Read and Update methods that allow the rest of
// the package to read and update a file.
// FileUpdater contains a lock to prevent multiple goroutines from trying to write at the
// same time and also prevents Reads from occuring whilst an update is being performed.
// The Update function checks the last modified time of the file before it updates it,
// if it has been modified since the last copy was read, then Update will return an
// error, and the Update method will attempt to re-read the contents of the file into
// the object supplied.
type FileUpdater interface {
	Create(obj Persistable) error
	Read(obj Persistable) error
	Update(obj Persistable) error
	Delete() error
	Path() string
}

// Create creates a file and saves the supplied object into the file identified by path
// It returns the FileUpdater that should be used to access the file
func (fd *FileDetails) Create(obj Persistable) error {
	lkPath, err := fd.createLockFile()
	if err != nil {
		return err
	}
	defer fd.remove(lkPath)
	exists, err := fd.fileExists(fd.Path())
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("file exists")
	}
	err = fd.createFile(fd.Path())
	if err != nil {
		return err
	}
	return fd.writeFile(obj)
}

// Read the file and unmarshal it into the supplied object. If the file hasn't changed
// since it was last read this function will do nothing. If the file is re-read then
// the Timestamp on the supplied object is updated
func (fd *FileDetails) Read(obj Persistable) error {
	fd.Lock.RLock()
	defer fd.Lock.RUnlock()
	exists, err := fd.fileExists(fd.lockpath())
	if err != nil || exists {
		return fmt.Errorf("file is locked by another process")
	}
	s, err := fd.stat(fd.Path())
	if err != nil {
		return err
	}
	if obj.Timestamp() == s.ModTime() {
		return nil
	}
	b, err := os.ReadFile(fd.Path())
	if err != nil {
		return err
	}
	obj.UpdateTimestamp(s.ModTime())
	return json.Unmarshal(b, obj)
}

// Update the contents of the file with the supplied persistable object. If the write
// to the file is successful, then the Timestamp of the persistable is updated too.
func (fd *FileDetails) Update(obj Persistable) error {
	fd.Lock.Lock()
	defer fd.Lock.Unlock()
	lkPath, err := fd.createLockFile()
	if err != nil {
		return err
	}
	defer fd.remove(lkPath)
	s, err := fd.stat(fd.Path())
	if err != nil {
		return err
	}
	if obj.Timestamp() != s.ModTime() {
		return fmt.Errorf("stale data")
	}
	return fd.writeFile(obj)
}

// Delete the file on the file system
func (fd *FileDetails) Delete() error {
	return fd.remove(fd.Path())
}

// Path to the file on the file system
func (fd *FileDetails) Path() string {
	return filepath.Join(fd.Rootpath, fd.Filepath)
}

// Actually writes data into a file
func (fd *FileDetails) writeFile(obj Persistable) error {
	fo, err := os.OpenFile(fd.Path(), os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fo.Close()
		return err
	}
	b, err := json.Marshal(obj)
	if err != nil {
		fo.Close()
		return err
	}
	_, err = fo.Write(b)
	if err != nil {
		fo.Close()
		return err
	}
	err = fo.Close()
	if err != nil {
		return err
	}
	s, err := fd.stat(fd.Path())
	if err != nil {
		return err
	}
	obj.UpdateTimestamp(s.ModTime())
	return nil
}

// Creates a lock file to avoid multiple instances clobbering each other
func (fd *FileDetails) createLockFile() (string, error) {
	lkPath := fd.lockpath()
	err := fd.createFile(lkPath)
	if err != nil {
		return lkPath, fmt.Errorf("unable to lock file '%s'", fd.Path())
	}
	return lkPath, nil
}

// Returns a true if the file currently exists
func (fd *FileDetails) fileExists(path string) (bool, error) {
	_, err := fd.stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// Returns the path to the lock file for this file
func (fd *FileDetails) lockpath() string {
	return fd.Path() + ".lk"
}

// Removes the specified file
func (fd *FileDetails) remove(path string) error {
	fd.DirLock.Lock()
	defer fd.DirLock.Unlock()
	return os.Remove(path)
}

// stat the given file path
func (fd *FileDetails) stat(path string) (os.FileInfo, error) {
	fd.DirLock.RLock()
	defer fd.DirLock.RUnlock()
	return os.Stat(path)
}

// createFile the given file path
func (fd *FileDetails) createFile(path string) error {
	fd.DirLock.Lock()
	defer fd.DirLock.Unlock()
	fl, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	return fl.Close()
}
