package srvr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/codeafix/orgnetsim/sim"
	"github.com/google/uuid"
	"github.com/spaceweasel/mango"
)

func CreateStepHandlerTestRouter(tfu *TestFileUpdater) *mango.Browser {
	r := CreateRouter(NewTestFileManager(tfu))
	// r.RequestLogger = func(l *mango.RequestLog) {
	// 	fmt.Println(l.CombinedFormat())
	// }
	// r.ErrorLogger = func(err error) {
	// 	fmt.Println(err.Error())
	// }

	br := mango.NewBrowser(r)

	return br
}

func TestReadSimStep(t *testing.T) {
	ts := &SimStep{
		ID:       uuid.New().String(),
		ParentID: uuid.New().String(),
	}
	tfu := &TestFileUpdater{
		Obj:      ts,
		Filepath: ts.Filepath(),
	}
	br := CreateStepHandlerTestRouter(tfu)

	hdrs := http.Header{}
	resp, err := br.Get(fmt.Sprintf("/api/simulation/%s/step/%s", ts.ParentID, ts.ID), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusOK, resp.Code, "Not OK")

	rs := &SimStep{}
	json.Unmarshal(resp.Body.Bytes(), rs)
	AreEqual(t, ts.ID, rs.ID, "Wrong object returned")
	AreEqual(t, ts.ParentID, rs.ParentID, "Wrong parent returned")
}

func TestReadSimStepError(t *testing.T) {
	ts := &SimStep{
		ID:       uuid.New().String(),
		ParentID: uuid.New().String(),
	}
	tfu := &TestFileUpdater{
		Obj:      ts,
		ReadErr:  fmt.Errorf("File not found"),
		Filepath: ts.Filepath(),
	}
	br := CreateStepHandlerTestRouter(tfu)

	hdrs := http.Header{}
	resp, err := br.Get(fmt.Sprintf("/api/simulation/%s/step/%s", ts.ParentID, ts.ID), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusInternalServerError, resp.Code, "Not Error")
}

func TestUpdateSimStepWithInvalidJson(t *testing.T) {
	ts := &SimStep{
		ID:       uuid.New().String(),
		ParentID: uuid.New().String(),
	}
	tfu := &TestFileUpdater{
		Obj:      ts,
		Filepath: ts.Filepath(),
	}
	br := CreateStepHandlerTestRouter(tfu)

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	data := `rg4t34to","code":"MGOdfbdb","categoryit"vcvbeaer}`
	resp, err := br.PutS(fmt.Sprintf("/api/simulation/%s/step/%s", ts.ParentID, ts.ID), data, hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusBadRequest, resp.Code, "Not Bad request")
}

func TestUpdateSimStepWithWrongHeaders(t *testing.T) {
	ts := &SimStep{
		ID:       uuid.New().String(),
		ParentID: uuid.New().String(),
	}
	tfu := &TestFileUpdater{
		Obj:      ts,
		Filepath: ts.Filepath(),
	}
	br := CreateStepHandlerTestRouter(tfu)

	hdrs := http.Header{}
	data := `{"name":"mango","code":"MGO","category":"fruit"}`
	resp, err := br.PutS(fmt.Sprintf("/api/simulation/%s/step/%s", ts.ParentID, ts.ID), data, hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusBadRequest, resp.Code, "Not Bad request")
}

func TestUpdateSimStepUpdatesCorrectly(t *testing.T) {
	ts := &SimStep{
		ID:       uuid.New().String(),
		ParentID: uuid.New().String(),
	}
	tfu := &TestFileUpdater{
		Obj:      ts,
		Filepath: ts.Filepath(),
	}
	br := CreateStepHandlerTestRouter(tfu)

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	data := fmt.Sprintf(`{"network":null,"results":{"iterations":5,"colors":null,"conversations":null},"id":"%s","parent":"%s"}`, ts.ID, ts.ParentID)
	resp, err := br.PutS(fmt.Sprintf("/api/simulation/%s/step/%s", ts.ParentID, ts.ID), data, hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusOK, resp.Code, "Not OK")

	rs := &SimStep{}
	json.Unmarshal(resp.Body.Bytes(), rs)
	AreEqual(t, ts.ID, rs.ID, "Wrong object returned")
	AreEqual(t, ts.ParentID, rs.ParentID, "Wrong parent returned")
	AreEqual(t, 5, rs.Results.Iterations, "Wrong iterations count returned")
	AreEqual(t, 5, tfu.Obj.(*SimStep).Results.Iterations, "Wrong iterations count written")
}

