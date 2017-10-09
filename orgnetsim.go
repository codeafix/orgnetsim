package orgnetsim

//RunSim runs the simulation
func RunSim(n RelationshipMgr, iterations int) ([][]int, []int) {
	colors := make([][]int, iterations+1, iterations+1)
	conversations := make([]int, iterations+1, iterations+1)

	colorCounts := make([]int, n.MaxColors(), n.MaxColors())
	agents := n.Agents()
	for _, a := range agents {
		colorCounts[a.GetColor()]++
	}
	colors[0] = colorCounts

	for i := 1; i <= iterations; i++ {
		hold := make(chan bool)
		convCount := make(chan int)

		nc := len(agents)

		for _, a := range agents {
			agent := a
			go func() {
				<-hold
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
		colors[i] = colorCounts
		conversations[i] = convTotal
	}

	return colors, conversations
}
