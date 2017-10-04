package orgnetsim

import (
	"fmt"
	"math"
	"math/rand"
)

//HierarchySpec provides parameters to the GenerateHierarchy functioning specifying the features of the
//Hierarchical network to generate
type HierarchySpec struct {
	Levels           int     `json:"levels"`
	TeamSize         int     `json:"teamSize"`
	TeamLinkLevel    int     `json:"teamLinkLevel"`
	LinkTeamPeers    bool    `json:"linkTeamPeers"`
	LinkTeams        bool    `json:"linkTeams"`
	InitColors       []Color `json:"initColors"`
	EvangelistAgents bool    `json:"evangelistAgents"`
	LoneEvangelist   bool    `json:"loneEvangelist"`
}

//GenerateHierarchy generates a hierarchical network
func GenerateHierarchy(s HierarchySpec) (*Network, error) {
	n := new(Network)
	nodeCount := new(int)
	*nodeCount = 1
	a := GenerateRandomAgent(nodeCount, s.InitColors)
	n.Nodes = append(n.Nodes, a)

	leafTeamCount := int(math.Pow(float64(s.TeamSize), float64(s.TeamLinkLevel-1)))
	leafTeams := make([][]*AgentWithMemory, 0, leafTeamCount)

	generateChildren(n, a, &leafTeams, nodeCount, 0, s)

	if s.LinkTeams {
		for i := 0; i < leafTeamCount; i++ {
			for j := i + 1; j < leafTeamCount; j++ {
				l := NewLink(leafTeams[i][0], leafTeams[j][0])
				n.Links = append(n.Links, l)
			}
		}
	}

	if s.EvangelistAgents {
		for i := 0; i < leafTeamCount; i++ {
			leafTeams[i][3].Susceptability = 5.0
			leafTeams[i][3].Color = Blue
		}
	}

	if s.LoneEvangelist {
		doa := GenerateRandomAgent(nodeCount, s.InitColors)
		doa.Susceptability = 5.0
		doa.Color = Blue
		n.Nodes = append(n.Nodes, doa)
		for i := 0; i < leafTeamCount; i++ {
			l := NewLink(doa, leafTeams[i][2])
			n.Links = append(n.Links, l)
		}
	}

	err := n.PopulateMaps()
	return n, err
}

func generateChildren(n *Network, parent Agent, leafTeams *[][]*AgentWithMemory, nodeCount *int, level int, s HierarchySpec) {
	level++
	if level >= s.Levels {
		return
	}

	peers := make([]*AgentWithMemory, s.TeamSize, s.TeamSize)
	for i := 0; i < s.TeamSize; i++ {
		a := GenerateRandomAgent(nodeCount, s.InitColors)
		peers[i] = a
		l := NewLink(parent, a)
		n.Nodes = append(n.Nodes, a)
		n.Links = append(n.Links, l)
		generateChildren(n, a, leafTeams, nodeCount, level, s)
	}

	//Add peer links
	if s.LinkTeamPeers {
		for i := 0; i < s.TeamSize; i++ {
			for j := i + 1; j < s.TeamSize; j++ {
				l := NewLink(peers[i], peers[j])
				n.Links = append(n.Links, l)
			}
		}
	}

	if level == s.TeamLinkLevel {
		*leafTeams = append(*leafTeams, peers)
	}
}

//GenerateRandomAgent creates an Agent with random properties
func GenerateRandomAgent(agentCount *int, initColors []Color) *AgentWithMemory {
	a := AgentWithMemory{
		AgentState{
			fmt.Sprintf("id_%d", *agentCount),
			Grey,
			rand.NormFloat64()*0.25 + 1,
			rand.NormFloat64()*0.25 + 1,
			rand.NormFloat64()*0.15 + 0.7,
			nil,
			0,
			"AgentWithMemory",
		},
		nil,
	}
	if len(initColors) > 0 {
		a.Color = initColors[rand.Intn(len(initColors))]
	}
	*agentCount++
	return &a
}

//NewLink returns a Link between to two passed agents
func NewLink(a1 Agent, a2 Agent) *Link {
	l := Link{
		a1.Identifier(),
		a2.Identifier(),
		0,
	}
	return &l
}
