package srvr

import (
	"fmt"

	"github.com/codeafix/orgnetsim/sim"
	"github.com/google/uuid"
)

//SimInfo contains all relevant information about a simulation
type SimInfo struct {
	TimestampHolder
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Steps       []string           `json:"steps"`
	Options     sim.NetworkOptions `json:"options"`
}

//CreateSimInfo creates a new SimInfo object with a new ID
func CreateSimInfo() *SimInfo {
	return &SimInfo{
		ID: uuid.New().String(),
	}
}

//NewSimInfo returns a SimInfo object for the passed ID
func NewSimInfo(id string) *SimInfo {
	return &SimInfo{
		ID: id,
	}
}

//CopyValues copies the values from the passed SimStep object to this object.
//Returns an error if the values could not be copied
func (si *SimInfo) CopyValues(obj Persistable) error {
	siToCopy, ok := obj.(*SimInfo)
	if !ok {
		return fmt.Errorf("Failed to copy values")
	}
	si.Name = siToCopy.Name
	si.Description = siToCopy.Description
	si.Options = siToCopy.Options
	return nil
}

//Filepath returns the Filepath used by this item
func (si *SimInfo) Filepath() string {
	return fmt.Sprintf("sim_%s.json", si.ID)
}

//RelPath returns the relative API path for this item
func (si *SimInfo) RelPath() string {
	return fmt.Sprintf("/api/simulation/%s", si.ID)
}

//GetItems returns the items in the specified list
func (si *SimInfo) GetItems(listname string) []string {
	return si.Steps
}

//UpdateItems updates the items in the specified list
func (si *SimInfo) UpdateItems(listname string, items []string) {
	si.Steps = items
}
