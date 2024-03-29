package srvr

import (
	"strings"

	"github.com/spaceweasel/mango"
)

//SimListHandlerState holds state data for the SimListHandler
type SimListHandlerState struct {
	ListHandlerState
	PersistableHandlerState
}

//SimListHandler provides Create/Delete methods for simulations on the simulation list
type SimListHandler interface {
	mango.Registerer
	GetSimulations(c *mango.Context)
	AddSimulation(c *mango.Context)
	UpdateNotes(c *mango.Context)
	DeleteSimulation(c *mango.Context)
}

//SimListExt is the struct that embeds SimInfo objects into the SimList
type SimListExt struct {
	TimestampHolder
	Items []*SimInfo `json:"simulations"`
	Notes string     `json:"notes"`
}

//NewSimListHandler returns a new instance of SimListHandler
func NewSimListHandler(fm FileManager) SimListHandler {
	sh := &SimListHandlerState{
		ListHandlerState{
			FileManager: fm,
		},
		PersistableHandlerState{
			FileManager: fm,
		},
	}
	sh.ListHandlerState.EncodeFunc = sh.EncodeFunc
	return sh
}

func (sh *SimListHandlerState) EncodeFunc(listHolder ListHolder, listname string) (interface{}, error) {
	paths := listHolder.GetItems(listname)
	items := []*SimInfo{}
	for _, path := range paths {
		elems := strings.Split(path, "/")
		newItem := NewSimInfo(elems[len(elems)-1])
		itemUpdater := sh.ListHandlerState.FileManager.Get(newItem.Filepath())
		err := itemUpdater.Read(newItem)
		if err != nil {
			return nil, err
		}
		items = append(items, newItem)
	}
	return SimListExt{
		TimestampHolder: TimestampHolder{
			Stamp: listHolder.Timestamp(),
		},
		Notes: listHolder.(*SimList).Notes,
		Items: items,
	}, nil
}

//Register the routes for this routehandler
func (sh *SimListHandlerState) Register(r *mango.Router) {
	r.Get("/api/simulation", sh.GetSimulations)
	r.Post("/api/simulation", sh.AddSimulation)
	r.Put("/api/simulation/notes", sh.UpdateNotes)
	r.Delete("/api/simulation/{sim_id}", sh.DeleteSimulation)
}

//GetSimulations gets the list of simulations
func (sh *SimListHandlerState) GetSimulations(c *mango.Context) {
	simlist := NewSimList()
	sh.GetList(simlist, c, "sim")
}

//AddSimulation adds a new simulation to the list of simulations
func (sh *SimListHandlerState) AddSimulation(c *mango.Context) {
	simlist := NewSimList()
	sim := CreateSimInfo()
	sh.AddItemWithContextBind(sim, simlist, c, "sim")
}

//UpdateNotes updates the notes on the sim list
func (sh *SimListHandlerState) UpdateNotes(c *mango.Context) {
	simlist := NewSimList()
	savedsimlist := NewSimList()
	sh.UpdateObjectWithContextBind(simlist, savedsimlist, c)
}

//DeleteSimulation removes a simulation from the list of simulations
func (sh *SimListHandlerState) DeleteSimulation(c *mango.Context) {
	simlist := NewSimList()
	sim := NewSimInfo(c.RouteParams["sim_id"])
	sh.DeleteItem(sim, simlist, c, "sim")
}
