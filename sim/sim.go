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

//RunSim runs the simulation
func RunSim(n RelationshipMgr, iterations int) Results {
	results := Results{
		Iterations:    iterations,
		Colors:        make([][]int, iterations+1, iterations+1),
		Conversations: make([]int, iterations+1, iterations+1),
	}
	//Seed rand to make sure random behaviour is evenly distributed
	rand.Seed(time.Now().UnixNano())

	colorCounts := make([]int, n.MaxColors(), n.MaxColors())
	agents := n.Agents()
	for _, a := range agents {
		colorCounts[a.GetColor()]++
	}
	results.Colors[0] = colorCounts

	for i := 1; i <= iterations; i++ {
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

		colorCounts := make([]int, n.MaxColors(), n.MaxColors())
		for _, a := range agents {
			color := a.ReadMail(n)
			colorCounts[color]++
		}
		results.Colors[i] = colorCounts
		results.Conversations[i] = convTotal
	}

	return results
}
