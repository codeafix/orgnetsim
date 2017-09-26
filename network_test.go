package orgnetsim

import (
	"strings"
	"testing"
)

func TestJsonSerialisationAgent(t *testing.T) {
	json := `{"links":null,"nodes":[{"id":"id_1","color":1,"susceptability":0.2,"influence":0.3,"contrariness":0.4,"change":5}]}`
	n := Network{}
	n.Nodes = append(n.Nodes, &Agent{"id_1", 1, 0.2, 0.3, 0.4, make(chan bool), make(chan string), 5})
	serJSON := n.Serialise()
	if json != serJSON {
		t.Errorf("Serialised json is not identical to original json %s", serJSON)
	}
}

func TestJsonSerialisationLink(t *testing.T) {
	json := `{"links":[{"agent1Id":"id_1","agent2Id":"id_2","strength":0.4}],"nodes":null}`
	n := Network{}
	n.Links = append(n.Links, &Link{"id_1", "id_2", 0.4})
	serJSON := n.Serialise()
	if json != serJSON {
		t.Errorf("Serialised json is not identical to original json %s", serJSON)
	}
}

func TestNewNetworkCreatesValidAgentMap(t *testing.T) {
	json := `{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"}],"links":[{"agent1Id":"id_1","agent2Id":"id_2"},{"agent1Id":"id_1","agent2Id":"id_3"}]}`
	CheckAgentMap(t, json)
}

func TestNewNetworkCreatesValidAgentMapIgnoresDuplicateAgents(t *testing.T) {
	json := `{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"},{"id":"id_1"},{"id":"id_2"}],"links":[{"agent1Id":"id_1","agent2Id":"id_2"},{"agent1Id":"id_1","agent2Id":"id_3"}]}`
	CheckAgentMap(t, json)
}

func CheckAgentMap(t *testing.T, json string) {
	n, err := NewNetwork(json)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(n.AgentsByID) != 3 {
		t.Errorf("Incorrect number of entries in the AgentMap expected 3 got %d", len(n.AgentsByID))
	}
	_, exists := n.AgentsByID["id_1"]
	if !exists {
		t.Errorf("Agent id_1 missing from AgentMap")
	}
	_, exists = n.AgentsByID["id_2"]
	if !exists {
		t.Errorf("Agent id_2 missing from AgentMap")
	}
	_, exists = n.AgentsByID["id_3"]
	if !exists {
		t.Errorf("Agent id_3 missing from AgentMap")
	}
}

func TestNewNetworkCreatesValidLinkMapHierarchy(t *testing.T) {
	json := `{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"}],"links":[{"agent1Id":"id_1","agent2Id":"id_2"},{"agent1Id":"id_1","agent2Id":"id_3"}]}`
	CheckLinkMapHierarchy(t, json)
}

func TestNewNetworkCreatesValidLinkMapHierarchyIgnoresDuplicateLinks(t *testing.T) {
	json := `{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"}],"links":[{"agent1Id":"id_1","agent2Id":"id_2"},{"agent1Id":"id_1","agent2Id":"id_3"},{"agent1Id":"id_1","agent2Id":"id_2"},{"agent1Id":"id_3","agent2Id":"id_1"}]}`
	CheckLinkMapHierarchy(t, json)
}

func CheckLinkMapHierarchy(t *testing.T, json string) {
	n, err := NewNetwork(json)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(n.AgentLinkMap) != 3 {
		t.Errorf("Incorrect number of entries in the AgentLinkMap expected 3 got %d", len(n.AgentLinkMap))
	}
	agent1links, exists := n.AgentLinkMap["id_1"]
	if !exists {
		t.Errorf("Agent id_1 missing from AgentLinkMap")
	}
	if len(agent1links) != 2 {
		t.Errorf("Incorrect number of related nodes to Agent id_1 expected 2 got %d", len(agent1links))
	}
	for id, agentLink := range agent1links {
		if id != agentLink.Agent.ID {
			t.Errorf("Id and Agent.ID not equal in AgentLinkMap expected %s got %s", id, agentLink.Agent.ID)
		}
		if agentLink.Agent.ID != "id_2" && agentLink.Agent.ID != "id_3" {
			t.Errorf("Unexpected agent related to Agent id_1 got %s", agentLink.Agent.ID)
		}
	}
	agent2links, exists := n.AgentLinkMap["id_2"]
	if !exists {
		t.Errorf("Agent id_2 missing from AgentLinkMap")
	}
	if len(agent2links) != 1 {
		t.Errorf("Incorrect number of related nodes to Agent id_2 expected 1 got %d", len(agent2links))
	}
	agent2Link, exists := agent2links["id_1"]
	if !exists || "id_1" != agent2Link.Agent.ID {
		t.Errorf("Unexpected agent related to Agent id_2 expected id_1 %v", agent2links)
	}
	agent3links := n.AgentLinkMap["id_3"]
	if agent3links == nil {
		t.Errorf("Agent id_3 missing from AgentLinkMap")
	}
	if len(agent3links) != 1 {
		t.Errorf("Incorrect number of related nodes to Agent id_3 expected 1 got %d", len(agent3links))
	}
	agent3Link, exists := agent3links["id_1"]
	if !exists || "id_1" != agent3Link.Agent.ID {
		t.Errorf("Unexpected agent related to Agent id_3 expected id_1 %v", agent3links)
	}
}

