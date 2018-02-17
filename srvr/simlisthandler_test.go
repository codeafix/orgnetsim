package srvr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/spaceweasel/mango"
)

func CreateSimListHandlerBrowser() (*mango.Browser, *TestFileUpdater, *TestFileUpdater, *TestFileManager) {
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
	dfu := &TestFileUpdater{}
	tfm.SetDefault(dfu)

	r := CreateRouter(tfm)
	br := mango.NewBrowser(r)

	return br, slfu, dfu, tfm
}

func TestUpdateSimListNotesSuccess(t *testing.T) {
	br, slfu, _, _ := CreateSimListHandlerBrowser()

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	data := `{"Notes":"some different notes"}`
	resp, err := br.PutS("/api/simulation/notes", data, hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusOK, resp.Code, "Not OK")

	AreEqual(t, 1, len(slfu.Obj.(*SimList).Items), "There should still be one item in the list")
	AreEqual(t, "/api/simulation/{someIdHere}", slfu.Obj.(*SimList).Items[0], "Sim list changed when updating notes")

	rsl := &SimList{}
	err = json.Unmarshal([]byte(resp.Body.String()), rsl)
	AssertSuccess(t, err)
	AreEqual(t, 1, len(rsl.Items), "There should be one item in the returned list")
	AreEqual(t, "/api/simulation/{someIdHere}", rsl.Items[0], "Returned Sim list changed when updating notes")
	AreEqual(t, "some different notes", rsl.Notes, "Returned notes not updated")
}

func TestGetSimList(t *testing.T) {
	br, _, _, _ := CreateSimListHandlerBrowser()

	hdrs := http.Header{}
	resp, err := br.Get("/api/simulation", hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusOK, resp.Code, "Not OK")

	rsl := NewSimList()
	json.Unmarshal([]byte(resp.Body.String()), rsl)
	AreEqual(t, 1, len(rsl.Items), "Wrong number of items in returned list")
	AreEqual(t, "/api/simulation/{someIdHere}", rsl.Items[0], "Wrong item in list")
	AreEqual(t, "Some notes", rsl.Notes, "Wrong notes")
}

func TestGetSimListReturnErrorWhenReadFails(t *testing.T) {
	br, slfu, _, _ := CreateSimListHandlerBrowser()
	slfu.ReadErr = fmt.Errorf("Access denied")

	hdrs := http.Header{}
	resp, err := br.Get("/api/simulation", hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusInternalServerError, resp.Code, "Not Error")
	AreEqual(t, "Access denied", resp.Body.String(), "Wrong notes")
}

func TestAddSimToSimListSuccess(t *testing.T) {
	br, slfu, dfu, _ := CreateSimListHandlerBrowser()

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	data := `{"Name":"mySim","Description":"A description of my sim"}`
	resp, err := br.PostS("/api/simulation", data, hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusCreated, resp.Code, "Not OK")

	AreEqual(t, 2, len(slfu.Obj.(*SimList).Items), "There should be two items in the list")
	AreEqual(t, fmt.Sprintf("/api/simulation/%s", dfu.Obj.(*SimInfo).ID), slfu.Obj.(*SimList).Items[1], "Wrong path added to sim list")

	rsim := &SimInfo{}
	err = json.Unmarshal([]byte(resp.Body.String()), rsim)
	AssertSuccess(t, err)
	AreEqual(t, dfu.Obj.(*SimInfo).ID, rsim.ID, "ID of returned SimInfo incorrect")
	AreEqual(t, "mySim", rsim.Name, "Name of returned SimInfo incorrect")
	AreEqual(t, "A description of my sim", rsim.Description, "Description of returned SimInfo incorrect")
}

