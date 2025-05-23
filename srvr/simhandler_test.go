package srvr

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/codeafix/orgnetsim/sim"
	"github.com/google/uuid"
	"github.com/spaceweasel/mango"
)

func CreateNetwork() sim.RelationshipMgr {
	rm := &sim.Network{}
	agent1 := sim.GenerateRandomAgent("Agent_1", "Agent 1", []sim.Color{sim.Blue}, false)
	rm.AddAgent(agent1)
	agent2 := sim.GenerateRandomAgent("Agent_2", "Agent 2", []sim.Color{sim.Blue}, false)
	rm.AddAgent(agent2)
	agent3 := sim.GenerateRandomAgent("Agent_3", "Agent 3", []sim.Color{sim.Blue}, false)
	rm.AddAgent(agent3)
	rm.AddLink(agent1, agent2)
	rm.AddLink(agent1, agent3)
	rm.SetMaxColors(4)
	rm.PopulateMaps()
	return rm
}

func CreateSimHandlerBrowserWithSteps(deleteItemIndex int) (*mango.Browser, *TestFileUpdater, *TestFileUpdater, *TestFileUpdater, []string, string) {
	simid := uuid.New().String()
	sim := NewSimInfo(simid)
	sim.Name = "mySavedSim"
	sim.Description = "A description of mySavedSim"
	ids := []string{
		uuid.NewString(),
		uuid.NewString(),
		uuid.NewString(),
	}
	steps := make([]string, len(ids))
	for i, id := range ids {
		steps[i] = fmt.Sprintf("/api/simulation/%s/step/%s", simid, id)
	}
	sim.Steps = steps
	simfu := &TestFileUpdater{
		Obj:      sim,
		Filepath: sim.Filepath(),
	}
	tfm := NewTestFileManager(simfu)

	ss := &SimStep{
		ID:       ids[deleteItemIndex],
		ParentID: simid,
		Network:  CreateNetwork(),
	}
	ssfu := &TestFileUpdater{
		Obj:      ss,
		Filepath: ss.Filepath(),
	}
	tfm.Add(ss.Filepath(), ssfu)
	dfu := &TestFileUpdater{}
	tfm.Default = dfu

	r := CreateRouter(tfm)
	br := mango.NewBrowser(r)

	return br, simfu, ssfu, dfu, steps, simid
}

func TestMarshalling(t *testing.T) {
	parent := uuid.NewString()
	s1 := &SimStep{
		ID:       uuid.NewString(),
		ParentID: parent,
		Network:  CreateNetwork(),
	}

	s2 := &SimStep{
		ID:       uuid.NewString(),
		ParentID: parent,
		Network:  CreateNetwork(),
	}
	steps := []*SimStep{s1, s2}
	b, err := json.Marshal(steps)
	AssertSuccess(t, err)
	rstep := []*SimStep{}
	err = json.Unmarshal(b, &rstep)
	AssertSuccess(t, err)
}

func TestGetSimSuccess(t *testing.T) {
	br, simfu, _, _, _, simid := CreateSimHandlerBrowserWithSteps(0)

	hdrs := http.Header{}
	resp, err := br.Get(fmt.Sprintf("/api/simulation/%s", simid), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusOK, resp.Code, "Not OK")

	rsim := &SimInfo{}
	err = json.Unmarshal(resp.Body.Bytes(), rsim)
	AssertSuccess(t, err)
	AreEqual(t, simfu.Obj.(*SimInfo).Name, rsim.Name, "Wrong name in returned SimInfo")
	AreEqual(t, simfu.Obj.(*SimInfo).Description, rsim.Description, "Wrong description in returned SimInfo")
}

func TestGetSimInvalidCommand(t *testing.T) {
	br, _, _, _, _, simid := CreateSimHandlerBrowserWithSteps(0)

	hdrs := http.Header{}
	resp, err := br.Get(fmt.Sprintf("/api/simulation/%s/somename", simid), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusNotFound, resp.Code, "Did not return Not Found")
}

