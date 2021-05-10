package srvr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

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
	data := `{"network":null,"results":{"iterations":5,"colors":null,"conversations":null}}`
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
	data := `{"network":null,"results":{"iterations":5,"colors":null,"conversations":null}}`
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
	data := `{"network":null,"results":{"iterations":5,"colors":null,"conversations":null}}`
	resp, err := br.PutS(fmt.Sprintf("/api/simulation/%s/step/%s", ts.ParentID, ts.ID), data, hdrs)
	AssertSuccess(t, err)
	AreEqual(t, http.StatusInternalServerError, resp.Code, "No error reported")
	Contains(t, "Access denied", resp.Body.String(), "Wrong error reported")
}
