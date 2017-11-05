package srvr

import (
	"net/http"

	"github.com/spaceweasel/mango"
)

//SimHandlerState holds state data for the SimHandler
type SimHandlerState struct {
	ListHandlerState
	PersistableHandlerState
}

//SimHandler provides Read/Update methods for simulations on the simulation list
type SimHandler interface {
	mango.Registerer
	Get(c *mango.Context)
	Put(c *mango.Context)
	GetSteps(c *mango.Context)
	PostRun(c *mango.Context)
	GenerateNetwork(c *mango.Context)
	GetResults(c *mango.Context)
	DeleteStep(c *mango.Context)
}

//NewSimHandler returns a new instance of SimHandler
func NewSimHandler(fm FileManager) SimHandler {
	return &SimHandlerState{
		ListHandlerState{
			FileManager: fm,
		},
		PersistableHandlerState{
			FileManager: fm,
		},
	}
}

//Register the routes for this routehandler
func (sh *SimHandlerState) Register(r *mango.Router) {
	r.Get("/api/simulation/{sim_id}", sh.Get)
	r.Put("/api/simulation/{sim_id}", sh.Put)
	r.Get("/api/simulation/{sim_id}/{step}", sh.GetSteps)
	r.Get("/api/simulation/{sim_id}/results", sh.GetResults)
	r.Post("/api/simulation/{sim_id}/run", sh.PostRun)
	r.Post("/api/simulation/{sim_id}/generate", sh.GenerateNetwork)
	r.Delete("/api/simulation/{sim_id}/step/{step_id}", sh.DeleteStep)
}

//Get returns an existing simulation
func (sh *SimHandlerState) Get(c *mango.Context) {
	siminfo := NewSimInfo(c.RouteParams["sim_id"])
	sh.GetObject(siminfo, c)
}

//Put updates an existing simulation
func (sh *SimHandlerState) Put(c *mango.Context) {
	siminfo := NewSimInfo(c.RouteParams["sim_id"])
	savedsiminfo := NewSimInfo(c.RouteParams["sim_id"])
	sh.UpdateObject(siminfo, savedsiminfo, c)
}

//GetSteps gets the list of steps in this simulation
func (sh *SimHandlerState) GetSteps(c *mango.Context) {
	siminfo := NewSimInfo(c.RouteParams["sim_id"])
	if "step" != c.RouteParams["step"] {
		c.Error("Not Found", http.StatusNotFound)
		return
	}
	sh.GetList(siminfo, c, "step")
}

//PostRun adds a new step to the list of simulations
func (sh *SimHandlerState) PostRun(c *mango.Context) {

}

//GenerateNetwork adds a new step to the list of simulations
func (sh *SimHandlerState) GenerateNetwork(c *mango.Context) {

}

//GetResults gets a concatenated set of results from all the steps in this simulation
func (sh *SimHandlerState) GetResults(c *mango.Context) {

}

//DeleteStep removes a simulation from the list of simulations
func (sh *SimHandlerState) DeleteStep(c *mango.Context) {
	siminfo := NewSimInfo(c.RouteParams["sim_id"])
	step := NewSimStep(c.RouteParams["step_id"], c.RouteParams["sim_id"])
	sh.DeleteItem(step, siminfo, c, "step")
}