func TestNewNetworkCreatesValidLinkMapTwoParents(t *testing.T) {
	json := `{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"}],"links":[{"agent1Id":"id_1","agent2Id":"id_3"},{"agent1Id":"id_2","agent2Id":"id_3"}]}`
	CheckLinkMapTwoParents(t, json)
}

func TestNewNetworkCreatesValidLinkMapTwoParentsIgnoresDuplicateLinks(t *testing.T) {
	json := `{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"}],"links":[{"agent1Id":"id_1","agent2Id":"id_3"},{"agent1Id":"id_2","agent2Id":"id_3"}, {"agent1Id":"id_1","agent2Id":"id_3"}]}`
	CheckLinkMapTwoParents(t, json)
}

func CheckLinkMapTwoParents(t *testing.T, json string) {
	n, err := NewNetwork(json)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(n.AgentLinkMap) != 3 {
		t.Errorf("Incorrect number of entries in the AgentLinkMap expected 3 got %d", len(n.AgentLinkMap))
	}
	agent1links, exists := n.AgentLinkMap["id_1"]
	if !exists {
		t.Errorf("Agent id_1 missing from AgentLinkMap")
	}
	if len(agent1links) != 1 {
		t.Errorf("Incorrect number of related nodes to Agent id_1 expected 1 got %d", len(agent1links))
	}
	agent1Link, exists := agent1links["id_3"]
	if !exists || "id_3" != agent1Link.Agent.ID {
		t.Errorf("Unexpected agent related to Agent id_1 expected id_3 %v", agent1links)
	}
	agent2links, exists := n.AgentLinkMap["id_2"]
	if !exists {
		t.Errorf("Agent id_2 missing from AgentLinkMap")
	}
	if len(agent2links) != 1 {
		t.Errorf("Incorrect number of related nodes to Agent id_2 expected 1 got %d", len(agent2links))
	}
	agent2Link, exists := agent2links["id_3"]
	if !exists || "id_3" != agent2Link.Agent.ID {
		t.Errorf("Unexpected agent related to Agent id_2 expected id_3 %v", agent2links)
	}
	agent3links := n.AgentLinkMap["id_3"]
	if agent3links == nil {
		t.Errorf("Agent id_1 missing from AgentLinkMap")
	}
	if len(agent3links) != 2 {
		t.Errorf("Incorrect number of related nodes to Agent id_3 expected 2 got %d", len(agent3links))
	}
	for id, agentLink := range agent3links {
		if id != agentLink.Agent.ID {
			t.Errorf("Id and Agent.ID not equal in AgentLinkMap expected %s got %s", id, agentLink.Agent.ID)
		}
		if agentLink.Agent.ID != "id_1" && agentLink.Agent.ID != "id_2" {
			t.Errorf("Unexpected agent related to Agent id_3 got %s", agentLink.Agent.ID)
		}
	}
}

func TestNewNetworkFailsWhenInvalidIdAppearsInLinks1(t *testing.T) {
	json := `{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"}],"links":[{"agent1Id":"id_1","agent2Id":"id_2"},{"agent1Id":"id_5","agent2Id":"id_3"}]}`
	_, err := NewNetwork(json)
	if err == nil {
		t.Errorf("Expecting error to be reported and it wasn't")
	}
	if !strings.Contains(err.Error(), "id_5") {
		t.Errorf("Incorrect error reported: %s", err.Error())
	}
}

func TestNewNetworkFailsWhenInvalidIdAppearsInLinks2(t *testing.T) {
	json := `{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"}],"links":[{"agent1Id":"id_1","agent2Id":"id_7"},{"agent1Id":"id_1","agent2Id":"id_3"}]}`
	_, err := NewNetwork(json)
	if err == nil {
		t.Errorf("Expecting error to be reported and it wasn't")
	}
	if !strings.Contains(err.Error(), "id_7") {
		t.Errorf("Incorrect error reported: %s", err.Error())
	}
}

