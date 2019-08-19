package rnet

import (
	"errors"
	"fmt"
	"runtime"
)

// newMDI () creates a new MDI that could be used to communicate with a PPO.
func newMDI (recipientPPO *PPO) (*mdi, error) {
	rk := newRack ()
	addBeginning:
	store := recipientPPO.getStore ()
	if store == nil {
		return nil, mdiErrNotConnected
	}
	errX := store.addRack (rk)
	if errX == strErrBeenHarvested {
		runtime.Gosched ()
		goto addBeginning
	} else if errX != nil {
		errMssg := fmt.Sprintf ("The rack created for the MDI could not be added to the recipient's store. [%s]", errX.Error ())
		return nil, errors.New (errMssg)
	}
	return &mdi {recipientPPO, rk}, nil
}

type mdi struct {
	recipientPPO *PPO  // The PPO messages would be sent to.
	senderRack *rack   // The rack of the sender, in the PPO's store.
}

func (mdi *mdi) getRecipientPPO () (*PPO) {
	return mdi.recipientPPO
}

func (mdi *mdi) getSenderRack () (*rack) {
	return mdi.senderRack
}

func (mdi *mdi) newRack () (error) {
	newSenderRack := newRack ()
	store := mdi.recipientPPO.getStore ()
	if store == nil {
		return mdiErrNotConnected
	}
	errY := store.addRack (newSenderRack)
	if errY != nil {
		errMssg := fmt.Sprintf ("Unable to add the new rack to the recipient's store. [%s]", errY.Error ())
		return errors.New (errMssg)
	}
	mdi.senderRack = newSenderRack
	return nil
}

var (
	mdiErrNotConnected error = errors.New ("Recipient is not on this network.")
)
