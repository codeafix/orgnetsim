package srvr

import (
	"fmt"
	"strings"

	"github.com/codeafix/orgnetsim/sim"
	"github.com/google/uuid"
)

//SimStep holds the results of each simulation step
type SimStep struct {
	TimestampHolder
	Network  sim.RelationshipMgr `json:"network"`
	Results  sim.Results         `json:"results"`
	ID       string              `json:"id"`
	ParentID string              `json:"parent"`
}

//CreateSimStep creates a new SimStep object with a new ID
func CreateSimStep(parentID string) *SimStep {
	rm, _ := sim.NewNetwork("")
	return &SimStep{
		ID:       uuid.New().String(),
		ParentID: parentID,
		Network:  rm,
		Results:  sim.Results{},
	}
}

//NewSimStep returns a SimStep object for the passed ID that will be persisted in directory root
func NewSimStep(id string, parentID string) *SimStep {
	rm, _ := sim.NewNetwork("")
	return &SimStep{
		ID:       id,
		ParentID: parentID,
		Network:  rm,
		Results:  sim.Results{},
	}
}

//NewSimStepFromRelPath returns a SimStep object extracting IDs from the relative path in the
//passed string
func NewSimStepFromRelPath(relPath string) *SimStep {
	elems := strings.Split(relPath, "/")
	rm, _ := sim.NewNetwork("")
	return &SimStep{
		ID:       elems[len(elems)-1],
		ParentID: elems[len(elems)-3],
		Network:  rm,
		Results:  sim.Results{},
	}
}

//CopyValues copies the values from the passed SimStep object to this object.
//Returns an error if the values could not be copied
func (ss *SimStep) CopyValues(obj Persistable) error {
	ssToCopy, ok := obj.(*SimStep)
	if !ok {
		return fmt.Errorf("failed to copy values")
	}
	ss.Network = ssToCopy.Network
	ss.Results = ssToCopy.Results
	return nil
}

//Filepath returns the Filepath used by this item
func (ss *SimStep) Filepath() string {
	return fmt.Sprintf("step_%s.json", ss.ID)
}

//RelPath returns the relative API path for this item
func (ss *SimStep) RelPath() string {
	return fmt.Sprintf("/api/simulation/%s/step/%s", ss.ParentID, ss.ID)
}
