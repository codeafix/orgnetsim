package srvr

import (
	"fmt"

	"github.com/codeafix/orgnetsim/sim"
)

//SimStep holds the results of each simulation step
type SimStep struct {
	TimestampHolder
	Network  sim.Network `json:"network"`
	Results  sim.Results `json:"results"`
	ID       string      `json:"id"`
	ParentID string      `json:"parent"`
}

//NewSimStep returns a SimStep object for the passed ID that will be persisted in directory root
func NewSimStep(id string, parentID string) *SimStep {
	return &SimStep{
		ID:       id,
		ParentID: parentID,
	}
}

//CopyValues copies the values from the passed SimStep object to this object.
//Returns an error if the values could not be copied
func (ss *SimStep) CopyValues(obj Persistable) error {
	ssToCopy, ok := obj.(*SimStep)
	if !ok {
		return fmt.Errorf("Failed to copy values")
	}
	ss.Network = ssToCopy.Network
	ss.Results = ssToCopy.Results
	return nil
}

//Filepath returns the Filepath used by this item
func (ss *SimStep) Filepath() string {
	return fmt.Sprintf("sim_%s.json", ss.ID)
}

//RelPath returns the relative API path for this item
func (ss *SimStep) RelPath() string {
	return fmt.Sprintf("/api/simulation/%s/step/%s", ss.ParentID, ss.ID)
}
