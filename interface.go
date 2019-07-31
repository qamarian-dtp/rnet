package rnet

import (
	"container/list"
	"errors"
	"fmt"
	"gopkg.in/qamarian-lib/str.v1"
	"sync"
)

type Interface struct {
	id string

	underlyingNet *Network
	user string
	netAddr string

	closed bool
	harvestBasket list.List
	deliveryStore store

	cache spCache
}

var (
	IntErrNotConnected error = errors.New ("Interface has been disconnected.")
)

type store struct {
	id string
	racks list.List
	wakeupSignal struct {
		waiting bool
		signalChan sync.Cond
	}
}

func (s *store) AddToRack () {}

type storeProtected struct {
	underlyingStore *store
	sendersAddr string
	lastKnownStore string
	rack list.List
}

func (s *storeProtected) AddToRack () {}

type spCache struct {
	locker sync.RWMutex
	storeP map[string]*storeProtected
}

// ---------- Section A ---------- //
// Comm from network

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
		underlyingStore: &(i.deliveryStore),
		senderAddr: netAddr,
		lastKnownStore: "",
		rack: nil,
	}
	return sp, nil
}

func (i *Interface) releaseAddr () {
	i.netAddr = ""
}

// ---------- Section B ---------- //
// Comm to network

func (i *Interface) getSPOfAnother () {
	//
}

// ---------- Section C ---------- //

func (i *Interface) Open () {}

func (i *Interface) Opened () () {}

func (i *Interface) Send () {}

func (i *Interface) Read () {}

func (i *Interface) Check () {}

func (i *Interface) Close () {}

func (i *Interface) NewIntf () {}

func (i *Interface) GetUser () {}

func (i *Interface) GetAddr () {}

func (i *Interface) Disconnect () {
	i.underlyingNet.disconnectMe (i.netAddr)
}
