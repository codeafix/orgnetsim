package srvr

import (
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

func NewTestFileManager(tfu *TestFileUpdater) FileManager {
	return &TestFileManager{
		tfu,
	}
}

type TestFileManager struct {
	FileUpdater FileUpdater
}

func (fm *TestFileManager) Get(path string) FileUpdater {
	return fm.FileUpdater
}

type TestFileUpdater struct {
	Obj       Persistable
	ReadErr   error
	UpdateErr error
	CreateErr error
	DeleteErr error
}

func (fu *TestFileUpdater) Create(obj Persistable) error {
	fu.Obj = obj
	return fu.CreateErr
}

func (fu *TestFileUpdater) Read(obj Persistable) error {
	obj = fu.Obj
	return fu.ReadErr
}

func (fu *TestFileUpdater) Update(obj Persistable) error {
	fu.Obj = obj
	return fu.UpdateErr
}

func (fu *TestFileUpdater) Delete() error {
	return fu.DeleteErr
}

func (fu *TestFileUpdater) Path() string {
	return ""
}
