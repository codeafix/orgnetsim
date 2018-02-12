package srvr

import (
	"net/http"

	"github.com/codeafix/orgnetsim/sim"
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
	r.Get("/api/simulation/{sim_id}/{stepResultsRunGenerate}", sh.GetStepsOrResults)
	r.Post("/api/simulation/{sim_id}/{stepResultsRunGenerate}", sh.RunOrGenerateNetwork)
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

	switch c.RouteParams["stepResultsRunGenerate"] {
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
	objUpdater := sh.ListHandlerState.FileManager.Get(siminfo.Filepath())
	err := objUpdater.Read(siminfo)
	if err != nil {
		c.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	switch c.RouteParams["stepResultsRunGenerate"] {
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

//GenerateNetwork generates a hierarchical network to be simulated.
func (sh *SimHandlerState) GenerateNetwork(siminfo *SimInfo, c *mango.Context) {
	hs := sim.HierarchySpec{}
	err := c.Bind(&hs)
	if err != nil {
		c.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	if len(siminfo.Steps) > 0 {
		c.Error("Simulation must have no steps when generating a new network", http.StatusBadRequest)
		return
	}
	rm, no, err := sim.GenerateHierarchy(hs)
	if err != nil {
		c.Error(err.Error(), http.StatusBadRequest)
		return
	}
	step := CreateSimStep(siminfo.ID)
	step.Network = rm
	step.Results = sim.Results{
		Iterations:    0,
		Colors:        make([][]int, 1, 1),
		Conversations: make([]int, 1, 1),
	}
	agents := rm.Agents()
	colorCounts := make([]int, rm.MaxColors(), rm.MaxColors())
	for _, a := range agents {
		colorCounts[a.GetColor()]++
	}
	step.Results.Colors[0] = colorCounts

	itemUpdater := sh.ListHandlerState.FileManager.Get(step.Filepath())
	err = itemUpdater.Create(step)
	if err != nil {
		c.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	listUpdater := sh.ListHandlerState.FileManager.Get(siminfo.Filepath())

	//Retry if there is a failure
	for i := 0; i < 2; i++ {
		err = listUpdater.Read(siminfo)
		if err != nil {
			continue
		}
		siminfo.Options = *no
		siminfo.Steps = append(siminfo.Steps, step.RelPath())
		err = listUpdater.Update(siminfo)
		if err == nil {
			break
		}
	}
	if err != nil {
		c.RespondWith(err.Error()).WithStatus(http.StatusInternalServerError)
	} else {
		c.RespondWith(step).WithStatus(http.StatusCreated)
	}
}

//DeleteStep removes a simulation from the list of simulations
func (sh *SimHandlerState) DeleteStep(c *mango.Context) {
	siminfo := NewSimInfo(c.RouteParams["sim_id"])
	step := NewSimStep(c.RouteParams["step_id"], c.RouteParams["sim_id"])
	sh.DeleteItem(step, siminfo, c, "step")
}