func TestAddsSimToEmptySimListSuccess(t *testing.T) {
	br, slfu, dfu, _ := CreateSimListHandlerBrowser()
	slfu.Obj = NewSimList()

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	data := `{"Name":"mySim","Description":"A description of my sim"}`
	resp, err := br.PostS("/api/simulation", data, hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusCreated, resp.Code, "Not OK")

	AreEqual(t, 1, len(slfu.Obj.(*SimList).Items), "There should be one item in the list")
	AreEqual(t, fmt.Sprintf("/api/simulation/%s", dfu.Obj.(*SimInfo).ID), slfu.Obj.(*SimList).Items[0], "Wrong path added to sim list")

	rsim := &SimInfo{}
	err = json.Unmarshal([]byte(resp.Body.String()), rsim)
	AssertSuccess(t, err)
	AreEqual(t, dfu.Obj.(*SimInfo).ID, rsim.ID, "ID of returned SimInfo incorrect")
	AreEqual(t, "mySim", rsim.Name, "Name of returned SimInfo incorrect")
	AreEqual(t, "A description of my sim", rsim.Description, "Description of returned SimInfo incorrect")
}

func TestAddSimListFailsWithWrongHeaders(t *testing.T) {
	br, _, _, _ := CreateSimListHandlerBrowser()

	hdrs := http.Header{}
	data := `{"Name":"mySim","Description":"A description of my sim"}`
	resp, err := br.PostS("/api/simulation", data, hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusBadRequest, resp.Code, "Not Bad request")
}

func TestAddSimListReturnsErrorWithFileReadErr(t *testing.T) {
	br, slfu, _, _ := CreateSimListHandlerBrowser()
	slfu.ReadErr = fmt.Errorf("Access denied")

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	data := `{"Name":"mySim","Description":"A description of my sim"}`
	resp, err := br.PostS("/api/simulation", data, hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusInternalServerError, resp.Code, "No error reported")
	Contains(t, "Access denied", resp.Body.String(), "Wrong error reported")
}

func TestAddSimListReturnsErrorWithFileUpdateErr(t *testing.T) {
	br, slfu, _, _ := CreateSimListHandlerBrowser()
	slfu.UpdateErr = fmt.Errorf("Access denied")

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	data := `{"Name":"mySim","Description":"A description of my sim"}`
	resp, err := br.PostS("/api/simulation", data, hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusInternalServerError, resp.Code, "No error reported")
	Contains(t, "Access denied", resp.Body.String(), "Wrong error reported")
}

func TestAddSimListReturnsErrorWithFileCreateErr(t *testing.T) {
	br, _, dfu, _ := CreateSimListHandlerBrowser()
	dfu.CreateErr = fmt.Errorf("Access denied")

	hdrs := http.Header{}
	hdrs.Set("Content-Type", "application/json")
	data := `{"Name":"mySim","Description":"A description of my sim"}`
	resp, err := br.PostS("/api/simulation", data, hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusInternalServerError, resp.Code, "No error reported")
	Contains(t, "Access denied", resp.Body.String(), "Wrong error reported")
}

func SetupForDeleteTests(deleteItemIndex int) (*mango.Browser, *TestFileUpdater, *TestFileUpdater, []string) {
	br, slfu, _, tfm := CreateSimListHandlerBrowser()
	sl := NewSimList()
	ids := []string{
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
	}
	items := make([]string, len(ids))
	for i, id := range ids {
		items[i] = fmt.Sprintf("/api/simulation/%s", id)
	}
	sl.Items = items
	slfu.Obj = sl
	sim := &SimInfo{
		ID: ids[deleteItemIndex],
	}
	simfu := &TestFileUpdater{
		Obj:      sim,
		Filepath: sim.Filepath(),
	}
	tfm.Add(sim.Filepath(), simfu)
	return br, slfu, simfu, items
}

func TestDeleteFirstSimFromSimListSuccess(t *testing.T) {
	br, slfu, simfu, items := SetupForDeleteTests(0)

	hdrs := http.Header{}
	resp, err := br.Delete(items[0], hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusOK, resp.Code, "Not OK")

	AreEqual(t, 2, len(slfu.Obj.(*SimList).Items), "There should be two items in the list")
	AreEqual(t, items[1], slfu.Obj.(*SimList).Items[0], "Wrong path in position 0 of sim list")
	AreEqual(t, items[2], slfu.Obj.(*SimList).Items[1], "Wrong path in position 1 of sim list")
	IsTrue(t, simfu.DeleteCalled, "Delete was not called on the correct fileupdater")
}

