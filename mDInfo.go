package rnet

import (
	"errors"
	"fmt"
	"runtime"
)

func newMDInfo (recipientIntf *Interface) (*mDInfo, error) {
	rk := newRack ()
	addBeginning:
	recipientStore := recipientIntf.getStore ()
	if recipientStore == nil {
		return nil, mdiErrNotConnected
	}
	errX := recipientStore.addRack (rk)
	if errX == strErrBeenHarvested {
		runtime.Gosched ()
		goto addBeginning
	} else if errX != nil {
		errMssg := fmt.Sprintf ("The rack created for the MDI could not be " +
			"added to the recipient's store. [%s]", errX.Error ())
		return nil, errors.New (errMssg)
	}
	return &mDInfo {recipientIntf, rk}, nil
}

type mDInfo struct {
	recipientIntf *Interface
	senderRack *rack
}

func (mdi *mDInfo) sendMssg (mssg interface {}) (error) {
	addBeginning:
	recipientStore := mdi.recipientIntf.getStore ()
	if recipientStore == nil {
		return mdiErrNotConnected
	}
	errX := mdi.senderRack.addMssg (mssg)
	if errX == rckErrBeenHarvested {
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
	mdiErrNotConnected error = errors.New ("Recipient is not connected to the " +
		"network.")
	mdiErrClosed error = errors.New ("Recipient is closed to new messages.")
)