func TestGetSimStepsSuccess(t *testing.T) {
	br, _, _, _, steps, simid := CreateSimHandlerBrowserWithSteps(0)

	hdrs := http.Header{}
	resp, err := br.Get(fmt.Sprintf("/api/simulation/%s/step", simid), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusOK, resp.Code, "Not OK")

	rsteps := []SimStepSummary{}
	err = json.Unmarshal(resp.Body.Bytes(), &rsteps)

	AssertSuccess(t, err)
	AreEqual(t, 3, len(rsteps), "Wrong number of Steps in returned SimInfo")
	// The 'steps' variable contains full paths, extract ID for comparison
	// Example step path: /api/simulation/SIM_ID/step/STEP_ID
	expectedId0 := steps[0][strings.LastIndex(steps[0], "/")+1:]
	AreEqual(t, expectedId0, rsteps[0].ID, "Wrong Step 0 ID")
	// Add ParentID assertion if needed, e.g.: AreEqual(t, simid, rsteps[0].ParentID, "Wrong ParentID for Step 0")

	if len(rsteps) > 1 {
		expectedId1 := steps[1][strings.LastIndex(steps[1], "/")+1:]
		AreEqual(t, expectedId1, rsteps[1].ID, "Wrong Step 1 ID")
	}
	if len(rsteps) > 2 {
		expectedId2 := steps[2][strings.LastIndex(steps[2], "/")+1:]
		AreEqual(t, expectedId2, rsteps[2].ID, "Wrong Step 2 ID")
	}
}

func TestUpdateSimSuccess(t *testing.T) {
	br, _, _, _, steps, simid := CreateSimHandlerBrowserWithSteps(0)

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	data := `{"Name":"myUpdatedSim","Description":"A description of mySavedSim"}`
	resp, err := br.PutS(fmt.Sprintf("/api/simulation/%s", simid), data, hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusOK, resp.Code, "Not OK")

	rsim := &SimInfo{}
	err = json.Unmarshal(resp.Body.Bytes(), rsim)
	AssertSuccess(t, err)
	AreEqual(t, "myUpdatedSim", rsim.Name, "Wrong name in returned SimInfo")
	AreEqual(t, "A description of mySavedSim", rsim.Description, "Wrong description in returned SimInfo")
	AreEqual(t, 3, len(rsim.Steps), "Wrong number of Steps in returned SimInfo")
	AreEqual(t, steps[0], rsim.Steps[0], "Wrong Step 0 in returned SimInfo")
	AreEqual(t, steps[1], rsim.Steps[1], "Wrong Step 1 in returned SimInfo")
	AreEqual(t, steps[2], rsim.Steps[2], "Wrong Step 2 in returned SimInfo")
}

func TestDeleteSimSuccess(t *testing.T) {
	br, simfu, ssfu, _, steps, _ := CreateSimHandlerBrowserWithSteps(0)

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	resp, err := br.Delete(steps[0], hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusOK, resp.Code, "Not OK")

	AreEqual(t, 2, len(simfu.Obj.(*SimInfo).Steps), "There should be two items in the list")
	AreEqual(t, steps[1], simfu.Obj.(*SimInfo).Steps[0], "Wrong path in position 0 of step list")
	AreEqual(t, steps[2], simfu.Obj.(*SimInfo).Steps[1], "Wrong path in position 1 of step list")
	IsTrue(t, ssfu.DeleteCalled, "Delete was not called on the correct fileupdater")
}

func TestGenerateNetworkFailsIfStepsExist(t *testing.T) {
	br, _, _, _, _, simid := CreateSimHandlerBrowserWithSteps(0)

	hs := sim.HierarchySpec{
		Levels:     2,
		TeamSize:   2,
		InitColors: []sim.Color{sim.Blue},
		MaxColors:  4,
	}
	hss, err := json.Marshal(hs)
	AssertSuccess(t, err)

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	resp, err := br.PostS(fmt.Sprintf("/api/simulation/%s/generate", simid), string(hss), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusBadRequest, resp.Code, "Not Bad Request")
	AreEqual(t, "Simulation must have no steps when generating a new network", strings.TrimSpace(resp.Body.String()), "Incorrect error response")
}

func TestParseNetworkFailsIfStepsExist(t *testing.T) {
	br, _, _, _, _, simid := CreateSimHandlerBrowserWithSteps(0)

	data := []string{
		"Header row is always skipped ,check_this_is_not_an_Id,,\n",
		"Should be ignored|||\n",
		"\n",
		"Strips ws around Id, my_id\n",
		"Blank lines are ignored\n",
	}
	var payload = []byte{}
	for _, s := range data {
		payload = append(payload, []byte(s)...)
	}
	pb := ParseBody{
		ParseOptions: sim.ParseOptions{
			Delimiter:  ",",
			Identifier: 1,
			Parent:     3,
		},
		Payload: payload,
	}

	pbs, err := json.Marshal(pb)
	AssertSuccess(t, err)

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	resp, err := br.PostS(fmt.Sprintf("/api/simulation/%s/parse", simid), string(pbs), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusBadRequest, resp.Code, "Not Bad Request")
	AreEqual(t, "Simulation must have no steps when parsing a new network", strings.TrimSpace(resp.Body.String()), "Incorrect error response")
}

