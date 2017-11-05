package srvr

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

var tmpFileCount int
var randVal int64

var dirLock = &sync.RWMutex{}

func GenerateFileName() string {
	tmpFileCount++
	return fmt.Sprintf("tmp%d_%d.json", randVal, tmpFileCount)
}

type TestPersistable struct {
	TimestampHolder
	Data     string `json:"data"`
	Filename string `json:"-"`
}

func (tp *TestPersistable) Filepath() string {
	return tp.Filename
}

func LockfileName(filename string) string {
	return filename + ".lk"
}

func NewTestPersistable() *TestPersistable {
	tmpFileCount++
	if randVal == 0 {
		randVal = time.Now().UnixNano()
	}
	tp := &TestPersistable{
		Data:     "Some information here",
		Filename: GenerateFileName(),
	}
	return tp
}

func TestDelete(t *testing.T) {
	tp := NewTestPersistable()
	filename := tp.Filepath()
	create(filename, t)
	write(filename, tp, t)
	defer remove(filename)

	fd := &FileDetails{
		Filepath: filename,
		DirLock:  dirLock,
	}
	err := fd.Delete()
	AssertSuccess(t, err)
	FileDoesNotExist(t, filename)
}

func TestUpdateFileFailsWhenLkFileExists(t *testing.T) {
	tp := NewTestPersistable()
	filename := tp.Filepath()
	create(filename, t)
	write(filename, tp, t)
	defer remove(filename)
	create(LockfileName(filename), t)
	defer remove(LockfileName(filename))

	tpr := &TestPersistable{
		Data: "Some new data here",
	}
	fd := &FileDetails{
		Filepath: filename,
		DirLock:  dirLock,
	}
	err := fd.Update(tpr)
	IsTrue(t, err != nil, "Update should fail since lk file exists")
}

func TestUpdateFileFailsWhenFileDoesNotExists(t *testing.T) {
	tpr := NewTestPersistable()
	filename := tpr.Filepath()
	fd := &FileDetails{
		Filepath: filename,
		DirLock:  dirLock,
	}
	err := fd.Update(tpr)
	IsTrue(t, err != nil, "Update should fail since target file does not exist")
}

func TestUpdateFileFailsWhenPersistableOutOfDate(t *testing.T) {
	tp := NewTestPersistable()
	filename := tp.Filepath()
	create(filename, t)
	write(filename, tp, t)
	defer remove(filename)

	tpr := &TestPersistable{
		Data: "Some new data here",
	}
	fd := &FileDetails{
		Filepath: filename,
		DirLock:  dirLock,
	}
	err := fd.Update(tpr)
	IsTrue(t, err != nil, "Update should fail since timestamp is out of date")
}

func TestUpdateFileSucceeds(t *testing.T) {
	tp := NewTestPersistable()
	filename := tp.Filepath()
	create(filename, t)
	write(filename, tp, t)
	st, err := stat(filename)
	AssertSuccess(t, err)
	stamp := st.ModTime()
	defer remove(filename)

	tpr := &TestPersistable{
		TimestampHolder: TimestampHolder{
			Stamp: stamp,
		},
		Data: "Some new data here",
	}
	fd := &FileDetails{
		Filepath: filename,
		DirLock:  dirLock,
	}
	err = fd.Update(tpr)
	AssertSuccess(t, err)
	NotEqual(t, stamp, tpr.Stamp, "Update should have changed the timestamp")
}

