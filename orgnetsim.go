package orgnetsim

import "fmt"

//RunSim runs the simulation
func RunSim() {
	iterations := 500
	n, _ := GenerateHierarchy()

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

	for c := 0; c < MaxColors; c++ {
		fmt.Printf("%s,", Color(c).String())
	}
	fmt.Printf("Conversations\n")
	for i := 0; i <= iterations; i++ {
		for j := 0; j < MaxColors; j++ {
			fmt.Printf("%d,", colors[i][j])
		}
		fmt.Printf("%d\n", conversations[i])
	}
}