func CreateSimHandlerBrowser() (*mango.Browser, *TestFileUpdater, *TestFileUpdater, string) {
	simid := uuid.New().String()
	nsim := NewSimInfo(simid)
	nsim.Name = "mySavedSim"
	nsim.Description = "A description of mySavedSim"
	nsim.Steps = []string{}
	nsim.Options.LinkedTeamList = []string{}
	nsim.Options.EvangelistList = []string{}
	nsim.Options.LoneEvangelist = []string{}
	nsim.Options.InitColors = []sim.Color{}
	simfu := &TestFileUpdater{
		Obj:      nsim,
		Filepath: nsim.Filepath(),
	}
	tfm := NewTestFileManager(simfu)
	ssfu := &TestFileUpdater{}
	tfm.SetDefault(ssfu)

	r := CreateRouter(tfm)
	br := mango.NewBrowser(r)

	return br, simfu, ssfu, simid
}

func TestParseNetworkSucceeds(t *testing.T) {
	br, simfu, ssfu, simid := CreateSimHandlerBrowser()
	savedsim, ok := simfu.Obj.(*SimInfo)
	IsTrue(t, ok, "Saved object would not cast to *SimInfo")
	savedsim.Options.MaxColors = 5

	data := []string{
		"Header always skipped ,check_this_is_not_an_Id\n",
		"Should be ignored|||\n",
		"\n",
		"Strips ws around Id, my_id\n",
		"Blank lines are ignored\n",
		"First agent, agent_1, some text,,\n",
		"Second agent, agent_2, more text, agent_1,\n",
	}
	var payload = []byte{}
	for _, s := range data {
		payload = append(payload, []byte(s)...)
	}
	pb := ParseBody{
		ParseOptions: sim.ParseOptions{
			Delimiter:  ",",
			Identifier: 1,
			Parent:     3,
		},
		Payload: payload,
	}

	pbs, err := json.Marshal(pb)
	AssertSuccess(t, err)

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	resp, err := br.PostS(fmt.Sprintf("/api/simulation/%s/parse", simid), string(pbs), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusCreated, resp.Code, "Not Created")
	simstep, ok := ssfu.Obj.(*SimStep)
	IsTrue(t, ok, "Saved object would not cast to *SimStep")
	AreEqual(t, 5, simstep.Network.MaxColors(), "Wrong MaxColors on network")
	IsTrue(t, simstep.Network.Links() != nil, "Links array is nil")
	AreEqual(t, len(simstep.Network.Links()), 1, "Links should have a single item")
	IsTrue(t, simstep.Network.Agents() != nil, "Agents array is nil")
	AreEqual(t, len(simstep.Network.Agents()), 3, "Agents array should have 3 items")
}

func TestParseNetworkFailsWithNoPayload(t *testing.T) {
	br, _, _, simid := CreateSimHandlerBrowser()

	pb := ParseBody{
		ParseOptions: sim.ParseOptions{
			Delimiter:  ",",
			Identifier: 0,
			Parent:     1,
		},
		Payload: []byte{},
	}

	pbs, err := json.Marshal(pb)
	AssertSuccess(t, err)

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	resp, err := br.PostS(fmt.Sprintf("/api/simulation/%s/parse", simid), string(pbs), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusBadRequest, resp.Code, "Not Bad request")
	AreEqual(t, "No links data in ParseOptions", strings.TrimSpace(resp.Body.String()), "Incorrect error response")
}

func TestParseNetworkFailsWithNoParseOptions(t *testing.T) {
	br, _, _, simid := CreateSimHandlerBrowser()

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	resp, err := br.PostS(fmt.Sprintf("/api/simulation/%s/parse", simid), "", hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusBadRequest, resp.Code, "Not Bad request")
	AreEqual(t, "EOF: Error reading ParseOptions", strings.TrimSpace(resp.Body.String()), "Incorrect error response")
}

