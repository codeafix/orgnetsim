package sim

import (
	"fmt"
	"strings"
	"testing"
)

func TestSerialisationOfAgentWithMemory(t *testing.T) {
	json := `{"links":null,"nodes":[{"id":"id_1","color":1,"susceptability":0.2,"influence":0.3,"contrariness":0.4,"change":5,"type":"AgentWithMemory"}],"maxColors":0}`
	n := Network{}
	a := AgentWithMemory{AgentState{"id_1", 1, 0.2, 0.3, 0.4, nil, 5, ""}, nil, nil, 0}
	a.Initialise(&n)
	n.Nodes = append(n.Nodes, &a)
	serJSON := n.Serialise()
	AreEqual(t, json, serJSON, "Serialised json is not identical to original json")
}

func TestDeserialisationOfAgentWithMemory(t *testing.T) {
	json := `{"links":null,"nodes":[{"id":"id_1","color":1,"susceptability":0.2,"influence":0.3,"contrariness":0.4,"change":5,"type":"AgentWithMemory"}],"maxColors":0}`
	n, err := NewNetwork(json)
	AssertSuccess(t, err)
	serJSON := n.Serialise()
	AreEqual(t, json, serJSON, "Serialised json is not identical to original json")
}

func TestJsonSerialisationAgent(t *testing.T) {
	json := `{"links":null,"nodes":[{"id":"id_1","color":1,"susceptability":0.2,"influence":0.3,"contrariness":0.4,"change":5,"type":"Agent"}],"maxColors":0}`
	n := Network{}
	n.Nodes = append(n.Nodes, &AgentState{"id_1", 1, 0.2, 0.3, 0.4, make(chan string), 5, "Agent"})
	serJSON := n.Serialise()
	AreEqual(t, json, serJSON, "Serialised json is not identical to original json")
}

func TestJsonSerialisationLink(t *testing.T) {
	json := `{"links":[{"source":"id_1","target":"id_2","strength":4}],"nodes":null,"maxColors":0}`
	n := Network{}
	n.Edges = append(n.Edges, &Link{"id_1", "id_2", 4})
	serJSON := n.Serialise()
	AreEqual(t, json, serJSON, "Serialised json is not identical to original json")
}

func TestNewNetworkCreatesValidAgentMap(t *testing.T) {
	json := `{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"}],"links":[{"source":"id_1","target":"id_2"},{"source":"id_1","target":"id_3"}]}`
	CheckAgentMap(t, json)
}

func TestNewNetworkCreatesValidAgentMapIgnoresDuplicateAgents(t *testing.T) {
	json := `{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"},{"id":"id_1"},{"id":"id_2"}],"links":[{"source":"id_1","target":"id_2"},{"source":"id_1","target":"id_3"}]}`
	CheckAgentMap(t, json)
}

func CheckAgentMap(t *testing.T, json string) {
	n, err := NewNetwork(json)
	AssertSuccess(t, err)
	AreEqual(t, 3, len(n.AgentsByID), "Incorrect number of entries in the AgentMap")
	_, exists := n.AgentsByID["id_1"]
	IsTrue(t, exists, "Agent id_1 missing from AgentMap")
	_, exists = n.AgentsByID["id_2"]
	IsTrue(t, exists, "Agent id_2 missing from AgentMap")
	_, exists = n.AgentsByID["id_3"]
	IsTrue(t, exists, "Agent id_3 missing from AgentMap")
}

func TestNewNetworkCreatesValidLinkMapHierarchy(t *testing.T) {
	json := `{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"}],"links":[{"source":"id_1","target":"id_2"},{"source":"id_1","target":"id_3"}]}`
	CheckLinkMapHierarchy(t, json)
}

func TestNewNetworkCreatesValidLinkMapHierarchyIgnoresDuplicateLinks(t *testing.T) {
	json := `{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"}],"links":[{"source":"id_1","target":"id_2"},{"source":"id_1","target":"id_3"},{"source":"id_1","target":"id_2"},{"source":"id_3","target":"id_1"}]}`
	CheckLinkMapHierarchy(t, json)
}

