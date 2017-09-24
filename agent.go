package orgnetsim

//An Agent is a node in the network
type Agent struct {
	ID             string  `json:"id"`
	Color          Color   `json:"color"`
	Susceptability float64 `json:"susceptability"`
	Influence      float64 `json:"influence"`
	Contrariness   float64 `json:"contrariness"`
}