func TestAddLinksSucceeds(t *testing.T) {
	br, _, ssfu, _, _, simid := CreateSimHandlerBrowserWithSteps(2)
	ssfu.Obj.(*SimStep).Network.AddAgent(sim.GenerateRandomAgent("Agent_4", "Agent 4", []sim.Color{sim.Blue}, false))
	ssfu.Obj.(*SimStep).Network.AddAgent(sim.GenerateRandomAgent("Agent_5", "Agent 5", []sim.Color{sim.Blue}, false))

	data := []string{
		"Header always skipped ,check_this_is_not_an_Id\n",
		"Should be ignored|||\n",
		"\n",
		"Agent_2, Agent_4\n",
		"Agent_2, Agent_5\n",
	}
	var payload = []byte{}
	for _, s := range data {
		payload = append(payload, []byte(s)...)
	}
	pb := ParseBody{
		ParseOptions: sim.ParseOptions{
			Delimiter:  ",",
			Identifier: 0,
			Parent:     1,
		},
		Payload: payload,
	}

	pbs, err := json.Marshal(pb)
	AssertSuccess(t, err)

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	resp, err := br.PutS(fmt.Sprintf("/api/simulation/%s/links", simid), string(pbs), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusOK, resp.Code, "Not Updated")
	simstep, ok := ssfu.Obj.(*SimStep)
	IsTrue(t, ok, "Saved object would not cast to *SimStep")
	IsTrue(t, simstep.Network.Links() != nil, "Links array is nil")
	AreEqual(t, len(simstep.Network.Links()), 4, "Links should have 4 items")
	IsTrue(t, simstep.Network.Agents() != nil, "Agents array is nil")
	AreEqual(t, len(simstep.Network.Agents()), 5, "Agents array should have 5 items")
}

func TestAddLinksIgnoresLinksWhenAgentDoesntExist(t *testing.T) {
	br, _, ssfu, _, _, simid := CreateSimHandlerBrowserWithSteps(2)

	data := []string{
		"Header always skipped ,check_this_is_not_an_Id\n",
		"Should be ignored|||\n",
		"\n",
		"Agent_2, Agent_4\n",
	}
	var payload = []byte{}
	for _, s := range data {
		payload = append(payload, []byte(s)...)
	}
	pb := ParseBody{
		ParseOptions: sim.ParseOptions{
			Delimiter:  ",",
			Identifier: 0,
			Parent:     1,
		},
		Payload: payload,
	}

	pbs, err := json.Marshal(pb)
	AssertSuccess(t, err)

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	resp, err := br.PutS(fmt.Sprintf("/api/simulation/%s/links", simid), string(pbs), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusOK, resp.Code, "Not Updated")
	simstep, ok := ssfu.Obj.(*SimStep)
	IsTrue(t, ok, "Saved object would not cast to *SimStep")
	IsTrue(t, simstep.Network.Links() != nil, "Links array is nil")
	AreEqual(t, len(simstep.Network.Links()), 2, "Links should have 2 items")
	IsTrue(t, simstep.Network.Agents() != nil, "Agents array is nil")
	AreEqual(t, len(simstep.Network.Agents()), 3, "Agents array should have 3 items")
}

func TestAddLinksFailsIfNoStepsExist(t *testing.T) {
	br, _, _, simid := CreateSimHandlerBrowser()

	data := []string{
		"Header always skipped ,check_this_is_not_an_Id\n",
		"Should be ignored|||\n",
		"\n",
		"Agent_2, Agent_4\n",
		"Agent_2, Agent_5\n",
	}
	var payload = []byte{}
	for _, s := range data {
		payload = append(payload, []byte(s)...)
	}
	pb := ParseBody{
		ParseOptions: sim.ParseOptions{
			Delimiter:  ",",
			Identifier: 0,
			Parent:     1,
		},
		Payload: payload,
	}

	pbs, err := json.Marshal(pb)
	AssertSuccess(t, err)

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	resp, err := br.PutS(fmt.Sprintf("/api/simulation/%s/links", simid), string(pbs), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusBadRequest, resp.Code, "Not Bad request")
	AreEqual(t, "The network cannot have links added without an initial step containing a network", strings.TrimSpace(resp.Body.String()), "Incorrect error response")
}

