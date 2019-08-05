package main

import (
	"container/list"
	"errors"
)

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
	id, errX := str.UniquePredsafeStr (32
	if errX != nil {
		errMssg := fmt.Sprintf ("ID could not be generated for store. [%s]",
			errX.Error ())
		return nil, errors.New (errMssg)
	}
	stre := store {
		id: id,
		state: 0,
		racks: list.List {},
		waiting: waitingData,
	}
	return stre, nil
}

func (s *store) addRack (senderRack *rack) (error) {
	ok := atomic.CompareAndSwapInt32 (&s.state, 0, 1)
	if ok == false && s.state == 2 {
		return StrErrToBeHarvested
	} else if ok == false {
		return errors.New ("This data type is buggy or in use by multiple "
			"routines.")
	}
	if s.racks.Len () == 0 {
		s.racks.PushFront (senderRack)
	} else {
		s.racks.PushBack (senderRack)
	}
	s.state = 0
}

func (s *store) getID () (string) {
	return s.id
}

var (
	StrErrToBeHarvested error = errors.New ("The store is about to be harvested.")
)
