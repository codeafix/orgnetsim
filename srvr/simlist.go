package srvr

//SimItem contains the information used to display the simulation list and holds a
//relative path to the directory containing all the simulation results
type SimItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
}

//SimList is the list of simulations in the root directory
type SimList struct {
	TimestampHolder
	Items []SimItem
}