func TestConcurrentUpdateToFileSucceedsOrFailsWithLockOrStaleErrors(t *testing.T) {
	tp := NewTestPersistable()
	filename := tp.Filepath()
	create(filename, t)
	write(filename, tp, t)
	defer remove(filename)

	hold := make(chan bool)
	fd := &FileDetails{
		Filepath: filename,
		DirLock:  dirLock,
	}
	result := make(chan int)

	for i := 0; i < 100; i++ {
		go func() {
			<-hold
			r := rand.Intn(10)
			time.Sleep(time.Duration(r) * time.Nanosecond)

			tpr := &TestPersistable{}
			fd.Read(tpr)
			rs := tpr.Timestamp()
			tpr.Data = fmt.Sprintf("Data %d", i)
			err := fd.Update(tpr)
			chk := &TestPersistable{}
			e := fd.Read(chk)
			if e != nil {
				fd.Read(chk)
			}
			switch {
			case err == nil && tpr.Timestamp() != rs && tpr.Data == chk.Data:
				result <- 1
			case err != nil && strings.Contains(err.Error(), "Unable to lock file"):
				result <- 2
			case err != nil && strings.Contains(err.Error(), "Stale data"):
				result <- 3
			default:
				result <- 0
			}
		}()
	}
	close(hold)
	updateCount := 0
	staleCount := 0
	lockCount := 0
	for i := 0; i < 100; i++ {
		val := <-result
		IsTrue(t, val > 0, fmt.Sprintf("%d update failed", i))
		switch {
		case val == 1:
			updateCount++
		case val == 2:
			lockCount++
		case val == 3:
			staleCount++
		}
	}
	IsTrue(t, updateCount > 0, "Expected at least 1 update to succeed")
	IsTrue(t, staleCount > 0, "Expected at least 1 update to be see stale data")
	close(result)
}

func TestConcurrentReadFromFileSucceeds(t *testing.T) {
	tp := NewTestPersistable()
	filename := tp.Filepath()
	create(filename, t)
	write(filename, tp, t)
	st, err := stat(filename)
	AssertSuccess(t, err)
	stamp := st.ModTime()
	defer remove(filename)

	hold := make(chan bool)
	fd := &FileDetails{
		Filepath: filename,
		DirLock:  dirLock,
	}
	success := make(chan bool)

	for i := 0; i < 100; i++ {
		go func() {
			<-hold
			r := rand.Intn(10)
			time.Sleep(time.Duration(r) * time.Nanosecond)

			tpr := &TestPersistable{}
			err := fd.Read(tpr)
			success <- err == nil && tpr.Stamp == stamp && tpr.Data == "Some information here"
		}()
	}
	close(hold)
	for i := 0; i < 100; i++ {
		IsTrue(t, <-success, fmt.Sprintf("%d read failed", i))
	}
	close(success)
}

func TestReadFromFileSucceeds(t *testing.T) {
	tp := NewTestPersistable()
	filename := tp.Filepath()
	create(filename, t)
	write(filename, tp, t)
	st, err := stat(filename)
	AssertSuccess(t, err)
	defer remove(filename)

	tpr := &TestPersistable{}
	fd := &FileDetails{
		Filepath: filename,
		DirLock:  dirLock,
	}
	err = fd.Read(tpr)
	AssertSuccess(t, err)
	AreEqual(t, "Some information here", tpr.Data, "Data was not read correctly from file")
	AreEqual(t, st.ModTime(), tpr.Timestamp(), "Timestamp value not updated during read")
}

func TestReadFromFileFailsWhenLkFileExists(t *testing.T) {
	tp := NewTestPersistable()
	filename := tp.Filepath()
	create(filename, t)
	write(filename, tp, t)
	defer remove(filename)
	create(LockfileName(filename), t)
	defer remove(LockfileName(filename))

	tpr := &TestPersistable{}
	fd := &FileDetails{
		Filepath: filename,
		DirLock:  dirLock,
	}
	err := fd.Read(tpr)
	IsTrue(t, err != nil, "Read should fail since lk file exists")
	AreEqual(t, "", tpr.Data, "Data should not be read from file")
}

func TestReadDoesNotReadFileIfSameTimestamp(t *testing.T) {
	tp := NewTestPersistable()
	filename := tp.Filepath()
	create(filename, t)
	write(filename, tp, t)
	st, err := stat(filename)
	AssertSuccess(t, err)
	defer remove(filename)

	tpr := &TestPersistable{}
	tpr.Stamp = st.ModTime()
	fd := &FileDetails{
		Filepath: filename,
		DirLock:  dirLock,
	}
	err = fd.Read(tpr)
	AssertSuccess(t, err)
	AreEqual(t, "", tpr.Data, "Data should not be read from file because Timestamps are same")
}

