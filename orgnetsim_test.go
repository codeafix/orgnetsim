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
func WriteOutput(t *testing.T, s HierarchySpec, n RelationshipMgr, colors [][]int, conversations []int) {
	f, err := os.Create("./out.csv")
	AssertSuccess(t, err)
	defer f.Close()

	var buffer bytes.Buffer

	for c := 0; c < MaxColors; c++ {
		buffer.WriteString(fmt.Sprintf("%s,", Color(c).String()))
	}

	buffer.WriteString("Conversations,,Node,Change Count,,Link,Strength,,Levels,TeamSize,TeamLinkLevel,LinkTeamPeers,LinkTeams,InitColors,EvangelistAgents,LoneEvangelist,AgentsWithMemory\n")

	agents := n.Agents()
	links := n.Links()
	iterations := len(conversations)
	agentCount := len(agents)
	linkCount := len(links)
	totalLines := iterations
	if agentCount > totalLines {
		totalLines = agentCount
	}
	if linkCount > totalLines {
		totalLines = linkCount
	}

	for i := 0; i < totalLines; i++ {
		if i < iterations {
			for j := 0; j < MaxColors; j++ {
				buffer.WriteString(fmt.Sprintf("%d,", colors[i][j]))
			}
			buffer.WriteString(fmt.Sprintf("%d", conversations[i]))
		} else {
			for j := 0; j < MaxColors; j++ {
				buffer.WriteString(",")
			}
		}
		if i < agentCount {
			buffer.WriteString(fmt.Sprintf(",,%s,%d", agents[i].Identifier(), agents[i].State().ChangeCount))
		} else {
			buffer.WriteString(",,,")
		}
		if i < linkCount {
			buffer.WriteString(fmt.Sprintf(",,%s-%s,%d", links[i].Agent1ID, links[i].Agent2ID, links[i].Strength))
		} else {
			buffer.WriteString(",,,")
		}
		if i == 0 {
			var initColors string
			for x := 0; x < len(s.InitColors); x++ {
				initColors = initColors + Color(s.InitColors[x]).String()
			}
			buffer.WriteString(fmt.Sprintf(",,%d,%d,%d,%t,%t,%s,%t,%t,%t\n", s.Levels, s.TeamSize, s.TeamLinkLevel, s.LinkTeamPeers, s.LinkTeams, initColors, s.EvangelistAgents, s.LoneEvangelist, s.AgentsWithMemory))
		} else {
			buffer.WriteString("\n")
		}
	}

	_, err = f.Write(buffer.Bytes())
	AssertSuccess(t, err)
}

func TestRunSim(t *testing.T) {

	s := HierarchySpec{
		4,                  //Levels
		5,                  //TeamSize
		3,                  //TeamLinkLevel
		true,               //LinkTeamPeers
		true,               //LinkTeams
		[]Color{Grey, Red}, //InitColors
		false,              //EvangelistAgents
		false,              //LoneEvangelist
		false,              //AgentsWithMemory
	}

	n, err := GenerateHierarchy(s)
	AssertSuccess(t, err)

	colors, conversations := RunSim(n, 500)
	WriteOutput(t, s, n, colors, conversations)
}