func CheckLinkMapHierarchy(t *testing.T, json string) {
	n, err := NewNetwork(json)
	AssertSuccess(t, err)
	AreEqual(t, 3, len(n.AgentLinkMap), "Incorrect number of entries in the AgentLinkMap")

	agent1links, exists := n.AgentLinkMap["id_1"]
	IsTrue(t, exists, "Agent id_1 missing from AgentLinkMap")
	AreEqual(t, 2, len(agent1links), "Incorrect number of related nodes to Agent id_1")
	for id, agentLink := range agent1links {
		AreEqual(t, id, agentLink.Agent.Identifier(), "Id and Agent.ID not equal in AgentLinkMap")
		IsTrue(t, agentLink.Agent.Identifier() == "id_2" || agentLink.Agent.Identifier() == "id_3", fmt.Sprintf("Unexpected agent related to Agent id_1 got %s", agentLink.Agent.Identifier()))
	}

	agent2links, exists := n.AgentLinkMap["id_2"]
	IsTrue(t, exists, "Agent id_2 missing from AgentLinkMap")
	AreEqual(t, 1, len(agent2links), "Incorrect number of related nodes to Agent id_2")
	agent2Link, exists := agent2links["id_1"]
	IsTrue(t, exists && "id_1" == agent2Link.Agent.Identifier(), fmt.Sprintf("Unexpected agent related to Agent id_2 expected id_1 %v", agent2links))

	agent3links, exists := n.AgentLinkMap["id_3"]
	IsTrue(t, exists, "Agent id_3 missing from AgentLinkMap")
	AreEqual(t, 1, len(agent3links), "Incorrect number of related nodes to Agent id_3")
	agent3Link, exists := agent3links["id_1"]
	IsTrue(t, exists && "id_1" == agent3Link.Agent.Identifier(), fmt.Sprintf("Unexpected agent related to Agent id_3 expected id_1 %v", agent3links))
}

func TestNewNetworkCreatesValidLinkMapTwoParents(t *testing.T) {
	json := `{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"}],"links":[{"source":"id_1","target":"id_3"},{"source":"id_2","target":"id_3"}]}`
	CheckLinkMapTwoParents(t, json)
}

func TestNewNetworkCreatesValidLinkMapTwoParentsIgnoresDuplicateLinks(t *testing.T) {
	json := `{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"}],"links":[{"source":"id_1","target":"id_3"},{"source":"id_2","target":"id_3"}, {"source":"id_1","target":"id_3"}]}`
	CheckLinkMapTwoParents(t, json)
}

func CheckLinkMapTwoParents(t *testing.T, json string) {
	n, err := NewNetwork(json)
	AssertSuccess(t, err)
	AreEqual(t, 3, len(n.AgentLinkMap), "Incorrect number of entries in the AgentLinkMap")

	agent1links, exists := n.AgentLinkMap["id_1"]
	IsTrue(t, exists, "Agent id_1 missing from AgentLinkMap")
	AreEqual(t, 1, len(agent1links), "Incorrect number of related nodes to Agent id_1")
	agent1Link, exists := agent1links["id_3"]
	IsTrue(t, exists && "id_3" == agent1Link.Agent.Identifier(), fmt.Sprintf("Unexpected agent related to Agent id_1 expected id_3 %v", agent1links))

	agent2links, exists := n.AgentLinkMap["id_2"]
	IsTrue(t, exists, "Agent id_2 missing from AgentLinkMap")
	AreEqual(t, 1, len(agent2links), "Incorrect number of related nodes to Agent id_2")
	agent2Link, exists := agent2links["id_3"]
	IsTrue(t, exists && "id_3" == agent2Link.Agent.Identifier(), fmt.Sprintf("Unexpected agent related to Agent id_2 expected id_3 %v", agent1links))

	agent3links, exists := n.AgentLinkMap["id_3"]
	IsTrue(t, exists, "Agent id_3 missing from AgentLinkMap")
	AreEqual(t, 2, len(agent3links), "Incorrect number of related nodes to Agent id_3")
	for id, agentLink := range agent3links {
		AreEqual(t, id, agentLink.Agent.Identifier(), "Id and Agent.ID not equal in AgentLinkMap")
		IsTrue(t, agentLink.Agent.Identifier() == "id_1" || agentLink.Agent.Identifier() == "id_2", fmt.Sprintf("Unexpected agent related to Agent id_3 got %s", agentLink.Agent.Identifier()))
	}
}

func TestNewNetworkFailsWhenInvalidIdAppearsInLinks1(t *testing.T) {
	json := `{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"}],"links":[{"source":"id_1","target":"id_2"},{"source":"id_5","target":"id_3"}]}`
	_, err := NewNetwork(json)
	IsFalse(t, err == nil, "Expecting error to be reported and it wasn't")
	IsTrue(t, strings.Contains(err.Error(), "id_5"), fmt.Sprintf("Incorrect error reported: %s", err.Error()))
}

func TestNewNetworkFailsWhenInvalidIdAppearsInLinks2(t *testing.T) {
	json := `{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"}],"links":[{"source":"id_1","target":"id_7"},{"source":"id_1","target":"id_3"}]}`
	_, err := NewNetwork(json)
	IsFalse(t, err == nil, "Expecting error to be reported and it wasn't")
	IsTrue(t, strings.Contains(err.Error(), "id_7"), fmt.Sprintf("Incorrect error reported: %s", err.Error()))
}

