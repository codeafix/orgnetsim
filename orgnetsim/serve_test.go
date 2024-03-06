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

func TestServeWithoutPortReturnsFalse(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"orgnetsim", "serve", "tmpDir", "-p"}
	success, _, _ := serveCommandLineOptions()
	IsFalse(t, success, "not returning false")
}

func TestServeWithPortReturnsTrue(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"orgnetsim", "serve", "tmpDir", "-p", "8081"}
	success, _, so := serveCommandLineOptions()
	IsTrue(t, success, "not returning true")
	AreEqual(t, so.Port, "8081", "Incorrect port")
}

func TestServeWithStaticDirReturnsTrue(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"orgnetsim", "serve", "tmpDir", "-s", "web"}
	success, staticDir, _ := serveCommandLineOptions()
	IsTrue(t, success, "not returning true")
	AreEqual(t, staticDir, "web", "Incorrect staticDir")
}

func TestServeWithPortAndStaticDirReturnsTrue(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"orgnetsim", "serve", "tmpDir", "-p", "8081", "-s", "web"}
	success, staticDir, so := serveCommandLineOptions()
	IsTrue(t, success, "not returning true")
	AreEqual(t, so.Port, "8081", "Incorrect port")
	AreEqual(t, staticDir, "web", "Incorrect staticDir")
}

func TestServeWithoutStaticPortReturnsFalse(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"orgnetsim", "serve", "tmpDir", "-s"}
	success, _, _ := serveCommandLineOptions()
	IsFalse(t, success, "not returning false")
}

func TestServeReturnsFalseWithUnrecognisedOptions(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"orgnetsim", "serve", "tmpDir", "-x"}
	success, _, _ := serveCommandLineOptions()
	IsFalse(t, success, "not returning false")
}
