package sim

import (
	"fmt"
	"math"
	"math/rand"
)

// HierarchySpec provides parameters to the GenerateHierarchy function specifying the features of the
// Hierarchical network to generate
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

// GenerateHierarchy generates a hierarchical network
func GenerateHierarchy(s HierarchySpec) (*Network, *NetworkOptions, error) {
	n := new(Network)
	n.MaxColorCount = s.MaxColors
	nodeCount := new(int)
	*nodeCount = 1
	a_id, a_name := generateIDAndName(nodeCount)
	a := GenerateRandomAgent(a_id, a_name, s.InitColors, s.AgentsWithMemory)
	n.AddAgent(a)

	leafTeamCount := int(math.Pow(float64(s.TeamSize), float64(s.TeamLinkLevel-1)))
	leafTeams := make([][]Agent, 0, leafTeamCount)

	generateChildren(n, a, &leafTeams, nodeCount, 0, s)
	err := n.PopulateMaps()
	o := CreateNetworkOptions(s)

	if err != nil {
		return n, o, err
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
		le_id, _ := generateIDAndName(nodeCount)
		o.LoneEvangelist = append(o.LoneEvangelist, le_id)
		for i := 0; i < leafTeamCount; i++ {
			o.LoneEvangelist = append(o.LoneEvangelist, leafTeams[i][2].Identifier())
		}
	}

	err = o.ModifyNetwork(n)
	if err != nil {
		return n, o, err
	}

	err = n.PopulateMaps()
	return n, o, err
}

func generateChildren(n *Network, parent Agent, leafTeams *[][]Agent, nodeCount *int, level int, s HierarchySpec) {
	level++
	if level >= s.Levels {
		return
	}

	peers := make([]Agent, s.TeamSize)
	for i := 0; i < s.TeamSize; i++ {
		id, name := generateIDAndName(nodeCount)
		a := GenerateRandomAgent(id, name, s.InitColors, s.AgentsWithMemory)
		peers[i] = a
		n.AddAgent(a)
		n.AddLink(parent, a)
		generateChildren(n, a, leafTeams, nodeCount, level, s)
	}

	if level == s.TeamLinkLevel {
		*leafTeams = append(*leafTeams, peers)
	}
}

func generateIDAndName(agentCount *int) (string, string) {
	id := fmt.Sprintf("id_%d", *agentCount)
	name := fmt.Sprintf("Agent %d", *agentCount)
	*agentCount++
	return id, name
}

// GenerateRandomAgent creates an Agent with random properties
func GenerateRandomAgent(id string, name string, initColors []Color, withMemory bool) Agent {
	as := AgentState{
		ID:             id,
		Name:           name,
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
