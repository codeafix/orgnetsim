package orgnetsim

import (
	"fmt"
	"math/rand"
)

//GenerateHierarchy generates a hierarchical network
func GenerateHierarchy() (*Network, error) {
	levels, teamsize := 4, 5
	peerLinks, teamLinks, initColors, unsusceptibleAgents := false, false, false, false

	n := new(Network)
	nodeCount := new(int)
	*nodeCount = 1
	a := GenerateRandomAgent(nodeCount, initColors)
	generateChildren(n, a, nodeCount, 0, levels, teamsize, peerLinks, teamLinks, initColors, unsusceptibleAgents)

	err := n.PopulateMaps()
	return n, err
}

func generateChildren(n *Network, parent *Agent, nodeCount *int, level int, levels int, teamsize int, peerLinks bool, teamLinks bool, initColors bool, unsusceptibleAgents bool) *Agent {
	level++
	if level >= levels {
		return nil
	}
	levelHasChildren := level < (levels - 1)

	peers := make([]*Agent, teamsize, teamsize)
	children := make([]*Agent, teamsize, teamsize)
	for i := 0; i < teamsize; i++ {
		a := GenerateRandomAgent(nodeCount, initColors)
		peers[i] = a
		l := NewLink(parent, a)
		n.Nodes = append(n.Nodes, a)
		n.Links = append(n.Links, l)
		children[i] = generateChildren(n, a, nodeCount, level, levels, teamsize, peerLinks, teamLinks, initColors, unsusceptibleAgents)
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

	//Add an additional network that links 1 member from each team
	if teamLinks && levelHasChildren {
		for i := 0; i < teamsize; i++ {
			for j := i + 1; j < teamsize; j++ {
				l := NewLink(children[i], children[j])
				n.Links = append(n.Links, l)
			}
		}
	}

	//Mark certain Agents as unsusceptible and give them a Color
	if unsusceptibleAgents && levelHasChildren {
		for i := 0; i < teamsize; i++ {
			children[i].Susceptability = 5.0
			children[i].Color = Blue
		}
	}

	return peers[0]
}

//GenerateRandomAgent creates an Agent with random properties
func GenerateRandomAgent(agentCount *int, initColors bool) *Agent {
	a := Agent{
		fmt.Sprintf("id_%d", *agentCount),
		Grey,
		rand.NormFloat64()*3 + 1,
		rand.NormFloat64()*3 + 1,
		rand.NormFloat64()*3 + 1,
		nil,
		0,
	}
	if initColors {
		if rand.Intn(2) == 1 {
			a.Color = Red
		} else {
			a.Color = Blue
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
