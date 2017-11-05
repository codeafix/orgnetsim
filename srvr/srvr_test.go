package srvr

import (
	"encoding/json"
	"strings"
	"testing"
)

func IsFalse(t *testing.T, condition bool, msg string) {
	if condition {
		t.Error(msg)
	}
}

func IsTrue(t *testing.T, condition bool, msg string) {
	if !condition {
		t.Error(msg)
	}
}

func AreEqual(t *testing.T, expected interface{}, actual interface{}, msg string) {
	if expected != actual {
		t.Errorf("%s Expected = '%v' Actual = '%v'", msg, expected, actual)
	}
}

func Contains(t *testing.T, expected string, actual string, msg string) {
	if !strings.Contains(actual, expected) {
		t.Errorf("%s Expected = '%v' Actual = '%v'", msg, expected, actual)
	}
}

func NotEqual(t *testing.T, expected interface{}, actual interface{}, msg string) {
	if expected == actual {
		t.Errorf("%s Expected = '%v' Actual = '%v'", msg, expected, actual)
	}
}

func AssertSuccess(t *testing.T, err error) {
	if err != nil {
		t.Errorf(err.Error())
	}
}

func NewTestFileManager(tfu *TestFileUpdater) *TestFileManager {
	return &TestFileManager{
		map[string]FileUpdater{
			tfu.Path(): tfu,
		},
		nil,
	}
}

type TestFileManager struct {
	FileUpdaters map[string]FileUpdater
	Default      FileUpdater
}

func (fm *TestFileManager) Get(path string) FileUpdater {
	fu, exists := fm.FileUpdaters[path]
	if !exists {
		return fm.Default
	}
	return fu
}

func (fm *TestFileManager) Add(path string, tfu *TestFileUpdater) {
	fm.FileUpdaters[path] = tfu
}

func (fm *TestFileManager) SetDefault(tfu *TestFileUpdater) {
	fm.Default = tfu
}

type TestFileUpdater struct {
	Obj          Persistable
	ReadErr      error
	UpdateErr    error
	CreateErr    error
	DeleteErr    error
	Filepath     string
	DeleteCalled bool
}

func (fu *TestFileUpdater) Create(obj Persistable) error {
	if fu.CreateErr != nil {
		return fu.CreateErr
	}
	fu.Obj = obj
	return nil
}

func (fu *TestFileUpdater) Read(obj Persistable) error {
	if fu.ReadErr != nil {
		return fu.ReadErr
	}
	js, _ := json.Marshal(fu.Obj)
	json.Unmarshal(js, obj)
	return nil
}

func (fu *TestFileUpdater) Update(obj Persistable) error {
	if fu.UpdateErr != nil {
		return fu.UpdateErr
	}
	fu.Obj = obj
	return nil
}

func (fu *TestFileUpdater) Delete() error {
	fu.DeleteCalled = true
	if fu.DeleteErr != nil {
		return fu.DeleteErr
	}
	return nil
}

func (fu *TestFileUpdater) Path() string {
	return fu.Filepath
}
