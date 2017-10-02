package orgnetsim

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"sort"
)

//A Network of Agents
type Network struct {
	Links        []*Link                         `json:"links"`
	Nodes        []*Agent                        `json:"nodes"`
	AgentsByID   map[string]*Agent               `json:"-"`
	AgentLinkMap map[string]map[string]AgentLink `json:"-"`
}

//AgentLink holds both the Link and the Agent in the AgentLinkMap
type AgentLink struct {
	Agent *Agent
	Link  *Link
}

//RelationshipMgr is an interface for the Network
type RelationshipMgr interface {
	GetRelatedAgents(a *Agent) []*Agent
	GetAgentByID(id string) *Agent
	IncrementLinkStrength(id1 string, id2 string) error
	Agents() []*Agent
}

//Agents returns a list of the Agents Communicating on the Network
func (n *Network) Agents() []*Agent {
	return n.Nodes
}

//Serialise returns a json representation of the Network
func (n *Network) Serialise() string {
	jsonBody, _ := json.Marshal(n)
	return string(jsonBody)
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
	n.AgentsByID = make(map[string]*Agent, len(n.Nodes))
	for _, agent := range n.Nodes {
		n.AgentsByID[agent.ID] = agent
		agent.Mail = make(chan string, 1)
		agent.Memory = make(map[Color]struct{}, MaxColors)
		agent.Memory[Grey] = struct{}{}
	}
	n.AgentLinkMap = make(map[string]map[string]AgentLink, len(n.Nodes))
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
			agent1Map = map[string]AgentLink{}
			n.AgentLinkMap[link.Agent1ID] = agent1Map
		}
		agent1Map[agent2.ID] = AgentLink{agent2, link}
		agent2Map, exists := n.AgentLinkMap[link.Agent2ID]
		if !exists {
			agent2Map = map[string]AgentLink{}
			n.AgentLinkMap[link.Agent2ID] = agent2Map
		}
		agent2Map[agent1.ID] = AgentLink{agent1, link}
	}
	if "" != err {
		return errors.New(err)
	}
	return nil
}

//GetRelatedAgents returns a slice of Agents adjacent in the Network to the passed Agent
//The returned slice of Agents is always deliberately shuffled into random order
func (n *Network) GetRelatedAgents(a *Agent) []*Agent {
	lnkdagents := n.AgentLinkMap[a.ID]
	acnt := len(n.AgentLinkMap[a.ID])
	keys := make([]float64, acnt)
	raMap := make(map[float64]*Agent, acnt)
	i := 0
	for _, agentLink := range lnkdagents {
		keys[i] = rand.Float64()
		raMap[keys[i]] = agentLink.Agent
		i = i + 1
	}
	sort.Float64s(keys)
	r := make([]*Agent, acnt)
	for i, key := range keys {
		r[i] = raMap[key]
	}
	return r
}

// GetAgentByID returns a reference to the Agent with the given ID or nil if it doesn't exist
func (n *Network) GetAgentByID(id string) *Agent {
	a := n.AgentsByID[id]
	return a
}

// IncrementLinkStrength updates the strength field of the link connecting Agents id1 and id2.
// Returns an error if no link is found
func (n *Network) IncrementLinkStrength(id1 string, id2 string) error {
	lnkdMap, exists := n.AgentLinkMap[id1]
	if !exists {
		return fmt.Errorf("Invalid Link id1=%s id2=%s", id1, id2)
	}
	agentLink, exists := lnkdMap[id2]
	if !exists {
		return fmt.Errorf("Invalid Link id1=%s id2=%s", id1, id2)
	}
	agentLink.Link.Strength++
	return nil
}
