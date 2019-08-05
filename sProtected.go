package main

import (
	"errors"
	"fmt"
)

func newSProtected (recipientIntf *Interface) (*storeProtected) {
	rk := newRack (RckToBeHarvestedState)
	return storeProtected {recipientIntf, &rk}
}

type storeProtected struct {
	recipientIntf *Interface
	senderRack *rack
}

func (s *storeProtected) addMessage (mssg interface {}) (error) {
	if s.senderRack.getState () == RckStateNew {
		errM := s.recipientIntf.getStore ().addRack (s.senderRack)
		if errM != nil {
			errMssg := fmt.Sprintf ("Unable to add new rack to the " +
				"store. [%s]", errM.Error ())
			return errors.New (errMssg)
		}
		s.senderRack.activate ()
	}

	oldRack := s.senderRack

	addBeginning:

	if s.recipientInt.getNetAddr () == "" {
		return StpErrNotConnected
	} else if s.recipientInt.getClosedSig () == true {
		return StpErrClosed
	}
	errX := s.senderRack.addMssg (mssg)
	if errX == RckErrToBeHarvested {
		s.senderRack = newRack ()
		errY := s.recipientIntf.getStore ().addRack (s.senderRack)
		if errY != nil {
			s.senderRack = oldRack
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
}

var (
	StpErrNotConnected error = errors.New ("Recipient is not connected to the " +
		"network.")
	StpErrClosed error = errors.New ("Recipient is closed to new messages.")
)