func TestGetRelatedAgentsReturnsCorrectList(t *testing.T) {
	json := `{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"},{"id":"id_4"},{"id":"id_5"}],"links":[{"source":"id_1","target":"id_2"},{"source":"id_1","target":"id_3"},{"source":"id_4","target":"id_1"},{"source":"id_5","target":"id_1"}]}`
	n, err := NewNetwork(json)
	AssertSuccess(t, err)
	relatedAgents := n.GetRelatedAgents(n.AgentsByID["id_1"])
	AreEqual(t, 4, len(relatedAgents), "Incorrect number of related Agents to Agent id_1")
	checks := map[string]bool{
		"id_2": false,
		"id_3": false,
		"id_4": false,
		"id_5": false,
	}
	for _, agent := range relatedAgents {
		_, exists := checks[agent.Identifier()]
		if exists {
			checks[agent.Identifier()] = true
		} else {
			t.Errorf("Unexpected related Agent %s", agent.Identifier())
		}
	}
	for id, check := range checks {
		IsTrue(t, check, fmt.Sprintf("Expected Agent not found in related agents list %s", id))
	}
}

func TestGetRelatedAgentsReturnDistributedResults(t *testing.T) {
	json := `{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"},{"id":"id_4"},{"id":"id_5"}],"links":[{"source":"id_1","target":"id_2"},{"source":"id_1","target":"id_3"},{"source":"id_4","target":"id_1"},{"source":"id_5","target":"id_1"}]}`
	n, err := NewNetwork(json)
	AssertSuccess(t, err)
	agentCounts := make(map[string]int, 4)
	iterations := 2000
	for i := 0; i < iterations; i++ {
		relatedAgents := n.GetRelatedAgents(n.AgentsByID["id_1"])
		agentCounts[relatedAgents[0].Identifier()]++
	}
	avgCount := iterations / 4
	maxCount := avgCount + avgCount/10
	minCount := avgCount - avgCount/10
	for id, count := range agentCounts {
		IsTrue(t, count >= minCount && count <= maxCount, fmt.Sprintf("Repeated calls to GetRelatedAgents does not distribute agents evenly. ID '%s' should be around %d but was %d", id, avgCount, count))
	}
}

func TestGetAgentByIDReturnsCorrectAgentWhenExists(t *testing.T) {
	json := `{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"},{"id":"id_4"},{"id":"id_5"}],"links":[{"source":"id_1","target":"id_2"},{"source":"id_1","target":"id_3"},{"source":"id_4","target":"id_1"},{"source":"id_5","target":"id_1"}]}`
	n, err := NewNetwork(json)
	AssertSuccess(t, err)
	agent := n.GetAgentByID("id_5")
	AreEqual(t, "id_5", agent.Identifier(), "Expected agent to be id_5")
}

func TestGetAgentByIDReturnsNilWhenNotExists(t *testing.T) {
	json := `{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"},{"id":"id_4"},{"id":"id_5"}],"links":[{"source":"id_1","target":"id_2"},{"source":"id_1","target":"id_3"},{"source":"id_4","target":"id_1"},{"source":"id_5","target":"id_1"}]}`
	n, err := NewNetwork(json)
	AssertSuccess(t, err)
	agent := n.GetAgentByID("id_9")
	AreEqual(t, (Agent)(nil), agent, "Expected nil to be returned")
}

func TestUpdateLinkStrength(t *testing.T) {
	json := `{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"},{"id":"id_4"},{"id":"id_5"}],"links":[{"source":"id_1","target":"id_2"},{"source":"id_1","target":"id_3"},{"source":"id_4","target":"id_1"},{"source":"id_5","target":"id_1"}]}`
	n, err := NewNetwork(json)
	AssertSuccess(t, err)
	err = n.IncrementLinkStrength("id_5", "id_1")
	AssertSuccess(t, err)
	AreEqual(t, 1, n.AgentLinkMap["id_5"]["id_1"].Link.Strength, "Link strength not updated")
	AreEqual(t, 1, n.AgentLinkMap["id_1"]["id_5"].Link.Strength, "Link strength not updated")
	err = n.IncrementLinkStrength("id_1", "id_5")
	AssertSuccess(t, err)
	AreEqual(t, 2, n.AgentLinkMap["id_5"]["id_1"].Link.Strength, "Link strength not updated")
	AreEqual(t, 2, n.AgentLinkMap["id_1"]["id_5"].Link.Strength, "Link strength not updated")
}

func TestUpdateLinkStrengthReportsErrorIfIdDoesntExist(t *testing.T) {
	json := `{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"},{"id":"id_4"},{"id":"id_5"}],"links":[{"source":"id_1","target":"id_2"},{"source":"id_1","target":"id_3"},{"source":"id_4","target":"id_1"},{"source":"id_5","target":"id_1"}]}`
	n, err := NewNetwork(json)
	AssertSuccess(t, err)

	err = n.IncrementLinkStrength("id_7", "id_1")
	NotEqual(t, nil, err, "Invalid link error not reported")
	if !strings.Contains(err.Error(), "id_7") || !strings.Contains(err.Error(), "id_1") {
		t.Errorf("Incorrect error reported: %s", err.Error())
	}

	err = n.IncrementLinkStrength("id_1", "id_7")
	NotEqual(t, nil, err, "Invalid link error not reported")
	if !strings.Contains(err.Error(), "id_7") || !strings.Contains(err.Error(), "id_1") {
		t.Errorf("Incorrect error reported: %s", err.Error())
	}
}
