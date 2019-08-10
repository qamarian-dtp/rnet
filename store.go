package rnet

import (
	"container/list"
	"errors"
	"fmt"
	"gopkg.in/qamarian-dtp/cart.v0"
	"gopkg.in/qamarian-lib/str.v1"
)

func newStore () (*store, error) {
	id, errX := str.UniquePredsafeStr (32)
	if errX != nil {
		errMssg := fmt.Sprintf ("ID could not be generated for store. [%s]",
			errX.Error ())
		return nil, errors.New (errMssg)
	}
	racksContainer, manager := cart.New ()
	stre := &store {
		id: id,
		racks: racksContainer,
		racksManager: manager,
		newMssg: false,
	}
	return stre, nil
}

type store struct {
	id string
	racks *cart.Cart
	racksManager *cart.AdminPanel
	newMssg bool
}

func (s *store) addRack (r *rack) (error) {
	errX := s.racks.Put (r)
	if errX == cart.ErrBeenHarvested {
		return StrErrBeenHarvested
	} else if errX != nil {
		errMssg := fmt.Sprintf ("Unable to add rack. [%s]", errX.Error ())
		return errors.New (errMssg)
	}
	return nil
}

func (s *store) checkNewMssg () (bool) {
	if s == nil {
		return false
	}
	return s.newMssg
}

func (s *store) sigNewMssg () {
	if s == nil {
		return
	}
	s.newMssg = true
}

func (s *store) Harvest () (*list.List, error) {
	extractMssgs := func (s *store) (*list.List, error) {
		racks, errX := s.racksManager.Harvest ()
		if errX != nil {
			return nil, errX
		}
		mssgs := list.New ()
		for e := racks.Front (); e != nil; e = e.Next () {
			rack, okX := e.Value.(*rack)
			if okX == false {
				return nil, errors.New ("A rack in this store is " +
					"corrupted.")
			}
			rackMssgs, errD := rack.harvest ()
			if errD != nil {
				errMssg := fmt.Sprintf ("A rack could not be " +
					"harvested. [%s]", errD.Error ())
				return nil, errors.New (errMssg)
			}
			mssgs.PushBackList (rackMssgs)
		}
		return mssgs, nil
	}
	mssgs, errZ := extractMssgs (s)
	if errZ == cart.ErrBeenHarvested {
		return nil, StrErrBeenHarvested
	} else if errZ != nil {
		errMssg := fmt.Sprintf ("This store's messages could not be " +
			"harvested. [%s]", errZ.Error ())
		return nil, errors.New (errMssg)
	}
	return mssgs, nil
}

var (
	StrErrBeenHarvested error = errors.New ("This store has already been harvested.")
)

