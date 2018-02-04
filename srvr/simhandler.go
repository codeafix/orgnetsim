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
	GetStepsOrResults(c *mango.Context)
	GetResults(siminfo *SimInfo, c *mango.Context)
	RunOrGenerateNetwork(c *mango.Context)
	PostRun(siminfo *SimInfo, c *mango.Context)
	GenerateNetwork(siminfo *SimInfo, c *mango.Context)
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
	r.Get("/api/simulation/{sim_id}/{stepOrResults}", sh.GetStepsOrResults)
	r.Post("/api/simulation/{sim_id}/{runOrGenerate}", sh.RunOrGenerateNetwork)
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

//GetStepsOrResults gets the list of steps in this simulation
func (sh *SimHandlerState) GetStepsOrResults(c *mango.Context) {
	siminfo := NewSimInfo(c.RouteParams["sim_id"])

	switch c.RouteParams["stepOrResults"] {
	case "step":
		sh.GetList(siminfo, c, "step")
		return
	case "results":
		sh.GetResults(siminfo, c)
		return
	default:
		c.Error("Not Found", http.StatusNotFound)
	}
	return
}

//GetResults gets a concatenated set of results from all the steps in this simulation
func (sh *SimHandlerState) GetResults(siminfo *SimInfo, c *mango.Context) {

}

//RunOrGenerateNetwork gets the list of steps in this simulation
func (sh *SimHandlerState) RunOrGenerateNetwork(c *mango.Context) {
	siminfo := NewSimInfo(c.RouteParams["sim_id"])

	switch c.RouteParams["runOrGenerate"] {
	case "run":
		sh.PostRun(siminfo, c)
		return
	case "generate":
		sh.GenerateNetwork(siminfo, c)
		return
	default:
		c.Error("Not Found", http.StatusNotFound)
	}
	return
}

//PostRun adds a new step to the list of simulations
func (sh *SimHandlerState) PostRun(siminfo *SimInfo, c *mango.Context) {

}

//GenerateNetwork adds a new step to the list of simulations
func (sh *SimHandlerState) GenerateNetwork(siminfo *SimInfo, c *mango.Context) {

}

//DeleteStep removes a simulation from the list of simulations
func (sh *SimHandlerState) DeleteStep(c *mango.Context) {
	siminfo := NewSimInfo(c.RouteParams["sim_id"])
	step := NewSimStep(c.RouteParams["step_id"], c.RouteParams["sim_id"])
	sh.DeleteItem(step, siminfo, c, "step")
}