func TestAddLinksFailsWithNoPayload(t *testing.T) {
	br, _, _, _, _, simid := CreateSimHandlerBrowserWithSteps(2)

	pb := ParseBody{
		ParseOptions: sim.ParseOptions{
			Delimiter:  ",",
			Identifier: 0,
			Parent:     1,
		},
		Payload: []byte{},
	}

	pbs, err := json.Marshal(pb)
	AssertSuccess(t, err)

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	resp, err := br.PutS(fmt.Sprintf("/api/simulation/%s/links", simid), string(pbs), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusBadRequest, resp.Code, "Not Bad request")
	AreEqual(t, "No links data in ParseOptions", strings.TrimSpace(resp.Body.String()), "Incorrect error response")
}

func TestAddLinksFailsWithNoParseOptions(t *testing.T) {
	br, _, _, _, _, simid := CreateSimHandlerBrowserWithSteps(2)

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	resp, err := br.PutS(fmt.Sprintf("/api/simulation/%s/links", simid), "", hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusBadRequest, resp.Code, "Not Bad request")
	AreEqual(t, "EOF: Error reading ParseOptions", strings.TrimSpace(resp.Body.String()), "Incorrect error response")
}

func CreateSimHandlerBrowserForCopyTests(stepcount int) (*mango.Browser, *TestFileUpdater, []*TestFileUpdater, string, *TestFileManager) {
	ssful := []*TestFileUpdater{}

	simid := uuid.New().String()
	nsim := NewSimInfo(simid)
	nsim.Name = "mySavedSim"
	nsim.Description = "A description of mySavedSim"
	nsim.Steps = []string{}
	simfu := &TestFileUpdater{
		Obj:      nsim,
		Filepath: nsim.Filepath(),
	}

	sl := NewSimList()
	sl.Items = []string{
		"/api/simulation/{someIdHere}",
	}
	sl.Notes = "Some notes"
	slfu := &TestFileUpdater{
		Obj:      sl,
		Filepath: sl.Filepath(),
	}

	tfm := NewTestFileManager(slfu)
	tfm.Add(simfu.Filepath, simfu)

	for i := 0; i < stepcount; i++ {
		rm := CreateNetwork()
		if i > 0 {
			rm.SetMaxColors(10)
		}
		ss := &SimStep{
			ID:       uuid.New().String(),
			ParentID: simid,
			Network:  rm,
		}
		ssfu := &TestFileUpdater{
			Obj:      ss,
			Filepath: ss.Filepath(),
		}
		nsim.Steps = append(nsim.Steps, ss.RelPath())
		tfm.Add(ss.Filepath(), ssfu)
		ssful = append(ssful, ssfu)
	}

	r := CreateRouter(tfm)
	br := mango.NewBrowser(r)

	return br, simfu, ssful, simid, tfm
}

func TestSimCopySucceedsWithNoSteps(t *testing.T) {
	br, simfu, _, simid, tfm := CreateSimHandlerBrowserForCopyTests(0)

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	resp, err := br.PostS(fmt.Sprintf("/api/simulation/%s/copy", simid), "", hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusCreated, resp.Code, "Not Created")
	siminfo, _ := simfu.Obj.(*SimInfo)

	cpsimfu := tfm.CreatedFileUpdaters(0)
	cpsim, ok := cpsimfu.Obj.(*SimInfo)
	IsTrue(t, ok, "Saved object would not cast to *SimInfo")
	AreEqual(t, cpsim.Name, "mySavedSim(copy)", "Wrong name in returned SimInfo")
	AreEqual(t, cpsim.Description, "A description of mySavedSim This is a copy of \"mySavedSim\".", "Wrong description in returned SimInfo")
	IsFalse(t, cpsim.ID == siminfo.ID, "ID should be different")
}

func TestSimCopySucceedsWithInitialStep(t *testing.T) {
	br, simfu, ssful, simid, tfm := CreateSimHandlerBrowserForCopyTests(1)

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	resp, err := br.PostS(fmt.Sprintf("/api/simulation/%s/copy", simid), "", hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusCreated, resp.Code, "Not Created")
	siminfo, _ := simfu.Obj.(*SimInfo)

	cpsimfu := tfm.CreatedFileUpdaters(0)
	cpsim, ok := cpsimfu.Obj.(*SimInfo)
	IsTrue(t, ok, "Saved object would not cast to *SimInfo")
	AreEqual(t, cpsim.Name, "mySavedSim(copy)", "Wrong name in returned SimInfo")
	AreEqual(t, cpsim.Description, "A description of mySavedSim This is a copy of \"mySavedSim\".", "Wrong description in returned SimInfo")
	IsFalse(t, cpsim.ID == siminfo.ID, "ID should be different")

	cpssfu := tfm.CreatedFileUpdaters(1)
	cpss, ok := cpssfu.Obj.(*SimStep)
	IsTrue(t, ok, "Saved object would not cast to *SimStep")
	AreEqual(t, cpss.ParentID, cpsim.ID, "ParentID should be the new sim ID")
	AreEqual(t, cpss.Network.MaxColors(), 4, "MaxColors should be the same")
	AreEqual(t, len(cpss.Network.Agents()), 3, "Agents should be the same")
	AreEqual(t, len(cpss.Network.Links()), 2, "Links should be the same")
	IsFalse(t, cpss.ID == ssful[0].Obj.(*SimStep).ID, "ID should be different")
}

