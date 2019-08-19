package rnet

import (
	"container/list"
	"errors"
	"fmt"
	"gopkg.in/qamarian-dtp/cart.v1"
)

// newStore () simply creates a new delivery store.
func newStore () (*store, error) {
	racks, racksManager := cart.New ()
	stre := &store {
		racks: racks,
		racksManager: racksManager,
		newMssg: false,
	}
	return stre, nil
}

type store struct {
	racks *cart.Cart               // The racks in the store.
	racksManager *cart.AdminPanel  // A data you could to harvest the racks in a store.
	newMssg bool                   /* Indicates if there is a new message in the store or not. True means there's a new message, while false means there's no message in the store. */
}

// addRack () adds a rack to the store.
func (s *store) addRack (r *rack) (error) {
	errX := s.racks.Put (r)
	if errX == cart.ErrBeenHarvested {
		return strErrBeenHarvested
	} else if errX != nil {
		errMssg := fmt.Sprintf ("Unable to the add rack. [%s]", errX.Error ())
		return errors.New (errMssg)
	}
	return nil
}

// checkNewMssg () checks if there's a new message in the store.
func (s *store) checkNewMssg () (bool) {
	if s == nil {
		return false
	}
	return s.newMssg
}

// mssgAdded () can be used to inidicate that a new message has been added to a rack in the store.
func (s *store) mssgAdded () {
	if s == nil {
		return
	}
	s.newMssg = true
}

// Harvest () harvests all messages in all racks in the store.
func (s *store) Harvest () (*list.List, error) {
	// This function extracts the messages in a given store.
	extractMssgs := func (s *store) (*list.List, error) {
		// Harvesting the racks in the store. { ...
		racks, errX := s.racksManager.Harvest ()
		if errX != nil {
			return nil, errX
		}
		// ... }
		// Harvesting the messages in the racks. { ...
		mssgs := list.New ()
		for e := racks.Front (); e != nil; e = e.Next () {
			rack, okX := e.Value.(*rack)
			if okX == false {
				return nil, errors.New ("A rack in this store is corrupted.")
			}
			rackMssgs, errD := rack.harvest ()
			if errD != nil {
				errMssg := fmt.Sprintf ("A rack could not be harvested. [%s]", errD.Error ())
				return nil, errors.New (errMssg)
			}
			mssgs.PushBackList (rackMssgs)
		}
		// ... }
		return mssgs, nil
	}

	mssgs, errZ := extractMssgs (s)
	if errZ == cart.ErrBeenHarvested {
		return nil, strErrBeenHarvested
	} else if errZ != nil {
		errMssg := fmt.Sprintf ("This store's messages could not be harvested. [%s]", errZ.Error ())
		return nil, errors.New (errMssg)
	}
	s.newMssg = false
	return mssgs, nil
}

var (
	strErrBeenHarvested error = errors.New ("This store has already been harvested.")
)

