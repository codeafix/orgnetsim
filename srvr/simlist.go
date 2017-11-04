package srvr

import "fmt"

//SimList is the list of simulations in the root directory
type SimList struct {
	TimestampHolder
	Items []string `json:"simulations"`
	Notes string   `json:"notes"`
}

//NewSimList returns a SimList object that will be persisted
func NewSimList() *SimList {
	return &SimList{}
}

//CopyValues copies the values from the passed SimStep object to this object.
//Returns an error if the values could not be copied
func (sl *SimList) CopyValues(obj Persistable) error {
	slToCopy, ok := obj.(*SimList)
	if !ok {
		return fmt.Errorf("Failed to copy values")
	}
	sl.Notes = slToCopy.Notes
	return nil
}

//GetItems returns the items in the specified list
func (sl *SimList) GetItems(listname string) []string {
	return sl.Items
}

//UpdateItems updates the items in the specified list
func (sl *SimList) UpdateItems(listname string, items []string) {
	sl.Items = items
}

//Filepath returns the Filepath used by this item
func (sl *SimList) Filepath() string {
	return fmt.Sprintf("sims.json")
}