func TestCreate(t *testing.T) {
	obj := NewTestPersistable()
	filename := obj.Filepath()
	fd := &FileDetails{
		Filepath: filename,
		DirLock:  dirLock,
	}
	err := fd.Create(obj)
	defer remove(filename)
	AssertSuccess(t, err)
	FileExists(t, fd)
}

func TestCreateUpdatesTimeStamp(t *testing.T) {
	obj := NewTestPersistable()
	filename := obj.Filepath()
	fd := &FileDetails{
		Filepath: filename,
		DirLock:  dirLock,
	}
	err := fd.Create(obj)
	defer remove(filename)
	AssertSuccess(t, err)
	FileExists(t, fd)
	st, err := stat(filename)
	AssertSuccess(t, err)
	AreEqual(t, st.ModTime(), obj.Stamp, "Timestamp has not been updated")
}

func TestCreateReturnsErrorWhenFileExists(t *testing.T) {
	obj := NewTestPersistable()
	filename := obj.Filepath()
	fd := &FileDetails{
		Filepath: filename,
		DirLock:  dirLock,
	}
	err := fd.Create(obj)
	defer remove(filename)
	AssertSuccess(t, err)
	FileExists(t, fd)
	err = fd.Create(obj)
	IsTrue(t, err != nil, "CreateNewFile should return error when File exists")
}

func TestCreateReturnsErrorWhenLKFileExists(t *testing.T) {
	obj := NewTestPersistable()
	filename := obj.Filepath()
	fd := &FileDetails{
		Filepath: filename,
		DirLock:  dirLock,
	}

	create(LockfileName(filename), t)
	defer remove(LockfileName(filename))
	err := fd.Create(obj)
	IsTrue(t, err != nil, "CreateNewFile should return error when lock file exists")
	FileDoesNotExist(t, filename)
}

func TestCreateReturnsErrorInvalidFilename(t *testing.T) {
	obj := NewTestPersistable()
	fd := &FileDetails{
		Filepath: "TestCreate///.json",
		DirLock:  dirLock,
	}
	err := fd.Create(obj)
	IsTrue(t, err != nil, "CreateNewFile should return error when invalid filename")
}

func TestCreateReturnsErrorEmptyFilename(t *testing.T) {
	obj := &TestPersistable{}
	fd := &FileDetails{
		DirLock: dirLock,
	}
	err := fd.Create(obj)
	IsTrue(t, err != nil, "CreateNewFile should return error when empty filename")
}

func FileExists(t *testing.T, fu FileUpdater) {
	var path string
	if fu != nil {
		path = fu.Path()
	} else {
		t.Errorf("FileUpdater is nil")
		return
	}
	_, err := stat(path)
	AssertSuccess(t, err)
}

func FileDoesNotExist(t *testing.T, path string) {
	var err error
	//Try twice because file system can return access denied on TestDelete
	for i := 0; i < 10; i++ {
		_, err = stat(path)
		if os.IsNotExist(err) {
			return
		}
		time.Sleep(5 * time.Nanosecond)
	}
	t.Errorf("File should not exist: %v", err)
}

func stat(path string) (os.FileInfo, error) {
	dirLock.RLock()
	defer dirLock.RUnlock()
	return os.Stat(path)
}

func remove(path string) error {
	dirLock.Lock()
	defer dirLock.Unlock()
	return os.Remove(path)
}

func create(path string, t *testing.T) {
	dirLock.Lock()
	defer dirLock.Unlock()
	fl, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		t.Error(err)
	}
	err = fl.Close()
	if err != nil {
		t.Error(err)
	}
}

func write(path string, p Persistable, t *testing.T) {
	fl, err := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		t.Error(err)
	}
	b, err := json.Marshal(p)
	if err != nil {
		t.Error(err)
	}
	_, err = fl.Write(b)
	if err != nil {
		t.Error(err)
	}
	err = fl.Close()
	if err != nil {
		t.Error(err)
	}
}
