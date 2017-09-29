package orgnetsim

import (
	"bytes"
	"fmt"
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

//Convenience method to dump the colors and conversations arrays into a csv file
func WriteOutput(t *testing.T, s HierarchySpec, colors [][]int, conversations []int) {
	f, err := os.Create("./out.csv")
	AssertSuccess(t, err)
	defer f.Close()

	var buffer bytes.Buffer

	for c := 0; c < MaxColors; c++ {
		buffer.WriteString(fmt.Sprintf("%s,", Color(c).String()))
	}

	buffer.WriteString("Conversations,Levels,TeamSize,TeamLinkLevel,LinkTeamPeers,LinkTeams,InitColors,EvangelistAgents,LoneEvangelist\n")
	for i := 0; i < len(conversations); i++ {
		for j := 0; j < MaxColors; j++ {
			buffer.WriteString(fmt.Sprintf("%d,", colors[i][j]))
		}
		buffer.WriteString(fmt.Sprintf("%d", conversations[i]))
		if i == 0 {
			var initColors string
			for x := 0; x < len(s.InitColors); x++ {
				initColors = initColors + Color(s.InitColors[x]).String()
			}
			buffer.WriteString(fmt.Sprintf(",%d,%d,%d,%t,%t,%s,%t,%t\n", s.Levels, s.TeamSize, s.TeamLinkLevel, s.LinkTeamPeers, s.LinkTeams, initColors, s.EvangelistAgents, s.LoneEvangelist))
		} else {
			buffer.WriteString("\n")
		}
	}

	_, err = f.Write(buffer.Bytes())
	AssertSuccess(t, err)
}

func TestRunSim(t *testing.T) {

	s := HierarchySpec{
		4,     //Levels
		5,     //TeamSize
		3,     //TeamLinkLevel
		true,  //LinkTeamPeers
		true,  //LinkTeams
		nil,   //InitColors
		false, //EvangelistAgents
		false, //LoneEvangelist
	}

	n, err := GenerateHierarchy(s)
	AssertSuccess(t, err)

	colors, conversations := RunSim(n, 500)
	WriteOutput(t, s, colors, conversations)
}
