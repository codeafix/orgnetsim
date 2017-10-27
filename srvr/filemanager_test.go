package srvr

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestGetsFileUpdaterForWindowsPath(t *testing.T) {
	fm := NewFileManager()
	path := "c:\\go_work\\src\\github.com\\randomtask\\orgnetsim\\srvr\\somefile.json"
	fu := fm.Get(path)
	AreEqual(t, path, fu.Path(), "Paths do not match")
	fu2 := fm.Get(path)
	AreEqual(t, fu, fu2, "FileUpdaters do not match")
}

func TestGetsFileUpdaterForNixPath(t *testing.T) {
	fm := NewFileManager()
	path := "/home/user/go_work/src/github.com/randomtask/orgnetsim/srvr/somefile.json"
	fu := fm.Get(path)
	AreEqual(t, path, fu.Path(), "Paths do not match")
	fu2 := fm.Get(path)
	AreEqual(t, fu, fu2, "FileUpdaters do not match")
}

func TestGetsFileUpdaterForUriPath(t *testing.T) {
	fm := NewFileManager()
	path := "file:///home/user/go_work/src/github.com/randomtask/orgnetsim/srvr/somefile.json"
	fu := fm.Get(path)
	AreEqual(t, path, fu.Path(), "Paths do not match")
	fu2 := fm.Get(path)
	AreEqual(t, fu, fu2, "FileUpdaters do not match")
}

func TestConcurrentGetFromFileManager(t *testing.T) {
	fm := NewFileManager()
	rootpath := "file:///home/user/go_work/src/github.com/randomtask/orgnetsim/srvr/"

	hold := make(chan bool)
	g1success := make(chan bool)
	g2success := make(chan bool)

	for i := 0; i < 100; i++ {
		go func() {
			<-hold
			path := fmt.Sprintf("%sfile%d.json", rootpath, i)
			fu := fm.Get(path)
			fu2 := fm.Get(path)
			g1success <- fu.Path() == path && fu == fu2
		}()
	}

	for i := 0; i < 100; i++ {
		go func() {
			<-hold
			path := fmt.Sprintf("%sfile%d.json", rootpath, i)
			r := rand.Intn(10)
			time.Sleep(time.Duration(r) * time.Nanosecond)
			fu := fm.Get(path)
			fu2 := fm.Get(path)
			g2success <- fu.Path() == path && fu == fu2
		}()
	}

	close(hold)

	for i := 200; i > 0; {
		select {
		case success := <-g1success:
			if !success {
				t.Errorf("Failure to get FileUpdater")
			}
			i--
		case success := <-g2success:
			if !success {
				t.Errorf("Failure to get FileUpdater")
			}
			i--
		}
	}

	close(g1success)
	close(g2success)
}
