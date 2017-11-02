package srvr

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
)

var tmpFileCount int
var randVal int64

//Used to ensure tests run sequentially
//This is needed because many of the tests are creating temporary files in a directory
//and running them concurrently can result in false failures if one test is creating
//or removing a file at the same time another test.
var wait = make(chan bool, 1)

func WaitForTurn() {
	wait <- true
}

func Next() {
	<-wait
}

type TestPersistable struct {
	TimestampHolder
	Data string `json:"data"`
}

func GenerateFileName(t *testing.T) string {
	tmpFileCount++
	if randVal == 0 {
		randVal = time.Now().UnixNano()
	}
	return fmt.Sprintf("tmp%d_%d.json", randVal, tmpFileCount)
}

func LockfileName(filename string) string {
	return filename + ".lk"
}

func TestDelete(t *testing.T) {
	WaitForTurn()
	filename := GenerateFileName(t)
	tp := &TestPersistable{
		Data: "Some information here",
	}
	fl, _ := os.OpenFile(filename, os.O_CREATE|os.O_EXCL, 0644)
	b, err := json.Marshal(tp)
	fl.Write(b)
	fl.Close()
	defer os.Remove(filename)

	fd := &FileDetails{
		Filepath: filename,
	}
	err = fd.Delete()
	AssertSuccess(t, err)
	FileDoesNotExist(t, filename)
	Next()
}

func TestUpdateFileFailsWhenLkFileExists(t *testing.T) {
	WaitForTurn()
	filename := GenerateFileName(t)
	tp := &TestPersistable{
		Data: "Some information here",
	}
	fl, _ := os.OpenFile(filename, os.O_CREATE|os.O_EXCL, 0644)
	b, err := json.Marshal(tp)
	fl.Write(b)
	fl.Close()
	defer os.Remove(filename)
	lk, _ := os.OpenFile(LockfileName(filename), os.O_CREATE|os.O_EXCL, 0644)
	lk.Close()
	defer os.Remove(LockfileName(filename))

	tpr := &TestPersistable{
		Data: "Some new data here",
	}
	fd := &FileDetails{
		Filepath: filename,
	}
	err = fd.Update(tpr)
	IsTrue(t, err != nil, "Update should fail since lk file exists")
	Next()
}

func TestUpdateFileFailsWhenFileDoesNotExists(t *testing.T) {
	WaitForTurn()
	filename := GenerateFileName(t)
	tpr := &TestPersistable{
		Data: "Some new data here",
	}
	fd := &FileDetails{
		Filepath: filename,
	}
	err := fd.Update(tpr)
	IsTrue(t, err != nil, "Update should fail since target file does not exist")
	Next()
}

func TestUpdateFileFailsWhenPersistableOutOfDate(t *testing.T) {
	WaitForTurn()
	filename := GenerateFileName(t)
	tp := &TestPersistable{
		Data: "Some information here",
	}
	fl, _ := os.OpenFile(filename, os.O_CREATE|os.O_EXCL, 0644)
	b, err := json.Marshal(tp)
	fl.Write(b)
	fl.Close()
	defer os.Remove(filename)

	tpr := &TestPersistable{
		Data: "Some new data here",
	}
	fd := &FileDetails{
		Filepath: filename,
	}
	err = fd.Update(tpr)
	IsTrue(t, err != nil, "Update should fail since timestamp is out of date")
	Next()
}

func TestUpdateFileSucceeds(t *testing.T) {
	WaitForTurn()
	filename := GenerateFileName(t)
	tp := &TestPersistable{
		Data: "Some information here",
	}
	fl, _ := os.OpenFile(filename, os.O_CREATE|os.O_EXCL, 0644)
	b, err := json.Marshal(tp)
	fl.Write(b)
	fl.Close()
	st, _ := os.Stat(filename)
	stamp := st.ModTime()
	defer os.Remove(filename)

	tpr := &TestPersistable{
		TimestampHolder: TimestampHolder{
			Stamp: stamp,
		},
		Data: "Some new data here",
	}
	fd := &FileDetails{
		Filepath: filename,
	}
	err = fd.Update(tpr)
	AssertSuccess(t, err)
	IsTrue(t, tpr.Stamp != stamp, "Update should have changed the timestamp")
	Next()
}

