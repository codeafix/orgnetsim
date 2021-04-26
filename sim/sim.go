package sim

import (
	"math/rand"
	"time"
)

//Results contains the results from a Sim run over a number of iterations
type Results struct {
	Iterations    int     `json:"iterations"`
	Colors        [][]int `json:"colors"`
	Conversations []int   `json:"conversations"`
}

//RunnerInfo specifies the number of iterations and steps to run and records the results
type RunnerInfo struct {
	RelationshipMgr RelationshipMgr `json:"network"`
	Iterations      int             `json:"iterations"`
}

//Runner is used to run a simulation for a specified number of steps on its network
type Runner interface {
	Run() Results
	GetRelationshipMgr() RelationshipMgr
}

//NewRunner returns an instance of a sim Runner
func NewRunner(n RelationshipMgr, iterations int) Runner {
	return &RunnerInfo{
		RelationshipMgr: n,
		Iterations:      iterations,
	}
}

//GetRelationshipMgr returns the internal network state
func (ri *RunnerInfo) GetRelationshipMgr() RelationshipMgr {
	return ri.RelationshipMgr
}

//Run runs the simulation
func (ri *RunnerInfo) Run() Results {
	results := Results{
		Iterations:    ri.Iterations,
		Colors:        make([][]int, ri.Iterations),
		Conversations: make([]int, ri.Iterations),
	}
	//Seed rand to make sure random behaviour is evenly distributed
	rand.Seed(time.Now().UnixNano())

	n := ri.RelationshipMgr
	agents := n.Agents()

	for i := 0; i < ri.Iterations; i++ {
		hold := make(chan bool)
		convCount := make(chan int)

		nc := len(agents)

		for _, a := range agents {
			agent := a
			go func() {
				<-hold
				r := rand.Intn(10)
				time.Sleep(time.Duration(r) * time.Nanosecond)
				convCount <- agent.SendMail(n)
			}()
		}
		close(hold)

		convTotal := 0
		for n := nc; n > 0; n-- {
			convTotal = convTotal + <-convCount
		}
		close(convCount)

		colorCounts := make([]int, n.MaxColors())
		for _, a := range agents {
			color := a.ReadMail(n)
			colorCounts[color]++
		}
		results.Colors[i] = colorCounts
		results.Conversations[i] = convTotal
	}

	return results
}
