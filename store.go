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
	failedStore []rack.Rack
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
		rack1, okX := racks.Front.(*rack.Rack)