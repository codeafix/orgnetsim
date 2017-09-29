package orgnetsim

//RunSim runs the simulation
func RunSim(n *Network, iterations int) ([][]int, []int) {
	colors := make([][]int, iterations+1, iterations+1)
	conversations := make([]int, iterations+1, iterations+1)

	colorCounts := make([]int, MaxColors, MaxColors)
	for _, a := range n.Nodes {
		colorCounts[a.Color]++
	}
	colors[0] = colorCounts

	for i := 1; i <= iterations; i++ {
		hold := make(chan bool)
		convCount := make(chan int)

		nc := len(n.Nodes)

		for _, a := range n.Nodes {
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

		colorCounts := make([]int, MaxColors, MaxColors)
		for _, a := range n.Nodes {
			color := a.ReadMail(n)
			colorCounts[color]++
		}
		colors[i] = colorCounts
		conversations[i] = convTotal
	}

	return colors, conversations
}
