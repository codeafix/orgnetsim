package main

import (
	"os"
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
		t.Errorf("%s Expected = %v Actual = %v", msg, expected, actual)
	}
}

func NotEqual(t *testing.T, expected interface{}, actual interface{}, msg string) {
	if expected == actual {
		t.Errorf("%s Expected = %v Actual = %v", msg, expected, actual)
	}
}

func AssertSuccess(t *testing.T, err error) {
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestParseReturnsFalseForHelp(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"orgnetsim", "parse", "-help"}
	success, _, _, _ := parseCommandLineOptions()
	IsFalse(t, success, "-help not returning false")
}

func TestParseReturnsTrue(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"orgnetsim", "parse", "../sim/tst.json"}
	success, _, _, _ := parseCommandLineOptions()
	IsTrue(t, success, "not returning true")
}

func TestParseReturnsTrueGetsArgs(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"orgnetsim", "parse", "../sim/tst.csv", "-awm", "-ltp", "-ic", "-mc", "7"}
	success, of, _, _ := parseCommandLineOptions()
	IsTrue(t, success, "not returning true")
	IsTrue(t, of.Network.AgentsWithMemory, "awm not true")
	IsTrue(t, of.Network.LinkTeamPeers, "ltp not true")
	AreEqual(t, of.Network.MaxColors, 7, "wrong max colors")
}

func TestParseReturnsFalseWithErrorArgs(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"orgnetsim", "parse", "../sim/tst.csv", "-be", "-lt", "-mc", "-opt"}
	success, _, opts, _ := parseCommandLineOptions()
	IsFalse(t, success, "command line with errors should return false")
	AreEqual(t, "-be-lt-mc-opt-mc4", opts, "options should be returned")
}
