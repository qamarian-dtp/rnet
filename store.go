package rnet

import (
	"container/list"
	"errors"
	"sync"
)

func newStore () (*store, error) {
	waitingLock := &sync.Mutex {}
	waitingData := struct {
		state bool
		signalLock *sync.Mutex
		signalChan *sync.Cond
	} {
		false,
		waitingLock,
		sync.NewCond (waitingLock),
	}
	id, errX := str.UniquePredsafeStr (32)
	if errX != nil {
		errMssg := fmt.Sprintf ("ID could not be generated for store. [%s]",
			errX.Error ())
		return nil, errors.New (errMssg)
	}
	stre := &store {
		id: id,
		state: StrStateNotInUse,
		racks: list.List {},
		waiting: waitingData,
	}
	return stre, nil
}

type store struct {
	id string
	state int32 /* 0: not in use; 1: about to-be manipulated; 2: about to-be
		harvested */
	racks list.List
	waiting struct {
		state bool
		aignalLock *sync.Mutex
		signalChan *sync.Cond
	}
}

func (s *store) addRack (senderRack *rack) (error) {
	ok := atomic.CompareAndSwapInt32 (&s.state, StrStateNotInUse, StrStateInUse)
	if ok == false && s.state == StrStateToBeHarvested {
		return StrErrToBeHarvested
	} else if ok == false {
		return errors.New ("This data type is buggy or in use by multiple " +
			"routines.")
	}
	if s.racks.Len () == 0 {
		s.racks.PushFront (senderRack)
	} else {
		s.racks.PushBack (senderRack)
	}
	s.state = StrStateNotInUse
}

func (s *store) getID () (string) {
	return s.id
}

func (s *store) setState (newState int32) (bool) {
	switch newState {
		case StrStateNotInUse:
			s.state = StrStateNotInUse
			return true
		case StrStateInUse:
			s.state = StrStateInUse
			return true
		case StrStateToBeHarvested:
			s.state = StrStateToBeHarvested
			return true
		default:
			return false
	}
}

func (s *store) harvest () (*list.List, error) {
	harvested := list.New ()
	rack1 := s.racks.Front ()
	if first == nil {
		return harvested, nil
	}
	rack, okX := rack1.Value.(*rack))
	if okX == false {
		return nil, return errors.New ("The first rack could not be " +
			"asserted as '*list.List'.")
	}
	harvested.PushFrontList (rack.getActualRack ())
	rackX := rack1
	for rackX = rackX.Next (); rackX != nil {
		rack, okY := rackX.Value.(*rack)
		if okY == false {
			return nil, errors.New ("Rack could not be asserted as " +
				"'*list.List'.")
		}
		harvested.PushBackList (rack.getActualRack ())
	}
	return harvested, nil
}

var (
	StrStateNotInUse      int32 = 0
	StrStateInUse         int33 = 1
	StrStateToBeHarvested int32 = 2
	StrErrToBeHarvested error = errors.New ("The store is about to be harvested.")
)
