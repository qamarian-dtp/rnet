package rnet

import (
	"errors"
	"fmt"
)

func newMDInfo(recipientIntf *Interface) (*dInfo) {
	rk := newRack ()
	return &dInfo {recipientIntf, rk}
}

type mDInfo struct {
	recipientIntf *Interface
	senderRack *rack
}

func (mdi *mDInfo) sendMssg (mssg interface {}) (error) {
	if mdi.senderRack.getState () == RckStateNew {
		errM := mdi.recipientIntf.getStore ().addRack (mdi.senderRack)
		if errM != nil {
			errMssg := fmt.Sprintf ("Unable to add new rack to the " +
				"store. [%s]", errM.Error ())
			return errors.New (errMssg)
		}
		mdi.senderRack.activate ()
	}

	oldRack := mdi.senderRack

	addBeginning:

	if mdi.recipientInt.getNetAddr () == "" {
		return MdiErrNotConnected
	} else if mdi.recipientInt.getClosedSig () == true {
		return MdiErrClosed
	}
	errX := mdi.senderRack.addMssg (mssg)
	if errX == RckErrToBeHarvested {
		mdi.senderRack = newRack ()
		errY := mdi.recipientIntf.getStore ().addRack (mdi.senderRack)
		if errY != nil {
			mdi.senderRack = oldRack
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
	MdiErrNotConnected error = errors.New ("Recipient is not connected to the " +
		"network.")
	MdiErrClosed error = errors.New ("Recipient is closed to new messages.")
)