func TestGetRelatedAgentsReturnsCorrectList(t *testing.T) {
	json := `{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"},{"id":"id_4"},{"id":"id_5"}],"links":[{"agent1Id":"id_1","agent2Id":"id_2"},{"agent1Id":"id_1","agent2Id":"id_3"},{"agent1Id":"id_4","agent2Id":"id_1"},{"agent1Id":"id_5","agent2Id":"id_1"}]}`
	n, err := NewNetwork(json)
	if err != nil {
		t.Errorf(err.Error())
	}
	relatedAgents := n.GetRelatedAgents(n.AgentsByID["id_1"])
	if len(relatedAgents) != 4 {
		t.Errorf("Incorrect number of related Agents to Agent id_1 expected 4 got %d", len(relatedAgents))
	}
	checks := map[string]bool{
		"id_2": false,
		"id_3": false,
		"id_4": false,
		"id_5": false,
	}
	for _, agent := range relatedAgents {
		_, exists := checks[agent.ID]
		if exists {
			checks[agent.ID] = true
		} else {
			t.Errorf("Unexpected related Agent %s", agent.ID)
		}
	}
	for id, check := range checks {
		if !check {
			t.Errorf("Expected Agent not found in related agents list %s", id)
		}
	}
}

func TestGetAgentByIDReturnsCorrectAgentWhenExists(t *testing.T) {
	json := `{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"},{"id":"id_4"},{"id":"id_5"}],"links":[{"agent1Id":"id_1","agent2Id":"id_2"},{"agent1Id":"id_1","agent2Id":"id_3"},{"agent1Id":"id_4","agent2Id":"id_1"},{"agent1Id":"id_5","agent2Id":"id_1"}]}`
	n, err := NewNetwork(json)
	if err != nil {
		t.Errorf(err.Error())
	}
	agent := n.GetAgentByID("id_5")
	if agent.ID != "id_5" {
		t.Errorf("Expected agent to be id_5 but got %v", agent)
	}
}

func TestGetAgentByIDReturnsNilWhenNotExists(t *testing.T) {
	json := `{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"},{"id":"id_4"},{"id":"id_5"}],"links":[{"agent1Id":"id_1","agent2Id":"id_2"},{"agent1Id":"id_1","agent2Id":"id_3"},{"agent1Id":"id_4","agent2Id":"id_1"},{"agent1Id":"id_5","agent2Id":"id_1"}]}`
	n, err := NewNetwork(json)
	if err != nil {
		t.Errorf(err.Error())
	}
	agent := n.GetAgentByID("id_9")
	if agent != nil {
		t.Errorf("Expected nil to be returned but got %v", agent)
	}
}

func TestUpdateLinkStrength(t *testing.T) {
	json := `{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"},{"id":"id_4"},{"id":"id_5"}],"links":[{"agent1Id":"id_1","agent2Id":"id_2"},{"agent1Id":"id_1","agent2Id":"id_3"},{"agent1Id":"id_4","agent2Id":"id_1"},{"agent1Id":"id_5","agent2Id":"id_1"}]}`
	n, err := NewNetwork(json)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = n.UpdateLinkStrength("id_5", "id_1", 3.1)
	if err != nil {
		t.Errorf(err.Error())
	}
	if n.AgentLinkMap["id_5"]["id_1"].Link.Strength != 3.1 {
		t.Errorf("Link strength not updated")
	}
	if n.AgentLinkMap["id_1"]["id_5"].Link.Strength != 3.1 {
		t.Errorf("Link strength not updated")
	}
}

func TestUpdateLinkStrengthReportsErrorIfIdDoesntExist(t *testing.T) {
	json := `{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"},{"id":"id_4"},{"id":"id_5"}],"links":[{"agent1Id":"id_1","agent2Id":"id_2"},{"agent1Id":"id_1","agent2Id":"id_3"},{"agent1Id":"id_4","agent2Id":"id_1"},{"agent1Id":"id_5","agent2Id":"id_1"}]}`
	n, err := NewNetwork(json)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = n.UpdateLinkStrength("id_7", "id_1", 3.1)
	if err == nil {
		t.Errorf("Invalid link error not reported")
	}
	if !strings.Contains(err.Error(), "id_7") || !strings.Contains(err.Error(), "id_1") {
		t.Errorf("Incorrect error reported: %s", err.Error())
	}
	err = n.UpdateLinkStrength("id_1", "id_7", 3.1)
	if err == nil {
		t.Errorf("Invalid link error not reported")
	}
	if !strings.Contains(err.Error(), "id_7") || !strings.Contains(err.Error(), "id_1") {
		t.Errorf("Incorrect error reported: %s", err.Error())
	}
}
