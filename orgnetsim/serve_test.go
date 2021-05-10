package main

import (
	"os"
	"testing"
)

func TestServeReturnsFalseForHelp(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"orgnetsim", "serve", "-help"}
	success, staticDir, so := serveCommandLineOptions()
	IsFalse(t, success, "-help not returning false")
	AreEqual(t, so.Port, "8080", "Incorrect default port")
	AreEqual(t, staticDir, "", "Default staticDir should be empty")
}

func TestServeReturnsTrue(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"orgnetsim", "serve", "tmpDir"}
	success, staticDir, so := serveCommandLineOptions()
	IsTrue(t, success, "not returning true")
	AreEqual(t, so.Port, "8080", "Incorrect default port")
	AreEqual(t, staticDir, "", "Default staticDir should be empty")
}
