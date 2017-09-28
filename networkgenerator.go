package orgnetsim

import (
	"fmt"
	"math"
	"math/rand"
)

//GenerateHierarchy generates a hierarchical network
func GenerateHierarchy() (*Network, error) {
	levels, teamsize, teamLinkLevel := 4, 5, 3
	peerLinks, teamLinks, initColors, unsusceptibleAgents, do := true, true, false, false, true

	n := new(Network)
	nodeCount := new(int)
	*nodeCount = 1
	a := GenerateRandomAgent(nodeCount, initColors)
	n.Nodes = append(n.Nodes, a)

	leafTeamCount := int(math.Pow(float64(teamsize), float64(teamLinkLevel-1)))
	leafTeams := make([][]*Agent, 0, leafTeamCount)

	generateChildren(n, a, &leafTeams, nodeCount, 0, levels, teamsize, teamLinkLevel, peerLinks, teamLinks, initColors, unsusceptibleAgents)

	if teamLinks {
		for i := 0; i < leafTeamCount; i++ {
			for j := i + 1; j < leafTeamCount; j++ {
				l := NewLink(leafTeams[i][0], leafTeams[j][0])
				n.Links = append(n.Links, l)
			}
		}
	}

	if unsusceptibleAgents {
		for i := 0; i < leafTeamCount; i++ {
			leafTeams[i][3].Susceptability = 5.0
			leafTeams[i][3].Color = Blue
		}
	}

	if do {
		doa := GenerateRandomAgent(nodeCount, initColors)
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

func generateChildren(n *Network, parent *Agent, leafTeams *[][]*Agent, nodeCount *int, level int, levels int, teamsize int, teamLinkLevel int, peerLinks bool, teamLinks bool, initColors bool, unsusceptibleAgents bool) {
	level++
	if level >= levels {
		return
	}

	peers := make([]*Agent, teamsize, teamsize)
	for i := 0; i < teamsize; i++ {
		a := GenerateRandomAgent(nodeCount, initColors)
		peers[i] = a
		l := NewLink(parent, a)
		n.Nodes = append(n.Nodes, a)
		n.Links = append(n.Links, l)
		generateChildren(n, a, leafTeams, nodeCount, level, levels, teamsize, teamLinkLevel, peerLinks, teamLinks, initColors, unsusceptibleAgents)
	}

	//Add peer links
	if peerLinks {
		for i := 0; i < teamsize; i++ {
			for j := i + 1; j < teamsize; j++ {
				l := NewLink(peers[i], peers[j])
				n.Links = append(n.Links, l)
			}
		}
	}

	if level == teamLinkLevel {
		*leafTeams = append(*leafTeams, peers)
	}
}

//GenerateRandomAgent creates an Agent with random properties
func GenerateRandomAgent(agentCount *int, initColors bool) *Agent {
	a := Agent{
		fmt.Sprintf("id_%d", *agentCount),
		Grey,
		rand.NormFloat64()*0.25 + 1,
		rand.NormFloat64()*0.25 + 1,
		rand.NormFloat64()*0.15 + 0.7,
		nil,
		0,
	}
	if initColors {
		if rand.Intn(2) == 1 {
			a.Color = Red
		} else {
			a.Color = Grey
		}
	}
	*agentCount++
	return &a
}

//NewLink returns a Link between to two passed agents
func NewLink(a1 *Agent, a2 *Agent) *Link {
	l := Link{
		a1.ID,
		a2.ID,
		0,
	}
	return &l
}
