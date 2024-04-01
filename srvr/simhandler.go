package srvr

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/codeafix/orgnetsim/sim"
	"github.com/spaceweasel/mango"
)

// SimHandlerState holds state data for the SimHandler
type SimHandlerState struct {
	ListHandlerState
	PersistableHandlerState
}

// SimHandler provides Read/Update methods for simulations on the simulation list
type SimHandler interface {
	mango.Registerer
	Get(c *mango.Context)
	Put(c *mango.Context)
	GetStepsOrResults(c *mango.Context)
	GetResults(c *mango.Context)
	RunGenerateParseNetwork(c *mango.Context)
	PostRun(siminfo *SimInfo, c *mango.Context)
	GenerateNetwork(siminfo *SimInfo, c *mango.Context)
	DeleteStep(c *mango.Context)
}

// ParseBody is the payload struct for uploading a network in a text file to be parsed
// together with the options that specify how to parse the text file
type ParseBody struct {
	sim.ParseOptions
	Payload []byte
}

// NewSimHandler returns a new instance of SimHandler
func NewSimHandler(fm FileManager) SimHandler {
	sh := &SimHandlerState{
		ListHandlerState{
			FileManager: fm,
		},
		PersistableHandlerState{
			FileManager: fm,
		},
	}
	sh.ListHandlerState.EncodeFunc = sh.EncodeStepList
	return sh
}

func (sh *SimHandlerState) EncodeStepList(listHolder ListHolder, listname string) (interface{}, error) {
	paths := listHolder.GetItems(listname)
	items := []*SimStep{}
	for _, path := range paths {
		elems := strings.Split(path, "/")
		newItem := NewSimStep(elems[len(elems)-1], elems[len(elems)-3])
		itemUpdater := sh.ListHandlerState.FileManager.Get(newItem.Filepath())
		err := itemUpdater.Read(newItem)
		if err != nil {
			return nil, err
		}
		items = append(items, newItem)
	}
	return items, nil
}

// Register the routes for this routehandler
func (sh *SimHandlerState) Register(r *mango.Router) {
	r.Get("/api/simulation/{sim_id}", sh.Get)
	r.Put("/api/simulation/{sim_id}", sh.Put)
	r.Get("/api/simulation/{sim_id}/{stepResultsRunGenerate}", sh.GetStepsOrResults)
	r.Post("/api/simulation/{sim_id}/{stepResultsRunGenerate}", sh.RunGenerateParseNetwork)
	r.Put("/api/simulation/{sim_id}/links", sh.AddLinks)
	r.Delete("/api/simulation/{sim_id}/step/{step_id}", sh.DeleteStep)
}

// Get returns an existing simulation
func (sh *SimHandlerState) Get(c *mango.Context) {
	siminfo := NewSimInfo(c.RouteParams["sim_id"])
	sh.GetObject(siminfo, c)
}

// Put updates an existing simulation
func (sh *SimHandlerState) Put(c *mango.Context) {
	siminfo := NewSimInfo(c.RouteParams["sim_id"])
	savedsiminfo := NewSimInfo(c.RouteParams["sim_id"])
	sh.UpdateObjectWithContextBind(siminfo, savedsiminfo, c)
}

// GetStepsOrResults gets the list of steps in this simulation
func (sh *SimHandlerState) GetStepsOrResults(c *mango.Context) {
	switch c.RouteParams["stepResultsRunGenerate"] {
	case "step":
		siminfo := NewSimInfo(c.RouteParams["sim_id"])
		sh.GetList(siminfo, c, "step")
		return
	case "results":
		for _, header := range c.Request.Header[http.CanonicalHeaderKey("content-type")] {
			if header == "text/csv" {
				sh.GetResultsCsv(c)
				return
			}
		}
		sh.GetResults(c)
		return
	default:
		c.Error("Not Found", http.StatusNotFound)
	}
}

// GetResultsCsv returns a concatenated set of results from all the steps in this simulation in text/csv format
func (sh *SimHandlerState) GetResultsCsv(c *mango.Context) {
	results, name, err := sh.collectAllResults(c)
	if err != nil {
		c.RespondWith(err.Error()).WithStatus(http.StatusInternalServerError)
	}
	if results.Iterations == 0 {
		c.RespondWith("this simulation has no iterations").WithStatus(http.StatusBadRequest)
	}
	var buffer bytes.Buffer

	maxColors := len(results.Colors[0])
	for c := 0; c < maxColors; c++ {
		buffer.WriteString(fmt.Sprintf("%s,", sim.Color(c).String()))
	}
	buffer.WriteString("Conversations\n")

	for i := 0; i < results.Iterations; i++ {
		for j := 0; j < maxColors; j++ {
			buffer.WriteString(fmt.Sprintf("%d,", results.Colors[i][j]))
		}
		buffer.WriteString(fmt.Sprintf("%d\n", results.Conversations[i]))
	}

	r := c.RespondWith(buffer.String())
	r.WithContentType("text/csv")
	r.WithHeader(http.CanonicalHeaderKey("Content-Disposition"), fmt.Sprintf("attachment; filename=\"%s.csv\"; filename*=\"%s.csv\"", name, name))
	r.WithStatus(http.StatusOK)

}

