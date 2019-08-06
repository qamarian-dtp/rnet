package rnet

import (
	"container/list"
	"errors"
	"runtime"
	"sync/atomic"
)

func newRack () (*rack) {
	rk := &rack {}
	rk.state = RckStateNew
	rk.rack = list.List {}
	return rk
}

type rack struct {
	state int32
	rack list.List
}

func (r *rack) activate () (error) {
	if r.state != RckStateNew {
		return errors.New ("This rack is not a new rack.")
	}
	r.state = RckStateNotInUse
	return nil
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

func (r *rack) harvest ()(*list.List, error) {
	harvestBeginning:
	
	runtime.Gosched ()
	switch r.state {
		case RckStateNew:
			return nil, RckErrToBeActivated
		case RckStateNotInUse:
			okY := atomic.CompareAndSwapInt32 (&r.state, RckStateNotInUse, RckStateHarvested)
			if okX == false {
				goto harvestBeginning
			}
			return &r.rack, nil
		case RckStateInUse:
			goto harvestBeginning
		case RckStateToBeHarvested:
			return nil, RckErrToBeHarvested
		case RckErrBeenHarvested:
			return nil, RckErrBeenHarvested
		default:
			return nil, RckErrInvalidState
	}
}

var (
	RckStateNew           int32 = 0
	RckStateNotInUse      int32 = 1
	RckStateInUse         int32 = 2
	RckStateToBeHarvested int32 = 3
	RckStateHarvested     int32 = 4

	RckErrToBeActivated error = errors.New ("This rack is yet to be activated.")
	RckErrNotNew        error = errors.New ("This rack is not a new rack.")
	RckErrToBeHarvested error = errors.New ("This rack is about to be " +
		"harvested.")
	RckErrBeenHarvested error = errors.New ("This rack has been harvested.")
	RckErrInvalidState  error = errors.New ("This rack is in an invalid state.")
)
