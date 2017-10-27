package srvr

import "time"

//TimestampHolder holds the current timestamp on a persistable object
type TimestampHolder struct {
	Stamp time.Time `json:"-"`
}

//Timestamp of the persistable object
func (ts *TimestampHolder) Timestamp() time.Time {
	return ts.Stamp
}

//UpdateTimestamp of the persistable object
func (ts *TimestampHolder) UpdateTimestamp(t time.Time) {
	ts.Stamp = t
}

//Persistable must be implemented by any obj that requires to be saved in a file using
//FileUpdater
type Persistable interface {
	Timestamp() time.Time
	UpdateTimestamp(t time.Time)
}
