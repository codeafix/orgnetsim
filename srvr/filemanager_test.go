package srvr

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"testing"
	"time"
)

func TestGetsFileUpdaterForWindowsPath(t *testing.T) {
	root := "c:\\go_work\\src\\github.com\\randomtask\\orgnetsim\\srvr\\"
	fm := NewFileManager(root)
	file := "dir\\somefile.json"
	path := filepath.Join(root, file)
	fu := fm.Get(file)
	AreEqual(t, path, fu.Path(), "Paths do not match")
	fu2 := fm.Get(file)
	AreEqual(t, fu, fu2, "FileUpdaters do not match")
}

func TestGetsFileUpdaterForNixPath(t *testing.T) {
	root := "/home/user/go_work/src/github.com/randomtask/orgnetsim/srvr"
	fm := NewFileManager(root)
	file := "dir/somefile.json"
	path := filepath.Join(root, file)
	fu := fm.Get(file)
	AreEqual(t, path, fu.Path(), "Paths do not match")
	fu2 := fm.Get(file)
	AreEqual(t, fu, fu2, "FileUpdaters do not match")
}

func TestGetsFileUpdaterForUriPath(t *testing.T) {
	root := "file:///home/user/go_work/src/github.com/randomtask/orgnetsim/srvr/"
	fm := NewFileManager(root)
	file := "dir/somefile.json"
	path := filepath.Join(root, file)
	fu := fm.Get(file)
	AreEqual(t, path, fu.Path(), "Paths do not match")
	fu2 := fm.Get(file)
	AreEqual(t, fu, fu2, "FileUpdaters do not match")
}

func TestConcurrentGetFromFileManager(t *testing.T) {
	root := "file:///home/user/go_work/src/github.com/randomtask/orgnetsim/srvr/"
	fm := NewFileManager(root)

	hold := make(chan bool)
	g1success := make(chan bool)
	g2success := make(chan bool)

	for i := 0; i < 100; i++ {
		go func(count int) {
			<-hold
			file := fmt.Sprintf("dir/file%d.json", count)
			path := filepath.Join(root, file)
			fu := fm.Get(file)
			fu2 := fm.Get(file)
			g1success <- fu.Path() == path && fu == fu2
		}(i)
	}

	for i := 0; i < 100; i++ {
		go func(count int) {
			<-hold
			file := fmt.Sprintf("dir/file%d.json", count)
			path := filepath.Join(root, file)
			r := rand.Intn(10)
			time.Sleep(time.Duration(r) * time.Nanosecond)
			fu := fm.Get(file)
			fu2 := fm.Get(file)
			g2success <- fu.Path() == path && fu == fu2
		}(i)
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
