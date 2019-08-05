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
	rk.state = 0
	rk.rack = list.List {}
	return rk
}

func (r *rack) addMssg (m interface {}) (error) {
	ok := atomic.CompareAndSwapInt32 (&r.state, 0, 1)
	if ok == false && r.state == 2 {
		return RckErrToBeHarvested
	} else if ok == false {
		return errors.New ("This data type is buggy or in use by multiple "
			"routines.")
	}
	if r.rack.Len () == 0 {
		r.rack.PushFront (m)
	} else {
		r.rack.PushBack (m)
	}
	atomic.CompareAndSwapInt32 (&r.state, 1, 0)
}

var (
	RckErrToBeHarvested error = errors.New ("This rack is about to be " +
		"harvested.")
)