func TestDeleteMiddleSimFromSimListSuccess(t *testing.T) {
	br, slfu, simfu, items := SetupForDeleteTests(1)

	hdrs := http.Header{}
	resp, err := br.Delete(items[1], hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusOK, resp.Code, "Not OK")

	AreEqual(t, 2, len(slfu.Obj.(*SimList).Items), "There should be two items in the list")
	AreEqual(t, items[0], slfu.Obj.(*SimList).Items[0], "Wrong path in position 0 of sim list")
	AreEqual(t, items[2], slfu.Obj.(*SimList).Items[1], "Wrong path in position 1 of sim list")
	IsTrue(t, simfu.DeleteCalled, "Delete was not called on the correct fileupdater")
}

func TestDeleteLastSimFromSimListSuccess(t *testing.T) {
	br, slfu, simfu, items := SetupForDeleteTests(2)

	hdrs := http.Header{}
	resp, err := br.Delete(items[2], hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusOK, resp.Code, "Not OK")

	AreEqual(t, 2, len(slfu.Obj.(*SimList).Items), "There should be two items in the list")
	AreEqual(t, items[0], slfu.Obj.(*SimList).Items[0], "Wrong path in position 0 of sim list")
	AreEqual(t, items[1], slfu.Obj.(*SimList).Items[1], "Wrong path in position 1 of sim list")
	IsTrue(t, simfu.DeleteCalled, "Delete was not called on the correct fileupdater")
}

func TestDeleteSimFromSimListReturnsNotFoundWhenListEmpty(t *testing.T) {
	br, slfu, _, _ := CreateSimListHandlerBrowser()
	sl := NewSimList()
	slfu.Obj = sl

	hdrs := http.Header{}
	resp, err := br.Delete(fmt.Sprintf("/api/simulation/%s", uuid.New().String()), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusNotFound, resp.Code, "Not Found not reported")
}

func TestDeleteSimFromSimListReturnsNotFoundWhenNotInList(t *testing.T) {
	br, _, _, _ := SetupForDeleteTests(0)

	hdrs := http.Header{}
	resp, err := br.Delete(fmt.Sprintf("/api/simulation/%s", uuid.New().String()), hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusNotFound, resp.Code, "Not Found not reported")
}

func TestDeleteReturnsErrorWhenSimListReadErr(t *testing.T) {
	br, slfu, simfu, items := SetupForDeleteTests(0)
	slfu.ReadErr = fmt.Errorf("Access denied")

	hdrs := http.Header{}
	resp, err := br.Delete(items[0], hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusInternalServerError, resp.Code, "Not error")
	Contains(t, "Access denied", resp.Body.String(), "Wrong error reported")
	IsFalse(t, simfu.DeleteCalled, "Delete was still called on the Sim fileupdater")
}

func TestDeleteReturnsErrorWhenSimListUpdateErr(t *testing.T) {
	br, slfu, simfu, items := SetupForDeleteTests(0)
	slfu.UpdateErr = fmt.Errorf("Access denied")

	hdrs := http.Header{}
	resp, err := br.Delete(items[0], hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusInternalServerError, resp.Code, "Not error")
	Contains(t, "Access denied", resp.Body.String(), "Wrong error reported")
	IsFalse(t, simfu.DeleteCalled, "Delete was still called on the Sim fileupdater")
}

func TestDeleteReturnsErrorWhenSimDeleteErr(t *testing.T) {
	br, _, simfu, items := SetupForDeleteTests(0)
	simfu.DeleteErr = fmt.Errorf("Access denied")

	hdrs := http.Header{}
	resp, err := br.Delete(items[0], hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusInternalServerError, resp.Code, "Not error")
	Contains(t, "Access denied", resp.Body.String(), "Wrong error reported")
	IsTrue(t, simfu.DeleteCalled, "Delete was still called on the Sim fileupdater")
}
