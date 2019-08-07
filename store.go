package rnet

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/qamarian-dtp/rack.v0"
	"sync"
)

func newStore () (*store, error) {
	id, errX := str.UniquePredsafeStr (32)
	if errX != nil {
		errMssg := fmt.Sprintf ("ID could not be generated for store. [%s]",
			errX.Error ())
		return nil, errors.New (errMssg)
	}
	stre := &store {id: id, store: rack.New (), newMssg: false}
	return stre, nil
}

type store struct {
	id string
	store *rack.Rack
	newMssg bool
}

func (s *store) checkNewMssg () (bool) {
	return s.newMssg
}

func (s *store) sigNewMssg () {
	s.newMssg = true
}

func (s *store) Harvest () (*list.List, error) {
	func extractMssgs (s *store) (*list.List, error) {
		racks, errX := s.store.Harvest ()
		if errX != nil {
			return nil, errX
		}
		mssgs := list.New ()
		for e := racks.Front; e != nil; e = e.Next () {
			rack, okX := e.Value.(*list.List)
			if okX == false {
				return nil, errors.New ("A box in this store is " +
					"corrupted.")
			}
			mssgs.PushBackList (rack)
		}
		return mssgs, nil
	}
	mssgs, errZ := extractMssgs (s.store)
	if errZ != nil {
		errMssg := fmt.Sprintf ("This store's messages could not be " +
			"harvested. [%s]", errZ.Error ())
		return nil, errors.New (errMssg)
	}
	return mssgs, nil
}