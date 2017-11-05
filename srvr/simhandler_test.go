package srvr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/spaceweasel/mango"
)

func CreateSimHandlerBrowser(deleteItemIndex int) (*mango.Browser, *TestFileUpdater, *TestFileUpdater, []string, string) {
	r := mango.NewRouter()

	simid := uuid.New().String()
	sim := NewSimInfo(simid)
	sim.Name = "mySavedSim"
	sim.Description = "A description of mySavedSim"
	ids := []string{
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
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
	}
	ssfu := &TestFileUpdater{
		Obj:      ss,
		Filepath: ss.Filepath(),
	}
	tfm.Add(ss.Filepath(), ssfu)

	r.RegisterModules([]mango.Registerer{
		NewSimHandler(tfm),
	})
	br := mango.NewBrowser(r)

	return br, simfu, ssfu, steps, simid
}

func TestGetSimSuccess(t *testing.T) {
	br, simfu, _, _, simid := CreateSimHandlerBrowser(0)

	hdrs := http.Header{}
	resp, err := br.Get(fmt.Sprintf("/api/simulation/%s", simid), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusOK, resp.Code, "Not OK")

	rsim := &SimInfo{}
	err = json.Unmarshal([]byte(resp.Body.String()), rsim)
	AssertSuccess(t, err)
	AreEqual(t, simfu.Obj.(*SimInfo).Name, rsim.Name, "Wrong name in returned SimInfo")
	AreEqual(t, simfu.Obj.(*SimInfo).Description, rsim.Description, "Wrong description in returned SimInfo")
}

func TestGetSimInvalidCommand(t *testing.T) {
	br, _, _, _, simid := CreateSimHandlerBrowser(0)

	hdrs := http.Header{}
	resp, err := br.Get(fmt.Sprintf("/api/simulation/%s/somename", simid), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusNotFound, resp.Code, "Did not return Not Found")
}

func TestGetSimStepsSuccess(t *testing.T) {
	br, _, _, steps, simid := CreateSimHandlerBrowser(0)

	hdrs := http.Header{}
	resp, err := br.Get(fmt.Sprintf("/api/simulation/%s/step", simid), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusOK, resp.Code, "Not OK")

	rsim := &SimInfo{}
	err = json.Unmarshal([]byte(resp.Body.String()), rsim)
	AssertSuccess(t, err)
	AreEqual(t, 3, len(rsim.Steps), "Wrong number of Steps in returned SimInfo")
	AreEqual(t, steps[0], rsim.Steps[0], "Wrong Step 0 in returned SimInfo")
	AreEqual(t, steps[1], rsim.Steps[1], "Wrong Step 1 in returned SimInfo")
	AreEqual(t, steps[2], rsim.Steps[2], "Wrong Step 2 in returned SimInfo")
}

func TestUpdateSimSuccess(t *testing.T) {
	br, _, _, steps, simid := CreateSimHandlerBrowser(0)

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	data := `{"Name":"myUpdatedSim","Description":"A description of mySavedSim"}`
	resp, err := br.PutS(fmt.Sprintf("/api/simulation/%s", simid), data, hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusOK, resp.Code, "Not OK")

	rsim := &SimInfo{}
	err = json.Unmarshal([]byte(resp.Body.String()), rsim)
	AssertSuccess(t, err)
	AreEqual(t, "myUpdatedSim", rsim.Name, "Wrong name in returned SimInfo")
	AreEqual(t, "A description of mySavedSim", rsim.Description, "Wrong description in returned SimInfo")
	AreEqual(t, 3, len(rsim.Steps), "Wrong number of Steps in returned SimInfo")
	AreEqual(t, steps[0], rsim.Steps[0], "Wrong Step 0 in returned SimInfo")
	AreEqual(t, steps[1], rsim.Steps[1], "Wrong Step 1 in returned SimInfo")
	AreEqual(t, steps[2], rsim.Steps[2], "Wrong Step 2 in returned SimInfo")
}

func TestDeleteSimSuccess(t *testing.T) {
	br, simfu, ssfu, steps, _ := CreateSimHandlerBrowser(0)

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