// GetResults returns a concatenated set of results from all the steps in this simulation in JSON format
func (sh *SimHandlerState) GetResults(c *mango.Context) {
	results, _, err := sh.collectAllResults(c)
	if err != nil {
		c.RespondWith(err.Error()).WithStatus(http.StatusInternalServerError)
	}
	c.RespondWith(results).WithStatus(http.StatusOK)
}

// collectAllResults gets a concatenated set of results from all the steps in this simulation
func (sh *SimHandlerState) collectAllResults(c *mango.Context) (sim.Results, string, error) {
	siminfo := sh.readSiminfo(c)
	results := sim.Results{
		Iterations:    0,
		Colors:        [][]int{},
		Conversations: []int{},
	}
	if siminfo == nil {
		return results, "", errors.New("unable to read simulation")
	}

	for _, spath := range siminfo.Steps {
		step := NewSimStepFromRelPath(spath)
		objUpdater := sh.ListHandlerState.FileManager.Get(step.Filepath())
		err := objUpdater.Read(step)
		if err != nil {
			return results, "", err
		}
		results.Iterations += step.Results.Iterations
		results.Colors = append(results.Colors, step.Results.Colors...)
		results.Conversations = append(results.Conversations, step.Results.Conversations...)
	}
	return results, siminfo.Name, nil
}

// RunGenerateParseNetwork handles three possible routes:
// /simulation/{id}/run Runs the simulation for the specified number of steps and iterations.
// /simulation/{id}/generate Generates a network to simulate, this will throw if the
// simulation already has steps.
// /simulation/{id}/parse Parses a network specified in a text file and sets it as the
// network to simulate. This will throw if the simulation already has steps.
// /simulation/{id}/links Parses a text file containing links and adds them to the network.
func (sh *SimHandlerState) RunGenerateParseNetwork(c *mango.Context) {
	siminfo := sh.readSiminfo(c)
	if siminfo == nil {
		return
	}
	switch c.RouteParams["stepResultsRunGenerate"] {
	case "run":
		sh.PostRun(siminfo, c)
		return
	case "generate":
		sh.GenerateNetwork(siminfo, c)
		return
	case "parse":
		sh.ParseNetwork(siminfo, c)
		return
	default:
		c.Error("Not Found", http.StatusNotFound)
	}
}

// RunSpec specifies the number of simulation steps to run, and the number of
// iterations that should be performed within each step
type RunSpec struct {
	Steps      int `json:"steps"`
	Iterations int `json:"iterations"`
}

