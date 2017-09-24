package orgnetsim

import (
	"encoding/json"
	"errors"
	"fmt"
)

//A Network of Agents
type Network struct {
	Links        []Link                      `json:"links"`
	Nodes        []Agent                     `json:"nodes"`
	AgentsByID   map[string]Agent            `json:"-"`
	AgentLinkMap map[string]map[string]Agent `json:"-"`
}

// NewNetwork creates a new Network structure from the passed json string
func NewNetwork(jsonBody string) (*Network, error) {
	n := Network{}
	json.Unmarshal([]byte(jsonBody), &n)
	err := n.PopulateMaps()
	return &n, err
}

// PopulateMaps creates the map lookups from the Links and Nodes arrays
func (n *Network) PopulateMaps() error {
	err := ""
	n.AgentsByID = make(map[string]Agent, len(n.Nodes))
	for _, agent := range n.Nodes {
		n.AgentsByID[agent.ID] = agent
	}
	n.AgentLinkMap = make(map[string]map[string]Agent, len(n.Nodes))
	for _, link := range n.Links {
		agent1, exists := n.AgentsByID[link.Agent1ID]
		if !exists {
			err = fmt.Sprintf("%sAgent1ID %s not found in list of Agents\n", err, link.Agent1ID)
			continue
		}
		agent2, exists := n.AgentsByID[link.Agent2ID]
		if !exists {
			err = fmt.Sprintf("%sAgent2ID %s not found in list of Agents\n", err, link.Agent2ID)
			continue
		}
		agent1Map, exists := n.AgentLinkMap[link.Agent1ID]
		if !exists {
			n.AgentLinkMap[link.Agent1ID] = map[string]Agent{agent2.ID: agent2}
		} else {
			agent1Map[agent2.ID] = agent2
		}
		agent2Map, exists := n.AgentLinkMap[link.Agent2ID]
		if !exists {
			n.AgentLinkMap[link.Agent2ID] = map[string]Agent{agent1.ID: agent1}
		} else {
			agent2Map[agent1.ID] = agent1
		}
	}
	if "" != err {
		return errors.New(err)
	}
	return nil
}

//GetRelatedAgents returns a list of Agents adjacent in the Network to the passed Agent
func (n *Network) GetRelatedAgents(a Agent) []Agent {
	r := make([]Agent, 0, len(n.AgentLinkMap[a.ID]))
	for _, agent := range n.AgentLinkMap[a.ID] {
		r = append(r, agent)
	}
	return r
}

// GetAgentByID returns the Agent with the given ID
func (n *Network) GetAgentByID(id string) Agent {
	return n.AgentsByID[id]
}
