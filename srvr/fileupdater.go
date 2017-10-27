package srvr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

//FileDetails contains information about a file that will be managed by the FileUpdater
type FileDetails struct {
	Filepath string
	Lock     sync.RWMutex
}

//FileUpdater manages a file and provides Read and Update methods that allow the rest of
//the package to read and update a file.
//FileUpdater contains a lock to prevent multiple goroutines from trying to write at the
//same time and also prevents Reads from occuring whilst an update is being performed.
//The Update function checks the last modified time of the file before it updates it,
//if it has been modified since the last copy was read, then Update will return an
//error, and the Update method will attempt to re-read the contents of the file into
//the object supplied.
type FileUpdater interface {
	Read(obj Persistable) error
	Update(obj Persistable) error
	Delete() error
	Path() string
}

//CreateNewFile creates a file and saves the supplied object into the file identified by path
//It returns the FileUpdater that should be used to access the file
func CreateNewFile(path string, obj Persistable) (FileUpdater, error) {
	fd := FileDetails{
		Filepath: path,
	}
	lkPath, err := fd.createLockFile()
	if err != nil {
		return nil, err
	}
	defer os.Remove(lkPath)
	exists, err := fd.fileExists(fd.Filepath)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("File exists")
	}
	return &fd, fd.writeFile(obj, true)
}

//Read the file and unmarshal it into the supplied object. If the file hasn't changed
//since it was last read this function will do nothing. If the file is re-read then
//the Timestamp on the supplied object is updated
func (fd *FileDetails) Read(obj Persistable) error {
	fd.Lock.RLock()
	defer fd.Lock.RUnlock()
	exists, err := fd.fileExists(fd.lockpath())
	if err != nil || exists {
		return fmt.Errorf("File is locked by another process")
	}
	s, err := os.Stat(fd.Filepath)
	if err != nil {
		return err
	}
	if obj.Timestamp() == s.ModTime() {
		return nil
	}
	b, err := ioutil.ReadFile(fd.Filepath)
	if err != nil {
		return err
	}
	obj.UpdateTimestamp(s.ModTime())
	return json.Unmarshal(b, obj)
}

//Update the contents of the file with the supplied persistable object. If the write
//to the file is successful, then the Timestamp of the persistable is updated too.
func (fd *FileDetails) Update(obj Persistable) error {
	fd.Lock.Lock()
	defer fd.Lock.Unlock()
	lkPath, err := fd.createLockFile()
	if err != nil {
		return err
	}
	defer os.Remove(lkPath)
	s, err := os.Stat(fd.Filepath)
	if err != nil {
		return err
	}
	if obj.Timestamp() != s.ModTime() {
		return fmt.Errorf("Stale data")
	}
	return fd.writeFile(obj, false)
}

//Delete the file on the file system
func (fd *FileDetails) Delete() error {
	return os.Remove(fd.Filepath)
}

//Path to the file on the file system
func (fd *FileDetails) Path() string {
	return fd.Filepath
}

//Actually writes data into a file
func (fd *FileDetails) writeFile(obj Persistable, create bool) error {
	var fo *os.File
	var err error
	if create {
		fo, err = os.OpenFile(fd.Filepath, os.O_CREATE|os.O_EXCL, 0644)
	} else {
		fo, err = os.OpenFile(fd.Filepath, os.O_TRUNC|os.O_WRONLY, 0644)
	}
	defer func() {
		s, err := os.Stat(fd.Path())
		if err != nil {
			return
		}
		obj.UpdateTimestamp(s.ModTime())
	}()
	defer fo.Close()
	if err != nil {
		return err
	}
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = fo.Write(b)
	return err
}

//Creates a lock file to avoid multiple instances clobbering each other
func (fd *FileDetails) createLockFile() (string, error) {
	lkPath := fd.lockpath()
	lk, err := os.OpenFile(lkPath, os.O_CREATE|os.O_EXCL, 0644)
	err2 := lk.Close()
	if err != nil || err2 != nil {
		return lkPath, fmt.Errorf("Unable to lock file '%s'", fd.Filepath)
	}
	return lkPath, nil
}

//Returns a true if the file currently exists
func (fd *FileDetails) fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

//Returns the path to the lock file for this file
func (fd *FileDetails) lockpath() string {
	return fd.Filepath + ".lk"
}
