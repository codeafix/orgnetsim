package srvr

import "github.com/spaceweasel/mango"

//StepHandlerState holds state data for the StepHandler
type StepHandlerState struct {
	PersistableHandlerState
}

//StepHandler provides Read/Update methods for steps of a simulation
type StepHandler interface {
	mango.Registerer
	Get(c *mango.Context)
	Put(c *mango.Context)
}

//NewStepHandler returns a new instance of StepHandler
func NewStepHandler(fm FileManager) StepHandler {
	return &StepHandlerState{
		PersistableHandlerState: PersistableHandlerState{
			FileManager: fm,
		},
	}
}

//Register the routes for this routehandler
func (sh *StepHandlerState) Register(r *mango.Router) {
	r.Get("/api/simulation/{sim_id}/step/{step_id}", sh.Get)
	r.Put("/api/simulation/{sim_id}/step/{step_id}", sh.Put)
}

//Get returns an existing step within a simulation
func (sh *StepHandlerState) Get(c *mango.Context) {
	step := NewSimStep(c.RouteParams["step_id"], c.RouteParams["sim_id"])
	sh.GetObject(step, c)
}

//Put updates an existing step within a simulation
func (sh *StepHandlerState) Put(c *mango.Context) {
	step := NewSimStep(c.RouteParams["step_id"], c.RouteParams["sim_id"])
	savedstep := NewSimStep(c.RouteParams["step_id"], c.RouteParams["sim_id"])
	sh.UpdateObjectWithContextBind(step, savedstep, c)
}
