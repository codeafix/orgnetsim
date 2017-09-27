package orgnetsim

//RunSim runs the simulation
func RunSim() {
	n, _ := NewNetwork("")

	colors := make([][]int, 100, 100)
	conversations := make([]int, 100, 100)

	for i := 0; i < 100; i++ {
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
}