func TestUpdateSimStepReturnsErrorWithFileReadErr(t *testing.T) {
	ts := &SimStep{
		ID:       uuid.New().String(),
		ParentID: uuid.New().String(),
	}
	tfu := &TestFileUpdater{
		Obj:      ts,
		ReadErr:  fmt.Errorf("Access denied"),
		Filepath: ts.Filepath(),
	}
	br := CreateStepHandlerTestRouter(tfu)

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	data := fmt.Sprintf(`{"network":null,"results":{"iterations":5,"colors":null,"conversations":null},"id":"%s","parent":"%s"}`, ts.ID, ts.ParentID)
	resp, err := br.PutS(fmt.Sprintf("/api/simulation/%s/step/%s", ts.ParentID, ts.ID), data, hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusInternalServerError, resp.Code, "No error reported")
	Contains(t, "Access denied", resp.Body.String(), "Wrong error reported")
}

func TestUpdateSimStepReturnsErrorWithFileUpdateErr(t *testing.T) {
	ts := &SimStep{
		ID:       uuid.New().String(),
		ParentID: uuid.New().String(),
	}
	tfu := &TestFileUpdater{
		Obj:       ts,
		UpdateErr: fmt.Errorf("Access denied"),
		Filepath:  ts.Filepath(),
	}
	br := CreateStepHandlerTestRouter(tfu)

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	data := fmt.Sprintf(`{"network":null,"results":{"iterations":5,"colors":null,"conversations":null},"id":"%s","parent":"%s"}`, ts.ID, ts.ParentID)
	resp, err := br.PutS(fmt.Sprintf("/api/simulation/%s/step/%s", ts.ParentID, ts.ID), data, hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusInternalServerError, resp.Code, "No error reported")
	Contains(t, "Access denied", resp.Body.String(), "Wrong error reported")
}

func TestGetNetworkForStepSuccess(t *testing.T) {
	// Use CreateSimHandlerBrowserWithSteps from simhandler_test.go
	// It sets up a SimInfo and a SimStep (ids[0]) with a network
	br, _, ssfu, _, _, simid := CreateSimHandlerBrowserWithSteps(0)
	mockStep := ssfu.Obj.(*SimStep) // The SimStep created by the helper

	hdrs := http.Header{}
	resp, err := br.Get(fmt.Sprintf("/api/simulation/%s/step/%s/network", simid, mockStep.ID), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusOK, resp.Code, "Not OK")

	rNet := &sim.Network{}
	err = json.Unmarshal(resp.Body.Bytes(), rNet)
	AssertSuccess(t, err)
	IsFalse(t, rNet == nil, "Returned network should not be nil")
	// Based on CreateNetwork() in simhandler_test.go, there should be 3 agents
	AreEqual(t, 3, len(rNet.Agents()), "Wrong number of agents in returned network")
	// Add more assertions based on the expected network structure from CreateNetwork()
	AreEqual(t, 4, rNet.MaxColors(), "MaxColors incorrect")
	AreEqual(t, 2, len(rNet.Links()), "Number of links incorrect")
}

func TestGetNetworkForStepNotErrorWithFileRead(t *testing.T) {
	ts := &SimStep{
		ID:       uuid.New().String(),
		ParentID: uuid.New().String(),
	}
	tfu := &TestFileUpdater{
		Obj:      ts,
		ReadErr:  fmt.Errorf("Not Found"),
		Filepath: ts.Filepath(),
	}
	br := CreateStepHandlerTestRouter(tfu)

	hdrs := http.Header{}
	resp, err := br.Get(fmt.Sprintf("/api/simulation/%s/step/%s/network", ts.ParentID, ts.ID), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusInternalServerError, resp.Code, "Not NotFound for non-existent step")
}

func TestGetAgentColorsForStepSuccess(t *testing.T) {
	br, _, ssfu, _, _, simid := CreateSimHandlerBrowserWithSteps(0)
	mockStep := ssfu.Obj.(*SimStep)
	// Explicitly set some results for this step, as CreateSimHandlerBrowserWithSteps might not
	mockStep.Results = sim.Results{Iterations: 10, Colors: [][]int{{1, 2, 3}, {4, 5, 6}}, Conversations: []int{5, 5}}
	ssfu.Obj = mockStep // Update the object in the mock file updater as it's a pointer

	hdrs := http.Header{}
	resp, err := br.Get(fmt.Sprintf("/api/simulation/%s/step/%s/agents", simid, mockStep.ID), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusOK, resp.Code, "Not OK")

	rResults := []*sim.AgentState{}
	err = json.Unmarshal(resp.Body.Bytes(), &rResults)
	AssertSuccess(t, err)
	AreEqual(t, 3, len(rResults), "Wrong number of agents")
	AreEqual(t, sim.Blue, rResults[0].GetColor(), "Wrong color for agent 1")
	AreEqual(t, sim.Blue, rResults[1].GetColor(), "Wrong color for agent 2")
	AreEqual(t, sim.Blue, rResults[2].GetColor(), "Wrong color for agent 3")
}

