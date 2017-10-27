package srvr

import "github.com/codeafix/orgnetsim/sim"

//SimInfo contains all relevant information about a simulation
type SimInfo struct {
	TimestampHolder
	Steps   []int              `json:"steps"`
	Options sim.NetworkOptions `json:"options"`
}
