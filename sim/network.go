package sim

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"sort"
)

//A Network of Agents
type Network struct {
	Edges         []*Link                         `json:"links"`
	Nodes         []Agent                         `json:"nodes"`
	AgentsByID    map[string]Agent                `json:"-"`
	AgentLinkMap  map[string]map[string]AgentLink `json:"-"`
	MaxColorCount int                             `json:"maxColors"`
}

type networkJSON struct {
	Links []json.RawMessage
	Nodes []json.RawMessage
}

//AgentLink holds both the Link and the Agent in the AgentLinkMap
type AgentLink struct {
	Agent Agent
	Link  *Link
}

//RelationshipMgr is an interface for the Network
type RelationshipMgr interface {
	GetRelatedAgents(a Agent) []Agent
	GetAgentByID(id string) Agent
	IncrementLinkStrength(id1 string, id2 string) error
	AddAgent(a Agent)
	AddLink(a1 Agent, a2 Agent)
	Agents() []Agent
	Links() []*Link
	MaxColors() int
}

//MaxColors returns the maximum number of color states that the agents are permitted on this network
func (n *Network) MaxColors() int {
	return n.MaxColorCount
}

//Agents returns a list of the Agents Communicating on the Network
func (n *Network) Agents() []Agent {
	return n.Nodes
}

//Links returns a list of the links between Agents on the Network
func (n *Network) Links() []*Link {
	return n.Edges
}

//AddAgent adds a new Agent to the network
func (n *Network) AddAgent(a Agent) {
	n.Nodes = append(n.Nodes, a)
}

//AddLink adds a new Link between the two passed agents
func (n *Network) AddLink(a1 Agent, a2 Agent) {
	l := Link{
		a1.Identifier(),
		a2.Identifier(),
		0,
	}
	n.Edges = append(n.Edges, &l)
}

//Serialise returns a json representation of the Network
func (n *Network) Serialise() string {
	jsonBody, _ := json.Marshal(n)
	return string(jsonBody)
}

//UnmarshalJSON implements unmarshalling of Agents of different types
func (n *Network) UnmarshalJSON(b []byte) error {
	var network map[string]json.RawMessage
	err := json.Unmarshal(b, &network)
	if err != nil {
		return err
	}
	json.Unmarshal(network["links"], &n.Edges)
	json.Unmarshal(network["maxColors"], &n.MaxColorCount)

	var nodes []json.RawMessage
	json.Unmarshal(network["nodes"], &nodes)

	n.Nodes = make([]Agent, len(nodes))

	for i, raw := range nodes {
		var fm map[string]interface{}
		json.Unmarshal(raw, &fm)
		switch fm["type"] {
		case "Agent":
			a := AgentState{}
			json.Unmarshal(raw, &a)
			n.Nodes[i] = &a
		case "AgentWithMemory":
			a := AgentWithMemory{}
			json.Unmarshal(raw, &a)
			n.Nodes[i] = &a
		default:
			a := AgentState{}
			json.Unmarshal(raw, &a)
			n.Nodes[i] = &a
		}
	}
	return nil
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
		n.AgentsByID[agent.Identifier()] = agent
		agent.Initialise(n)
	}
	n.AgentLinkMap = make(map[string]map[string]AgentLink, len(n.Nodes))
	for _, link := range n.Edges {
		agent1, exists := n.AgentsByID[link.Agent1ID]
		if !exists {
			err = fmt.Sprintf("%sAgent1ID '%s' not found in list of Agents\n", err, link.Agent1ID)
			continue
		}
		agent2, exists := n.AgentsByID[link.Agent2ID]
		if !exists {
			err = fmt.Sprintf("%sAgent2ID '%s' not found in list of Agents\n", err, link.Agent2ID)
			continue
		}
		agent1Map, exists := n.AgentLinkMap[link.Agent1ID]
		if !exists {
			agent1Map = map[string]AgentLink{}
			n.AgentLinkMap[link.Agent1ID] = agent1Map
		}
		agent1Map[agent2.Identifier()] = AgentLink{agent2, link}
		agent2Map, exists := n.AgentLinkMap[link.Agent2ID]
		if !exists {
			agent2Map = map[string]AgentLink{}
			n.AgentLinkMap[link.Agent2ID] = agent2Map
		}
		agent2Map[agent1.Identifier()] = AgentLink{agent1, link}
	}
	if "" != err {
		return errors.New(err)
	}
	return nil
}

//GetRelatedAgents returns a slice of Agents adjacent in the Network to the passed Agent
//The returned slice of Agents is always deliberately shuffled into random order
func (n *Network) GetRelatedAgents(a Agent) []Agent {
	lnkdagents := n.AgentLinkMap[a.Identifier()]
	acnt := len(n.AgentLinkMap[a.Identifier()])
	keys := make([]float64, acnt)
	raMap := make(map[float64]Agent, acnt)
	i := 0
	for _, agentLink := range lnkdagents {
		keys[i] = rand.Float64()
		raMap[keys[i]] = agentLink.Agent
		i = i + 1
	}
	sort.Float64s(keys)
	r := make([]Agent, acnt)
	for i, key := range keys {
		r[i] = raMap[key]
	}
	return r
}

// GetAgentByID returns a reference to the Agent with the given ID or nil if it doesn't exist
func (n *Network) GetAgentByID(id string) Agent {
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