func TestSimCopySucceedsWithMultipleSteps(t *testing.T) {
	br, simfu, ssful, simid, tfm := CreateSimHandlerBrowserForCopyTests(3)

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	resp, err := br.PostS(fmt.Sprintf("/api/simulation/%s/copy", simid), "", hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusCreated, resp.Code, "Not Created")
	siminfo, _ := simfu.Obj.(*SimInfo)

	cpsimfu := tfm.CreatedFileUpdaters(0)
	cpsim, ok := cpsimfu.Obj.(*SimInfo)
	IsTrue(t, ok, "Saved object would not cast to *SimInfo")
	AreEqual(t, cpsim.Name, "mySavedSim(copy)", "Wrong name in returned SimInfo")
	AreEqual(t, cpsim.Description, "A description of mySavedSim This is a copy of \"mySavedSim\".", "Wrong description in returned SimInfo")
	IsFalse(t, cpsim.ID == siminfo.ID, "ID should be different")
	IsTrue(t, len(cpsim.Steps) == 1, "Steps should contain initial step only")

	cpssfu := tfm.CreatedFileUpdaters(1)
	cpss, ok := cpssfu.Obj.(*SimStep)
	IsTrue(t, ok, "Saved object would not cast to *SimStep")
	AreEqual(t, cpss.ParentID, cpsim.ID, "ParentID should be the new sim ID")
	AreEqual(t, cpss.Network.MaxColors(), 4, "MaxColors should be the same")
	AreEqual(t, len(cpss.Network.Agents()), 3, "Agents should be the same")
	AreEqual(t, len(cpss.Network.Links()), 2, "Links should be the same")
	IsFalse(t, cpss.ID == ssful[0].Obj.(*SimStep).ID, "ID should be different")
}

func TestGenerateNetworkSucceeds(t *testing.T) {
	br, simfu, ssfu, simid := CreateSimHandlerBrowser()

	hs := sim.HierarchySpec{
		Levels:     2,
		TeamSize:   2,
		InitColors: []sim.Color{sim.Green},
		MaxColors:  4,
	}
	hss, err := json.Marshal(hs)
	AssertSuccess(t, err)

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	resp, err := br.PostS(fmt.Sprintf("/api/simulation/%s/generate", simid), string(hss), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusCreated, resp.Code, "Not Created")
	simstep, ok := ssfu.Obj.(*SimStep)
	IsTrue(t, ok, "Saved object would not cast to *SimStep")
	AreEqual(t, 4, simstep.Network.MaxColors(), "Wrong MaxColors on network")
	AreEqual(t, 3, len(simstep.Network.Agents()), "Wrong number of agents on network")
	AreEqual(t, 4, len(simstep.Results.Colors[0]), "Wrong number of Colors in Color results array")
	AreEqual(t, 3, simstep.Results.Colors[0][3], "Wrong Green Color count in results array")
	sim, ok := simfu.Obj.(*SimInfo)
	IsTrue(t, ok, "Saved object would not cast to *SimInfo")
	AreEqual(t, 1, len(sim.Options.InitColors), "Wrong InitColors on sim options")
	AreEqual(t, hs.InitColors[0], sim.Options.InitColors[0], "Wrong InitColors on sim options")
	AreEqual(t, hs.MaxColors, sim.Options.MaxColors, "Wrong MaxColors on sim options")
}

func TestGenerateNetworkFailsWithNoHierarchySpec(t *testing.T) {
	br, _, _, simid := CreateSimHandlerBrowser()

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	resp, err := br.PostS(fmt.Sprintf("/api/simulation/%s/generate", simid), "", hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusBadRequest, resp.Code, "Not Bad request")
	AreEqual(t, "EOF: Error reading HierarchySpec", strings.TrimSpace(resp.Body.String()), "Incorrect error response")
}

