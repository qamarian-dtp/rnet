package rnet

import (
	"container/list"
	"errors"
	"fmt"
	"gopkg.in/qamarian-lib/str.v1"
	"sync"
	"sync/atmoic"
)

type Interface struct {
	id string

	underlyingNet *Network
	user string
	netAddr string

	closed bool
	harvestBasket list.List
	deliveryStore *store

	cache spCache
}

func (i *Interface) getID () (string) {
	return i.id
}

func (i *Interface) getUNet () (*Network) {
	return i.underlyingNet
}

func (i *Interface) getUser () (string) {
	return i.user
}

func (i *Interface) getNetAddr () (string) {
	return i.netAddr
}

func (i *Interface) getClosedSig () (bool) {
	return i.closed
}

func (i *Interface) getStore () (*store) {
	return i.deliveryStore
}

type store struct {
	id string
	racks struct {
		state int32 // 0: not in use; 1: about to-be manipulated; 2: about to-be harvested
		racks list.List
	}
	wakeupSignal struct {
		waiting bool
		signalChan sync.Cond
	}
}

func (s *store) getID () (string) {
	return s.id
}

func (s *store) lockRacks () {
	for {
		ok := atomic.CompareAndSwapInt32 (*s.racks.state, 0, 1)
		if ok == true {
			break
		}
	}
}

func (s *store) unlockRacks () {
	s.racks.state = 0
}

func (s *store) addRack (senderRack *rack) (error) {
	if s.racks.state == 2 {
		return StrErrToBeHarvested
	}
	if s.racks.Len () == 0 {
		s.racks.PushFront (senderRack)
	} else {
		s.racks.PushBack (senderRack)
	}
	s.racks.state = 0
}

var (
	StrErrToBeHarvested error = errors.New ("The store is about to be harvested.")
)

type spCache struct {
	locker sync.RWMutex
	storeP map[string]*storeProtected
}

type storeProtected struct {
	recipientInt *Interface
	senderAddr string
	lastKnownStore string
	senderRack *list.List
}

func (s *storeProtected) addMessage (mssg interface {}) (error) {
	if s.recipientInt.getNetAddr () == "" {
		return StpErrNotConnected
	} else if s.recipientInt.getClosedSig () == true {
		return StpErrClosed
	}
	s.recipientInt.getStore ().lockRacks ()
	defer s.recipientInt.getStore ().unlockRacks ()
	if lastKnownStore != s.recipientInt.getStore ().getID () {
		s.senderRack = &list.List {}
		errX := s.recipientInt.getStore ().addRack (s.senderRack)
		if errX != nil {
			errMssg := fmt.Sprintf ("Message could not be added to the rack. [%s]",
				errX.Error ())
			return errors.New (errMssg)
		}
	}
	if s.senderRack.Len () == 0 {
		s.senderRack.PushFront (mssg)
	} else {
		s.senderRack.PushBack (mssg)
	}
}

var (
	StpErrNotConnected error = errors.New ("Recipient is not connected to the network.")
	StpErrClosed error = errors.New ("Recipient is closed to new messages.")
)

// ---------- Section A ---------- //
// Ambassadors

func (i *Interface) init (underlyingNet *Network, user, netAddr string) (error) {
	var errX error
	i.id, errX = str.UniquePredsafeStr (32)
	if errX != nil {
		errMssg := fmt.Sprintf ("Unable to generate ID for interface. [%s]", errX.Error ())
 		return errors.New (errMssg)
	}

	i.underlyingNet = underlyingNet
	i.user = user
	i.netAddr = netAddr

	i.closed = false
	i.harvestBasket = list.List {}
	i.deliveryStore = store {}

	i.cache = spCache {
		locker: sync.RWMutex {},
		storeP: make (map[string]*storeProtected),
	}

	return nil
}

func (i *Interface) getUser () (string) {
	return i.user
}

func (i *Interface) provideSP (netAddr string) (*storeProtected, error) {
	if i.netAddr == "" {
		return nil, IntErrNotConnected
	}
	sp := &storeProtected {
		underlyingInt: i,
		senderAddr: netAddr,
		lastKnownStore: "",
		rack: list.List {},
	}
	return sp, nil
}

var (
	IntErrNotConnected error = errors.New ("Interface has been disconnected.")
)

func (i *Interface) releaseAddr () {
	i.netAddr = ""
}

// ---------- Section B ---------- //
// Originals

func (i *Interface) Open () {
	i.closed = true
}

func (i *Interface) Opened () (bool) {
	return i.closed
}

func (i *Interface) Send () {}

func (i *Interface) Read () {}

func (i *Interface) Wait () {}

func (i *Interface) Close () {
	i.closed = false
}
