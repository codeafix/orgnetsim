package srvr

import (
	"github.com/codeafix/orgnetsim/sim"
)

//SimStep holds the results of each simulation step
type SimStep struct {
	TimestampHolder
	sim.Network
	sim.Results
}