func TestGetAgentColorsForStepNotFound(t *testing.T) {
	{
		ts := &SimStep{
			ID:       uuid.New().String(),
			ParentID: uuid.New().String(),
		}
		tfu := &TestFileUpdater{
			Obj:      ts,
			ReadErr:  fmt.Errorf("Not Found"),
			Filepath: ts.Filepath(),
		}
		br := CreateStepHandlerTestRouter(tfu)

		hdrs := http.Header{}
		resp, err := br.Get(fmt.Sprintf("/api/simulation/%s/step/%s/agents", ts.ParentID, ts.ID), hdrs)
		AssertSuccess(t, err)
		AreEqual(t, http.StatusInternalServerError, resp.Code, "Not NotFound for non-existent step")
	}
}

func TestPutStepNetworkData_Success(t *testing.T) {
	simID := uuid.New().String()
	stepID := uuid.New().String()

	// Initial SimStep in our mock file system
	initialStep := &SimStep{
		ID:       stepID,
		ParentID: simID,
		Network:  &sim.Network{}, // Start with an empty network
	}
	tfu := &TestFileUpdater{
		Obj:      initialStep,
		Filepath: initialStep.Filepath(),
	}
	br := CreateStepHandlerTestRouter(tfu)

	// New network data to send
	newNetwork := &sim.Network{
		Nodes: []sim.Agent{
			&sim.AgentState{ID: "agent1", X: 10, Y: 20, Color: 1},
			&sim.AgentState{ID: "agent2", X: 30, Y: 40, Color: 2},
		},
		Edges: []*sim.Link{
			{Agent1ID: "agent1", Agent2ID: "agent2", Strength: 15},
		},
		MaxColorCount: 2,
	}
	newNetwork.PopulateMaps() // Ensure internal maps are populated before comparison

	networkData, err := json.Marshal(newNetwork)
	AssertSuccess(t, err)

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")

	resp, err := br.PutS(fmt.Sprintf("/api/simulation/%s/step/%s/network", simID, stepID), string(networkData), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusOK, resp.Code, "Expected StatusOK")

	// Check response body
	returnedNetwork := &sim.Network{}
	err = json.Unmarshal(resp.Body.Bytes(), returnedNetwork)
	AssertSuccess(t, err)
	AreEqual(t, len(newNetwork.Nodes), len(returnedNetwork.Nodes), "Returned network node count mismatch")
	AreEqual(t, len(newNetwork.Edges), len(returnedNetwork.Edges), "Returned network edge count mismatch")
	if len(newNetwork.Nodes) > 0 && len(returnedNetwork.Nodes) > 0 {
		AreEqual(t, newNetwork.Nodes[0].Identifier(), returnedNetwork.Nodes[0].Identifier(), "Returned network node ID mismatch")
	}

	// Check that the file updater's object was updated
	updatedStep, ok := tfu.Obj.(*SimStep)
	IsTrue(t, ok, "TestFileUpdater object is not a SimStep")
	IsFalse(t, updatedStep.Network == nil, "Updated step network should not be nil")
	AreEqual(t, len(newNetwork.Nodes), len(updatedStep.Network.Agents()), "Persisted network node count mismatch")
	AreEqual(t, len(newNetwork.Edges), len(updatedStep.Network.Links()), "Persisted network edge count mismatch")
	if len(newNetwork.Nodes) > 0 && len(updatedStep.Network.Agents()) > 0 {
		AreEqual(t, newNetwork.Nodes[0].Identifier(), updatedStep.Network.Agents()[0].Identifier(), "Persisted network node ID mismatch")
	}
}

