package main

import (
	"errors"
	"fmt"
)

type storeProtected struct {
	recipientInt *Interface
	senderRack *rack
}

func (s *storeProtected) new (

func (s *storeProtected) addMessage (mssg interface {}) (error) {
	addBeginning:

	if s.recipientInt.getNetAddr () == "" {
		return StpErrNotConnected
	} else if s.recipientInt.getClosedSig () == true {
		return StpErrClosed
	}
	errX := s.senderRack.lock ()
	defer s.senderRack.unlock ()
	if errX == RckErrToBeHarvested {
		s.senderRack = &rack {}
		errY := s.recipientInt.getStore ().addRack (s.senderRack)
		if errY != nil {
			errMssg := fmt.Sprintf ("Unable to add new rack to the " +
				"store. [%s]", errY.Error ())
			return errors.New (errMssg)
		}
		goto addBeginning
	} else if errX != nil {
		errrMssg := fmt.Sprintf ("Unable to add new rack to the store. " +
			"[%s]", errX.Error ())
		return errors.New (errMssg)
	}
	if s.senderRack.Len () == 0 {
		s.senderRack.PushFront (mssg)
	} else {
		s.senderRack.PushBack (mssg)
	}
}

var (
	StpErrNotConnected error = errors.New ("Recipient is not connected to the " +
		"network.")
	StpErrClosed error = errors.New ("Recipient is closed to new messages.")
)
