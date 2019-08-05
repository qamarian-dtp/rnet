package main

import (
	"errors"
	"fmt"
)

func newDInfo(recipientIntf *Interface) (*dInfo) {
	rk := newRack ()
	return &dInfo {recipientIntf, rk}
}

type dInfo struct {
	recipientIntf *Interface
	senderRack *rack
}

func (di *dInfo) sendMssg (mssg interface {}) (error) {
	if di.senderRack.getState () == RckStateNew {
		errM := di.recipientIntf.getStore ().addRack (di.senderRack)
		if errM != nil {
			errMssg := fmt.Sprintf ("Unable to add new rack to the " +
				"store. [%s]", errM.Error ())
			return errors.New (errMssg)
		}
		di.senderRack.activate ()
	}

	oldRack := di.senderRack

	addBeginning:

	if di.recipientInt.getNetAddr () == "" {
		return DInErrNotConnected
	} else if di.recipientInt.getClosedSig () == true {
		return DInErrClosed
	}
	errX := di.senderRack.addMssg (mssg)
	if errX == RckErrToBeHarvested {
		di.senderRack = newRack ()
		errY := di.recipientIntf.getStore ().addRack (di.senderRack)
		if errY != nil {
			di.senderRack = oldRack
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
	DInErrNotConnected error = errors.New ("Recipient is not connected to the " +
		"network.")
	DInErrClosed error = errors.New ("Recipient is closed to new messages.")
)
