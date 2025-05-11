package srvr

import (
	"fmt"
	"net/http"

	"github.com/codeafix/orgnetsim/sim"
	"github.com/spaceweasel/mango"
)

// StepHandlerState holds state data for the StepHandler
type StepHandlerState struct {
	PersistableHandlerState
}

// StepHandler provides Read/Update methods for steps of a simulation
type StepHandler interface {
	mango.Registerer
	Get(c *mango.Context)
	Put(c *mango.Context)
	GetStepData(c *mango.Context)
	PutStepNetworkData(c *mango.Context)
}

// NewStepHandler returns a new instance of StepHandler
func NewStepHandler(fm FileManager) StepHandler {
	return &StepHandlerState{
		PersistableHandlerState: PersistableHandlerState{
			FileManager: fm,
		},
	}
}

// Register the routes for this routehandler
func (sh *StepHandlerState) Register(r *mango.Router) {
	r.Get("/api/simulation/{sim_id}/step/{step_id}", sh.Get)
	r.Put("/api/simulation/{sim_id}/step/{step_id}", sh.Put)
	// New combined route for specific step data (netdata or agentcolors)
	r.Get("/api/simulation/{sim_id}/step/{step_id}/{datatype}", sh.GetStepData)
	r.Put("/api/simulation/{sim_id}/step/{step_id}/{datatype}", sh.PutStepNetworkData) // Corrected typo here
}

// Get returns an existing step within a simulation
func (sh *StepHandlerState) Get(c *mango.Context) {
	step := NewSimStep(c.RouteParams["step_id"], c.RouteParams["sim_id"])
	sh.GetObject(step, c)
}

// Put updates an existing step within a simulation
func (sh *StepHandlerState) Put(c *mango.Context) {
	step := NewSimStep(c.RouteParams["step_id"], c.RouteParams["sim_id"])
	savedstep := NewSimStep(c.RouteParams["step_id"], c.RouteParams["sim_id"])
	sh.UpdateObjectWithContextBind(step, savedstep, c)
}

// GetStepData returns specific data (network or agents) for a simulation step.
func (sh *StepHandlerState) GetStepData(c *mango.Context) {
	if c.RouteParams == nil {
		c.Error("RouteParams is nil in GetStepData", http.StatusInternalServerError)
		return
	}
	simID := c.RouteParams["sim_id"]
	stepID := c.RouteParams["step_id"]
	dataType := c.RouteParams["datatype"]

	if simID == "" || stepID == "" {
		c.Error(fmt.Sprintf("Simulation ID ('%s') or StepID ('%s') from route is empty", simID, stepID), http.StatusBadRequest)
		return
	}

	step := NewSimStep(stepID, simID)
	objUpdater := sh.FileManager.Get(step.Filepath())
	err := objUpdater.Read(step)
	if err != nil {
		c.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	if step.Network == nil {
		c.Error(fmt.Sprintf("Network data not found for stepID '%s', simID '%s'", stepID, simID), http.StatusNotFound)
		return
	}

	switch dataType {
	case "network":
		c.RespondWith(step.Network).WithStatus(http.StatusOK)
	case "agents":
		agents := step.Network.Agents()
		c.RespondWith(agents).WithStatus(http.StatusOK)
	default:
		c.Error("Not Found", http.StatusNotFound)
	}
}

// PutStepNetworkData updates the network data for a specific simulation step.
func (sh *StepHandlerState) PutStepNetworkData(c *mango.Context) {
	if c.RouteParams == nil {
		c.Error("RouteParams is nil in PutStepNetworkData", http.StatusInternalServerError)
		return
	}
	simID := c.RouteParams["sim_id"]
	stepID := c.RouteParams["step_id"]

	if simID == "" || stepID == "" {
		c.Error(fmt.Sprintf("Missing IDs in route: Simulation ID ('%s') or StepID ('%s') is empty", simID, stepID), http.StatusBadRequest)
		return
	}

	dataType := c.RouteParams["datatype"]
	if dataType != "network" {
		c.Error(fmt.Sprintf("Update to ('%s') from route is not supported. Only direct updates to 'network' are available.", dataType), http.StatusBadRequest)
		return
	}

	step := NewSimStep(stepID, simID)
	objUpdater := sh.FileManager.Get(step.Filepath())
	err := objUpdater.Read(step)
	if err != nil {
		c.Error(fmt.Sprintf("Error reading step: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	var newNetwork sim.Network
	if err := c.Bind(&newNetwork); err != nil {
		c.Error(fmt.Sprintf("Error binding network data: %s", err.Error()), http.StatusBadRequest)
		return
	}

	step.Network = &newNetwork
	if err := objUpdater.Update(step); err != nil {
		c.Error(fmt.Sprintf("Error updating step with new network: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	c.RespondWith(step.Network).WithStatus(http.StatusOK)
}