func TestPostRunFailsIfNoStepsExist(t *testing.T) {
	br, _, _, simid := CreateSimHandlerBrowser()

	rs := RunSpec{
		Iterations: 5,
		Steps:      5,
	}
	rss, err := json.Marshal(rs)
	AssertSuccess(t, err)

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	resp, err := br.PostS(fmt.Sprintf("/api/simulation/%s/run", simid), string(rss), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusBadRequest, resp.Code, "Not Bad request")
	AreEqual(t, "The simulation cannot be run without an initial step containing a network", strings.TrimSpace(resp.Body.String()), "Incorrect error response")
}

func TestPostRunFailsWithZeroIterations(t *testing.T) {
	br, _, _, _, _, simid := CreateSimHandlerBrowserWithSteps(2)

	rs := RunSpec{
		Iterations: 0,
		Steps:      5,
	}
	rss, err := json.Marshal(rs)
	AssertSuccess(t, err)

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	resp, err := br.PostS(fmt.Sprintf("/api/simulation/%s/run", simid), string(rss), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusBadRequest, resp.Code, "Not Bad request")
	AreEqual(t, "Steps and Iterations cannot be zero", strings.TrimSpace(resp.Body.String()), "Incorrect error response")
}

func TestPostRunFailsWithZeroStepCount(t *testing.T) {
	br, _, _, _, _, simid := CreateSimHandlerBrowserWithSteps(2)

	rs := RunSpec{
		Iterations: 0,
		Steps:      5,
	}
	rss, err := json.Marshal(rs)
	AssertSuccess(t, err)

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	resp, err := br.PostS(fmt.Sprintf("/api/simulation/%s/run", simid), string(rss), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusBadRequest, resp.Code, "Not Bad request")
	AreEqual(t, "Steps and Iterations cannot be zero", strings.TrimSpace(resp.Body.String()), "Incorrect error response")
}

func TestPostRunFailsWithNoRunSpec(t *testing.T) {
	br, _, _, _, _, simid := CreateSimHandlerBrowserWithSteps(2)

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	resp, err := br.PostS(fmt.Sprintf("/api/simulation/%s/run", simid), "", hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusBadRequest, resp.Code, "Not Bad request")
	AreEqual(t, "EOF: Error reading RunSpec", strings.TrimSpace(resp.Body.String()), "Incorrect error response")
}

func TestPostRunSucceedsWithOneStep(t *testing.T) {
	br, _, _, dfu, _, simid := CreateSimHandlerBrowserWithSteps(2)
	rs := RunSpec{
		Iterations: 5,
		Steps:      1,
	}
	rss, err := json.Marshal(rs)
	AssertSuccess(t, err)

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	resp, err := br.PostS(fmt.Sprintf("/api/simulation/%s/run", simid), string(rss), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusCreated, resp.Code, "Not created")
	ns := dfu.Obj.(*SimStep)
	NotEqual(t, nil, ns, "New step is nil")
	AreEqual(t, 4, ns.Network.MaxColors(), "MaxColors not correct on new network")
	//Check there are 5 items in the results arrays, one for each iteration
	AreEqual(t, 5, len(ns.Results.Colors), "Wrong number of items in the Colors array")
	AreEqual(t, 5, len(ns.Results.Conversations), "Wrong number of items in the Conversations array")
}

func CreateResults(iterations, maxColors int) sim.Results {
	results := sim.Results{
		Iterations:    iterations,
		Colors:        make([][]int, iterations),
		Conversations: make([]int, iterations),
	}
	for i := 0; i < iterations; i++ {
		results.Conversations[i] = i + 1
		colorCounts := make([]int, maxColors)
		for j := 0; j < maxColors; j++ {
			colorCounts[j] = i + j
		}
		results.Colors[i] = colorCounts
	}
	return results
}

