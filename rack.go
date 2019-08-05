package rnet

import (
	"container/list"
	"errors"
	"sync/atomic"
)

type rack struct {
	state int32
	rack list.List
}

func newRack () (*rack) {
	rk := &rack {}
	rk.state = RckStateNew
	rk.rack = list.List {}
	return rk
}

func (r *rack) activate () {
	r.state = RckStateNotInUse
}

func (r *rack) addMssg (m interface {}) (error) {
	ok := atomic.CompareAndSwapInt32 (&r.state, RckStateNotInUse, RckStateInUse)
	if ok == false && r.state == RckStateToBeHarvested {
		return RckErrToBeHarvested
	} else if ok == false {
		return errors.New ("Message could not be added to the rack.")
	}
	if r.rack.Len () == 0 {
		r.rack.PushFront (m)
	} else {
		r.rack.PushBack (m)
	}
	atomic.CompareAndSwapInt32 (&r.state, RckStateInUse, RckStateNotInUse)
}

func (r *rack) getState () (int32) {
	return r.state
}

var (
	RckStateNew           int32 = 0
	RckStateNotInUse      int32 = 1
	RckStateInUse         int32 = 2
	RckStateToBeHarvested int32 = 3
	RckErrToBeHarvested error = errors.New ("This rack is about to be " +
		"harvested.")
)
