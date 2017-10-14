package orgnetsim

//A Link between Agents in the Network
type Link struct {
	Agent1ID string `json:"source"`
	Agent2ID string `json:"target"`
	Strength int    `json:"strength"`
}