func CreateSimHandlerBrowserWithStepsAndResults() (*mango.Browser, string) {
	simid := uuid.New().String()
	sim := NewSimInfo(simid)
	sim.Name = "mySavedSim"
	sim.Description = "A description of mySavedSim"
	ids := []string{
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
	}
	simfu := &TestFileUpdater{
		Obj:      sim,
		Filepath: sim.Filepath(),
	}
	tfm := NewTestFileManager(simfu)
	steps := make([]string, len(ids))
	for i, id := range ids {
		steps[i] = fmt.Sprintf("/api/simulation/%s/step/%s", simid, id)
		ss := &SimStep{
			ID:       id,
			ParentID: simid,
			Results:  CreateResults(i+1, 4),
		}
		ssfu := &TestFileUpdater{
			Obj:      ss,
			Filepath: ss.Filepath(),
		}
		tfm.Add(ss.Filepath(), ssfu)
	}
	sim.Steps = steps
	r := CreateRouter(tfm)
	br := mango.NewBrowser(r)

	return br, simid
}

func TestGetResultsSucceeds(t *testing.T) {
	br, simid := CreateSimHandlerBrowserWithStepsAndResults()

	hdrs := http.Header{}
	resp, err := br.Get(fmt.Sprintf("/api/simulation/%s/results", simid), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusOK, resp.Code, "Not OK")
	rs := &sim.Results{}
	json.Unmarshal(resp.Body.Bytes(), rs)
	AreEqual(t, 6, rs.Iterations, "Wrong number of iterations")
	AreEqual(t, 1, rs.Conversations[0], "Wrong conversation count")
	AreEqual(t, 1, rs.Conversations[1], "Wrong conversation count")
	AreEqual(t, 2, rs.Conversations[2], "Wrong conversation count")
	AreEqual(t, 1, rs.Conversations[3], "Wrong conversation count")
	AreEqual(t, 2, rs.Conversations[4], "Wrong conversation count")
	AreEqual(t, 3, rs.Conversations[5], "Wrong conversation count")
	AreEqual(t, 3, rs.Colors[0][3], "Wrong color count")
	AreEqual(t, 3, rs.Colors[1][3], "Wrong color count")
	AreEqual(t, 4, rs.Colors[2][3], "Wrong color count")
	AreEqual(t, 3, rs.Colors[3][3], "Wrong color count")
	AreEqual(t, 4, rs.Colors[4][3], "Wrong color count")
	AreEqual(t, 5, rs.Colors[5][3], "Wrong color count")
}

func TestGetResultsSucceedsAsCsv(t *testing.T) {
	br, simid := CreateSimHandlerBrowserWithStepsAndResults()

	hdrs := http.Header{
		"Content-Type": []string{"text/csv"},
	}
	resp, err := br.Get(fmt.Sprintf("/api/simulation/%s/results", simid), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusOK, resp.Code, "Not OK")
	csv := resp.Body.String()
	scanner := bufio.NewScanner(strings.NewReader(csv))

	//check the headers are correct
	ct := false
	for _, header := range resp.Header()[http.CanonicalHeaderKey("content-type")] {
		if header == "text/csv" {
			ct = true
			break
		}
	}
	IsTrue(t, ct, "content-type header incorrect or missing")
	IsTrue(t, len(resp.Header()[http.CanonicalHeaderKey("content-disposition")]) > 0, "content-disposition missing")

	//Read the header line
	scanner.Scan()
	endCol := len(strings.Split(scanner.Text(), ",")) - 1

	//convert the csv into an array of int arrays
	rs := [][]int{}
	i := 0
	for scanner.Scan() {
		strs := strings.Split(scanner.Text(), ",")
		vals := make([]int, endCol+1)
		for j, val := range strs {
			vals[j], _ = strconv.Atoi(val)
		}
		rs = append(rs, vals)
		i++
	}
	//Check the results are correct
	AreEqual(t, 6, len(rs), "Wrong number of iterations")
	AreEqual(t, 1, rs[0][endCol], "Wrong conversation count")
	AreEqual(t, 1, rs[1][endCol], "Wrong conversation count")
	AreEqual(t, 2, rs[2][endCol], "Wrong conversation count")
	AreEqual(t, 1, rs[3][endCol], "Wrong conversation count")
	AreEqual(t, 2, rs[4][endCol], "Wrong conversation count")
	AreEqual(t, 3, rs[5][endCol], "Wrong conversation count")
	AreEqual(t, 3, rs[0][3], "Wrong color count")
	AreEqual(t, 3, rs[1][3], "Wrong color count")
	AreEqual(t, 4, rs[2][3], "Wrong color count")
	AreEqual(t, 3, rs[3][3], "Wrong color count")
	AreEqual(t, 4, rs[4][3], "Wrong color count")
	AreEqual(t, 5, rs[5][3], "Wrong color count")
}