func TestPutStepNetworkData_ReadError(t *testing.T) {
	simID := uuid.New().String()
	stepID := uuid.New().String()
	initialStep := &SimStep{ID: stepID, ParentID: simID} // Only needed for filepath

	tfu := &TestFileUpdater{
		ReadErr:  fmt.Errorf("Simulated read error"),
		Filepath: initialStep.Filepath(), // Filepath still needed for Get
	}
	br := CreateStepHandlerTestRouter(tfu)

	newNetwork := &sim.Network{Nodes: []sim.Agent{&sim.AgentState{ID: "agent1"}}}
	networkData, _ := json.Marshal(newNetwork)
	hdrs := http.Header{"Content-Type": {"application/json"}}

	resp, err := br.PutS(fmt.Sprintf("/api/simulation/%s/step/%s/network", simID, stepID), string(networkData), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusInternalServerError, resp.Code, "Expected InternalServerError on read error")
	Contains(t, "Error reading step:", resp.Body.String(), "Response body does not contain expected error message for read error")
}

func TestPutStepNetworkData_BindError(t *testing.T) {
	simID := uuid.New().String()
	stepID := uuid.New().String()
	initialStep := &SimStep{ID: stepID, ParentID: simID, Network: &sim.Network{}}
	tfu := &TestFileUpdater{Obj: initialStep, Filepath: initialStep.Filepath()}
	br := CreateStepHandlerTestRouter(tfu)

	invalidNetworkData := "this is not json"
	hdrs := http.Header{"Content-Type": {"application/json"}}

	resp, err := br.PutS(fmt.Sprintf("/api/simulation/%s/step/%s/network", simID, stepID), invalidNetworkData, hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusBadRequest, resp.Code, "Expected BadRequest on bind error")
	Contains(t, "Error binding network data", resp.Body.String(), "Response body does not contain expected error message for bind error")
}

func TestPutStepNetworkData_UpdateError(t *testing.T) {
	simID := uuid.New().String()
	stepID := uuid.New().String()
	initialStep := &SimStep{ID: stepID, ParentID: simID, Network: &sim.Network{}}
	tfu := &TestFileUpdater{
		Obj:       initialStep,
		Filepath:  initialStep.Filepath(),
		UpdateErr: fmt.Errorf("simulated update error"),
	}
	br := CreateStepHandlerTestRouter(tfu)

	newNetwork := &sim.Network{Nodes: []sim.Agent{&sim.AgentState{ID: "agent1"}}}
	networkData, _ := json.Marshal(newNetwork)
	hdrs := http.Header{"Content-Type": {"application/json"}}

	resp, err := br.PutS(fmt.Sprintf("/api/simulation/%s/step/%s/network", simID, stepID), string(networkData), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusInternalServerError, resp.Code, "Expected InternalServerError on update error")
	Contains(t, "Error updating step with new network", resp.Body.String(), "Response body does not contain expected error message for update error")
}

func TestPutStepNetworkData_MissingIDs(t *testing.T) {
	br := CreateStepHandlerTestRouter(&TestFileUpdater{}) // TFU doesn't matter much here

	newNetwork := &sim.Network{Nodes: []sim.Agent{&sim.AgentState{ID: "agent1"}}}
	networkData, _ := json.Marshal(newNetwork)
	hdrs := http.Header{"Content-Type": {"application/json"}}

	// Missing sim_id
	resp, err := br.PutS(fmt.Sprintf("/api/simulation//step/%s/network", uuid.New().String()), string(networkData), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusBadRequest, resp.Code, "Expected BadRequest for missing sim_id")
	Contains(t, "Missing IDs in route", resp.Body.String(), "Response body does not contain expected error message for missing sim_id")

	// Missing step_id
	resp, err = br.PutS(fmt.Sprintf("/api/simulation/%s/step//network", uuid.New().String()), string(networkData), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusBadRequest, resp.Code, "Expected BadRequest for missing step_id")
	Contains(t, "Missing IDs in route", resp.Body.String(), "Response body does not contain expected error message for missing step_id")
}

func TestPutStepData_InvalidDataType(t *testing.T) {
	br := CreateStepHandlerTestRouter(&TestFileUpdater{}) // TFU doesn't matter much here

	newNetwork := &sim.Network{Nodes: []sim.Agent{&sim.AgentState{ID: "agent1"}}}
	networkData, _ := json.Marshal(newNetwork)
	hdrs := http.Header{"Content-Type": {"application/json"}}

	// Unsupported data type
	resp, err := br.PutS(fmt.Sprintf("/api/simulation/%s/step/%s/invalid", uuid.New().String(), uuid.New().String()), string(networkData), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusBadRequest, resp.Code, "Expected BadRequest for invalid data type")
	Contains(t, "Only direct updates to 'network' are available.", resp.Body.String(), "Response body does not contain expected error message for missing sim_id")
}
