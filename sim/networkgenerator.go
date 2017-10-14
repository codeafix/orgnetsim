package sim

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
	MaxColors        int     `json:"maxColors"`
	EvangelistAgents bool    `json:"evangelistAgents"`
	LoneEvangelist   bool    `json:"loneEvangelist"`
	AgentsWithMemory bool    `json:"agentsWithMemory"`
}

//GenerateHierarchy generates a hierarchical network
func GenerateHierarchy(s HierarchySpec) (*Network, error) {
	n := new(Network)
	n.MaxColorCount = s.MaxColors
	nodeCount := new(int)
	*nodeCount = 1
	a := GenerateRandomAgent(generateID(nodeCount), s.InitColors, s.AgentsWithMemory)
	n.AddAgent(a)

	leafTeamCount := int(math.Pow(float64(s.TeamSize), float64(s.TeamLinkLevel-1)))
	leafTeams := make([][]Agent, 0, leafTeamCount)

	generateChildren(n, a, &leafTeams, nodeCount, 0, s)

	if s.LinkTeams {
		for i := 0; i < leafTeamCount; i++ {
			for j := i + 1; j < leafTeamCount; j++ {
				n.AddLink(leafTeams[i][0], leafTeams[j][0])
			}
		}
	}

	if s.EvangelistAgents {
		for i := 0; i < leafTeamCount; i++ {
			leafTeams[i][3].State().Susceptability = 5.0
			leafTeams[i][3].State().Color = Blue
		}
	}

	if s.LoneEvangelist {
		doa := GenerateRandomAgent(generateID(nodeCount), s.InitColors, s.AgentsWithMemory)
		doa.State().Susceptability = 5.0
		doa.State().Color = Blue
		n.AddAgent(doa)
		for i := 0; i < leafTeamCount; i++ {
			n.AddLink(doa, leafTeams[i][2])
		}
	}

	err := n.PopulateMaps()
	return n, err
}

func generateChildren(n *Network, parent Agent, leafTeams *[][]Agent, nodeCount *int, level int, s HierarchySpec) {
	level++
	if level >= s.Levels {
		return
	}

	peers := make([]Agent, s.TeamSize, s.TeamSize)
	for i := 0; i < s.TeamSize; i++ {
		a := GenerateRandomAgent(generateID(nodeCount), s.InitColors, s.AgentsWithMemory)
		peers[i] = a
		n.AddAgent(a)
		n.AddLink(parent, a)
		generateChildren(n, a, leafTeams, nodeCount, level, s)
	}

	//Add peer links
	if s.LinkTeamPeers {
		for i := 0; i < s.TeamSize; i++ {
			for j := i + 1; j < s.TeamSize; j++ {
				n.AddLink(peers[i], peers[j])
			}
		}
	}

	if level == s.TeamLinkLevel {
		*leafTeams = append(*leafTeams, peers)
	}
}

func generateID(agentCount *int) string {
	s := fmt.Sprintf("id_%d", *agentCount)
	*agentCount++
	return s
}

//GenerateRandomAgent creates an Agent with random properties
func GenerateRandomAgent(id string, initColors []Color, withMemory bool) Agent {
	as := AgentState{
		ID:             id,
		Color:          Grey,
		Influence:      rand.NormFloat64()*0.25 + 1,
		Susceptability: rand.NormFloat64()*0.25 + 1,
		Contrariness:   rand.NormFloat64()*0.15 + 0.7,
		Mail:           nil,
		ChangeCount:    0,
		Type:           "AgentWithMemory",
	}
	if len(initColors) > 0 {
		as.Color = initColors[rand.Intn(len(initColors))]
	}
	if withMemory {
		a := AgentWithMemory{
			AgentState:     as,
			PreviousColors: nil,
			ShortMemory:    nil,
			MaxColors:      0,
		}
		return &a
	}

	return &as
}
