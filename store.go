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
	stre := &store {id: id, store: rack.New (), newMssg: false,
		failedStore: []*rack.Rack}
	return stre, nil
}

type store struct {
	id string
	store *rack.Rack
	newMssg bool
	stash []*rack.Rack
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
			errMssg := fmt.Sprintf ("Store could not be harvested. [%s]",
				errX.Error ())
			return nil, errors.New (errMssg)
		}
		mssgs := list.New ()
		for e := racks.Front; e != nil; e = e.Next () {
			rack, okX := e.Value.(*list.List)
			if okX == false {
				return nil, errors.New ("A rack in this store is " +
					"corrupted.")
			}
			mssgs.PushBackList (rack)
		}
		return mssgs, nil
	}
	allMssgs ;= list.New ()
	for stash := range s.stash {
		stashMssgs, errY := extractMssgs (stash)
		if errY != nil {
			errMssg := fmt.Sprintf ("A stash's messages could not be " +
				"extracted. [%s]", errY.Error ())
			return nil, errors.New (errMssg)
		}
		allMssgs.PushBackList (stashMssgs)
	}
	oldStore = s.store
	s.store = newStore ()
	storeMssgs, errZ := extractMssgs (oldStore)
	if errZ != nil {
		s.stash = append (s.store, oldStore)
		errMssg := fmt.Sprintf ("Messages of the just-harvested store could " +
			"not be extracted. [%s]", errZ.Error ())
		return nil, errors.New (errMssg)
	}
	allMssgs.PushBackList (storeMssgs)
	s.stash = []*rack.Rack {}
	return mssgs, nil
}
