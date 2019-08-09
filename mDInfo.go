package rnet

import (
	"errors"
	"fmt"
)

func newMDInfo(recipientIntf *Interface) (*mDInfo) {
	rk := newRack ()
	return &mDInfo {recipientIntf, rk}
}

type mDInfo struct {
	recipientIntf *Interface
	senderRack *rack
}

func (mdi *mDInfo) sendMssg (mssg interface {}) (error) {
	addBeginning:
	recipientStore := mdi.recipientIntf.getStore ()
	if recipientStore == nil {
		return MdiErrNotConnected
	}
	errX := mdi.senderRack.addMssg (mssg)
	if errX == RckErrBeenHarvested {
		oldRack := mdi.senderRack
		mdi.senderRack = newRack ()
		errY := recipientStore.addRack (mdi.senderRack)
		if errY != nil {
			mdi.senderRack = oldRack
			errMssg := fmt.Sprintf ("Unable to add new rack to the " +
				"store. [%s]", errY.Error ())
			return errors.New (errMssg)
		}
		goto addBeginning
	} else if errX != nil {
		errMssg := fmt.Sprintf ("Unable to add message to the store. " +
			"[%s]", errX.Error ())
		return errors.New (errMssg)
	}
	recipientStore.sigNewMssg ()
	return nil
}

var (
	MdiErrNotConnected error = errors.New ("Recipient is not connected to the " +
		"network.")
	MdiErrClosed error = errors.New ("Recipient is closed to new messages.")
)
