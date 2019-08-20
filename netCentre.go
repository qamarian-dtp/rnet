package rnet

import (
	"errors"
	"fmt"
	"sync"
)

func New () (*NetCentre) {
	return &NetCentre {sync.Map {}}
}

type NetCentre struct {
	ppos sync.Map /*
		This map keeps track of the IDs being used by PPOs of the network.
		key: ppo id; val: *PPO

		*/
}

// NewPPO () helps create a new personal post office (PPO).
//
// Inputs
//
// input 0: A desired ID for the PPO. Value may be any string, but not an empty string.
//
// Outpts
//
// outpt 0: The PPO. If an error should occur during the creation, value would be nil.
//
// outpt 1: Possible error values include: NcErrIdInUse.
func (nc *NetCentre) NewPPO (id string) (*PPO, error) {
	// Input data validation. { ...
	if id == "" {
		return nil, errors.New ("A PPO's ID can not be an empty string.")
	}
	// ... }
	// Checking if ID is already in use. { ...
	_, ok := nc.ppos.Load (id)
	if ok == true {
		return nil, NcErrIdInUse
	}
	// ...}
	// Creating a new PPO. { ...
	ppo, errX := newPPO (id, nc)
	if errX != nil {
		errMssg := fmt.Sprintf ("Unable to create new a PPO. [%s]",
			errX.Error ())
		return nil, errors.New (errMssg)
	}
	nc.ppos.Store (id, ppo)
	// ... }
	return ppo, nil
}

// Disconnect () disconnects a PPO from the network. The input of this method should be
// the ID of the PPO that should be disconnected.
func (nc *NetCentre) Disconnect (id string) (error) {
	// Trying to get the PPO, using the ID provided. { ...
	somePPO, okX := nc.ppos.Load (id)
	if okX == false {
		return nil
	}
	// ... }
	// Disconnecting the PPO. { ...
	nc.ppos.Delete (id)
	ppo, okY := somePPO.(*PPO)
	if okY == false {
		return errors.New ("Assertion on PPO failed.")
	}
	errX := ppo.destroy ()
	if errX != nil {
		errMssg := fmt.Sprintf ("The PPO could not be successfully destroyed. " +
			"[%s]", errX.Error ())
		return errors.New (errMssg)
	}
	// ... }
	return nil
}

// provideMDI () provides an MDI that could be used to send messages to another PPO.
//
//	Inputs
//	input 0: The ID of the PPO messages would be sent to.
//
//	Outpts
//	outpt 0: An MDI that could be used to message the PPO specified as input 0.
//
//	outpt 1: On success, value would be nil. On failure value would be an error. If
//		the ID provided as input 0 is not in use, value would be NcErrIdNotInUse.
func (n *NetCentre) provideMDI (id string) (*mdi, error) {
	// Fetching the PPO whose ID was provided. { ...
	somePPO, ok := n.ppos.Load (id)
	if ok == false {
		return nil, NcErrIdNotInUse
	}
	// ... }
	// Requesting an MDI from the PPO. { ...
	ppo, okX := somePPO.(*PPO)
	if okX == false {
		return nil, errors.New ("Assertion on PPO failed.")
	}
	mdi, errX := ppo.getMDI ()
	if errX == PpoErrNotConnected {
		return nil, NcErrIdNotInUse
	} else if errX != nil {
		errMssg := fmt.Sprintf ("Unable to get message delivery info from " +
			"recipient. [%s]", errX.Error ())
		return nil, errors.New (errMssg)
	}
	// ... }
	return mdi, nil
}

var (
	NcErrIdInUse error = errors.New ("PPO ID already in use.")
	NcErrIdNotInUse error = errors.New ("PPO ID not in use.")
)
