package main

import (
	"os"
	"testing"
)

func TestCommandLineInvokesHelp(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"orgnetsim", "-help"}
	main()
}

func TestCommandLineParse(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"orgnetsim", "parse", "./network.json.tst"}
	main()
}

func TestCommandLineServe(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"orgnetsim", "serve"}
	main()
}