// PostRun adds a new step to the list of simulations
func (sh *SimHandlerState) PostRun(siminfo *SimInfo, c *mango.Context) {
	rs := RunSpec{}
	err := c.Bind(&rs)
	if err != nil {
		c.Error(err.Error()+": Error reading RunSpec", http.StatusBadRequest)
		return
	}
	if len(siminfo.Steps) == 0 {
		c.Error("The simulation cannot be run without an initial step containing a network", http.StatusBadRequest)
		return
	}
	if rs.Steps <= 0 || rs.Iterations <= 0 {
		c.Error("Steps and Iterations cannot be zero", http.StatusBadRequest)
		return
	}
	ls := NewSimStepFromRelPath(siminfo.Steps[len(siminfo.Steps)-1])
	objUpdater := sh.ListHandlerState.FileManager.Get(ls.Filepath())
	err = objUpdater.Read(ls)
	if err != nil {
		c.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	err = ls.Network.PopulateMaps()
	if err != nil {
		c.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	r := sim.NewRunner(ls.Network, rs.Iterations)
	var ns *SimStep
	for i := 0; i < rs.Steps; i++ {
		ns = CreateSimStep(siminfo.ID)
		ns.Results = r.Run()
		ns.Network = r.GetRelationshipMgr()
		err = sh.AddItem(ns, siminfo, c, "step")
		if err != nil {
			c.Error(err.Error(), http.StatusInternalServerError)
			return
		}
	}
	c.RespondWith(ns).WithStatus(http.StatusCreated)
}

// GenerateNetwork generates a hierarchical network to be simulated.
func (sh *SimHandlerState) GenerateNetwork(siminfo *SimInfo, c *mango.Context) {
	hs := sim.HierarchySpec{}
	err := c.Bind(&hs)
	if err != nil {
		c.Error(err.Error()+": Error reading HierarchySpec", http.StatusBadRequest)
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
	siminfo.Options = *no
	savedsiminfo := NewSimInfo(c.RouteParams["sim_id"])
	err = sh.UpdateObject(siminfo, savedsiminfo, c)
	if err != nil {
		c.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	sh.createFirstSimStep(savedsiminfo, rm, c)
}

// createFirstSimStep creates a first simulation step in the passed simulation
// assigns the passed network to it and saves it all
func (sh *SimHandlerState) createFirstSimStep(siminfo *SimInfo, rm sim.RelationshipMgr, c *mango.Context) {
	step := CreateSimStep(siminfo.ID)
	step.Network = rm
	step.Results = sim.Results{
		Iterations:    0,
		Colors:        make([][]int, 1),
		Conversations: make([]int, 1),
	}
	agents := rm.Agents()
	colorCounts := make([]int, rm.MaxColors())
	for _, a := range agents {
		colorCounts[a.GetColor()]++
	}
	step.Results.Colors[0] = colorCounts
	err := sh.AddItem(step, siminfo, c, "step")
	if err != nil {
		c.Error(err.Error(), http.StatusInternalServerError)
	} else {
		c.RespondWith(step).WithStatus(http.StatusCreated)
	}
}

// ParseNetwork parses a network from a text file uploaded in the body of the post.
// The network is modified according to the options already stored in the simulation
// and then set as the starting point for the simulation. This will throw if a
// the simulation already has steps.
func (sh *SimHandlerState) ParseNetwork(siminfo *SimInfo, c *mango.Context) {
	of := ParseBody{}
	err := c.Bind(&of)
	if err != nil {
		c.Error(err.Error()+": Error reading ParseOptions", http.StatusBadRequest)
		return
	}
	if len(siminfo.Steps) > 0 {
		c.Error("Simulation must have no steps when parsing a new network", http.StatusBadRequest)
		return
	}
	r := []string{}
	br := bytes.NewReader(of.Payload)
	s := bufio.NewScanner(br)
	for s.Scan() {
		r = append(r, s.Text())
	}
	err = s.Err()
	if err != nil {
		c.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	if len(r) == 0 {
		c.Error("No links data in ParseOptions", http.StatusBadRequest)
		return
	}

	rm, err := of.ParseDelim(r)
	if err != nil {
		c.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	crm, err := siminfo.Options.CloneModify(rm)
	if err != nil {
		c.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	sh.createFirstSimStep(siminfo, crm, c)
}

// AddLinks parses links from a text file uploaded in the body of the post and adds
// them to the network in the latest step of the simulation. Unlike ParseNetwork, this
// method does not throw if the simulation already has steps, and it also does not
// modify the network according to the options in the simulation.
func (sh *SimHandlerState) AddLinks(c *mango.Context) {
	siminfo := sh.readSiminfo(c)
	if siminfo == nil {
		return
	}

	of := ParseBody{}
	err := c.Bind(&of)
	if err != nil {
		c.Error(err.Error()+": Error reading ParseOptions", http.StatusBadRequest)
		return
	}

	if len(siminfo.Steps) == 0 {
		c.Error("The network cannot have links added without an initial step containing a network", http.StatusBadRequest)
		return
	}

	r := []string{}
	br := bytes.NewReader(of.Payload)
	s := bufio.NewScanner(br)
	for s.Scan() {
		r = append(r, s.Text())
	}
	err = s.Err()
	if err != nil {
		c.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	if len(r) == 0 {
		c.Error("No links data in ParseOptions", http.StatusBadRequest)
		return
	}

	ls := NewSimStepFromRelPath(siminfo.Steps[len(siminfo.Steps)-1])
	objUpdater := sh.ListHandlerState.FileManager.Get(ls.Filepath())
	err = objUpdater.Read(ls)
	if err != nil {
		c.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	err = ls.Network.PopulateMaps()
	if err != nil {
		c.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	crm, err := of.ParseEdges(r, ls.Network)

	if err != nil {
		c.Error(err.Error(), http.StatusBadRequest)
		return
	}

	ls.Network = crm

	err = objUpdater.Update(ls)

	if err != nil {
		c.Error(err.Error(), http.StatusInternalServerError)
	} else {
		c.RespondWith(ls).WithStatus(http.StatusOK)
	}
}

// readSiminfo reads the SimInfo object specifed by the Id in the route
// parameter into a SimInfo struct and returns the pointer to it
func (sh *SimHandlerState) readSiminfo(c *mango.Context) *SimInfo {
	siminfo := NewSimInfo(c.RouteParams["sim_id"])
	objUpdater := sh.ListHandlerState.FileManager.Get(siminfo.Filepath())
	err := objUpdater.Read(siminfo)
	if err != nil {
		c.Error(err.Error(), http.StatusInternalServerError)
		return nil
	}
	return siminfo
}

// DeleteStep removes a simulation from the list of simulations
func (sh *SimHandlerState) DeleteStep(c *mango.Context) {
	siminfo := NewSimInfo(c.RouteParams["sim_id"])
	step := NewSimStep(c.RouteParams["step_id"], c.RouteParams["sim_id"])
	sh.DeleteItem(step, siminfo, c, "step")
}