func TestConcurrentUpdateToFileSucceedsOrFailsWithLockOrStaleErrors(t *testing.T) {
	WaitForTurn()
	filename := GenerateFileName(t)
	tp := &TestPersistable{
		Data: "Some information here",
	}
	fl, _ := os.OpenFile(filename, os.O_CREATE|os.O_EXCL, 0644)
	b, _ := json.Marshal(tp)
	fl.Write(b)
	fl.Close()
	defer os.Remove(filename)

	hold := make(chan bool)
	fd := &FileDetails{
		Filepath: filename,
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
	Next()
}

func TestConcurrentReadFromFileSucceeds(t *testing.T) {
	WaitForTurn()
	filename := GenerateFileName(t)
	tp := &TestPersistable{
		Data: "Some information here",
	}
	fl, _ := os.OpenFile(filename, os.O_CREATE|os.O_EXCL, 0644)
	b, _ := json.Marshal(tp)
	fl.Write(b)
	fl.Close()
	st, _ := os.Stat(filename)
	stamp := st.ModTime()
	defer os.Remove(filename)

	hold := make(chan bool)
	fd := &FileDetails{
		Filepath: filename,
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
	Next()
}

func TestReadFromFileSucceeds(t *testing.T) {
	WaitForTurn()
	filename := GenerateFileName(t)
	tp := &TestPersistable{
		Data: "Some information here",
	}
	fl, _ := os.OpenFile(filename, os.O_CREATE|os.O_EXCL, 0644)
	b, err := json.Marshal(tp)
	fl.Write(b)
	fl.Close()
	st, _ := os.Stat(filename)
	defer os.Remove(filename)

	tpr := &TestPersistable{}
	fd := &FileDetails{
		Filepath: filename,
	}
	err = fd.Read(tpr)
	AssertSuccess(t, err)
	AreEqual(t, "Some information here", tpr.Data, "Data was not read correctly from file")
	AreEqual(t, st.ModTime(), tpr.Timestamp(), "Timestamp value not updated during read")
	Next()
}

func TestReadFromFileFailsWhenLkFileExists(t *testing.T) {
	WaitForTurn()
	filename := GenerateFileName(t)
	tp := &TestPersistable{
		Data: "Some information here",
	}
	fl, _ := os.OpenFile(filename, os.O_CREATE|os.O_EXCL, 0644)
	b, err := json.Marshal(tp)
	fl.Write(b)
	fl.Close()
	defer os.Remove(filename)
	lk, _ := os.OpenFile(LockfileName(filename), os.O_CREATE|os.O_EXCL, 0644)
	lk.Close()
	defer os.Remove(LockfileName(filename))

	tpr := &TestPersistable{}
	fd := &FileDetails{
		Filepath: filename,
	}
	err = fd.Read(tpr)
	IsTrue(t, err != nil, "Read should fail since lk file exists")
	AreEqual(t, "", tpr.Data, "Data should not be read from file")
	Next()
}

func TestReadDoesNotReadFileIfSameTimestamp(t *testing.T) {
	WaitForTurn()
	filename := GenerateFileName(t)
	tp := &TestPersistable{
		Data: "Some information here",
	}
	fl, _ := os.OpenFile(filename, os.O_CREATE|os.O_EXCL, 0644)
	b, err := json.Marshal(tp)
	fl.Write(b)
	fl.Close()
	st, _ := os.Stat(filename)
	defer os.Remove(filename)

	tpr := &TestPersistable{}
	tpr.Stamp = st.ModTime()
	fd := &FileDetails{
		Filepath: filename,
	}
	err = fd.Read(tpr)
	AssertSuccess(t, err)
	AreEqual(t, "", tpr.Data, "Data should not be read from file because Timestamps are same")
	Next()
}

func TestCreate(t *testing.T) {
	WaitForTurn()
	filename := GenerateFileName(t)
	obj := &TimestampHolder{}
	fu, err := CreateNewFile(filename, obj)
	defer os.Remove(filename)
	AssertSuccess(t, err)
	FileExists(t, fu)
	Next()
}

func TestCreateUpdatesTimeStamp(t *testing.T) {
	WaitForTurn()
	filename := GenerateFileName(t)
	obj := &TimestampHolder{}
	fu, err := CreateNewFile(filename, obj)
	defer os.Remove(filename)
	AssertSuccess(t, err)
	FileExists(t, fu)
	st, err := os.Stat(filename)
	AreEqual(t, st.ModTime(), obj.Stamp, "Timestamp has not been updated")
	Next()
}

func TestCreateReturnsErrorWhenFileExists(t *testing.T) {
	WaitForTurn()
	filename := GenerateFileName(t)
	obj := &TimestampHolder{}
	fu, err := CreateNewFile(filename, obj)
	defer os.Remove(filename)
	AssertSuccess(t, err)
	FileExists(t, fu)
	fu, err = CreateNewFile(filename, obj)
	IsTrue(t, fu == nil, "CreateNewFile returns a FileUpdater but should have returned nil when file existed")
	IsTrue(t, err != nil, "CreateNewFile should return error when File exists")
	Next()
}

func TestCreateReturnsErrorWhenLKFileExists(t *testing.T) {
	WaitForTurn()
	filename := GenerateFileName(t)
	obj := &TimestampHolder{}
	lk, _ := os.OpenFile(LockfileName(filename), os.O_CREATE|os.O_EXCL, 0644)
	lk.Close()
	defer os.Remove(LockfileName(filename))
	fu, err := CreateNewFile(filename, obj)
	IsTrue(t, fu == nil, "CreateNewFile returns a FileUpdater but should have returned nil when lock file exists")
	IsTrue(t, err != nil, "CreateNewFile should return error when lock file exists")
	FileDoesNotExist(t, filename)
	Next()
}

func TestCreateReturnsErrorInvalidFilename(t *testing.T) {
	obj := &TimestampHolder{}
	fu, err := CreateNewFile("TestCreate??\\.json", obj)
	IsTrue(t, fu == nil, "CreateNewFile returns a FileUpdater but should have returned nil when invalid filename")
	IsTrue(t, err != nil, "CreateNewFile should return error when invalid filename")
}

func TestCreateReturnsErrorEmptyFilename(t *testing.T) {
	obj := &TimestampHolder{}
	_, err := CreateNewFile("", obj)
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
	_, err := os.Stat(path)
	AssertSuccess(t, err)
}

func FileDoesNotExist(t *testing.T, path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return
	}
	t.Errorf("File should not exist: %v", err)
}