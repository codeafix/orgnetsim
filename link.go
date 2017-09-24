package orgnetsim

//A Link between Agents in the Network
type Link struct {
	Agent1ID string  `json:"agent1Id"`
	Agent2ID string  `json:"agent2Id"`
	Strength float64 `json:"strength"`
}
