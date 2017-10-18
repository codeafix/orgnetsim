package sim

import (
	"fmt"
	"math"
	"math/rand"
)

//HierarchySpec provides parameters to the GenerateHierarchy function specifying the features of the
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
	err := n.PopulateMaps()
	if err != nil {
		return n, err
	}

	o := NetworkOptions{
		LinkTeamPeers:    s.LinkTeamPeers,
		InitColors:       s.InitColors,
		MaxColors:        s.MaxColors,
		AgentsWithMemory: s.AgentsWithMemory,
	}

	if s.LinkTeams {
		for i := 0; i < leafTeamCount; i++ {
			o.LinkedTeamList = append(o.LinkedTeamList, leafTeams[i][0].Identifier())
		}
	}

	if s.EvangelistAgents {
		for i := 0; i < leafTeamCount; i++ {
			o.EvangelistList = append(o.EvangelistList, leafTeams[i][3].Identifier())
		}
	}

	if s.LoneEvangelist {
		o.LoneEvangelist = append(o.LoneEvangelist, generateID(nodeCount))
		for i := 0; i < leafTeamCount; i++ {
			o.LoneEvangelist = append(o.LoneEvangelist, leafTeams[i][2].Identifier())
		}
	}

	err = ModifyNetwork(n, o)
	if err != nil {
		return n, err
	}

	err = n.PopulateMaps()
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
		Type:           "Agent",
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
		a.Type = "AgentWithMemory"
		return &a
	}

	return &as
}
