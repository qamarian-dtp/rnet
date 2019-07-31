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
	disconnected bool

	underlyingNet *Network
	user string
	netAddr string

	closed bool
	harvestBasket list.List
	deliveryStore store

	cache spCache
}

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

func (i *Interface) init (underlyingNet *Network, user, netAddr string) (error) {
	var errX error
	i.id, errX = str.UniquePredsafeStr (32)
	if errX != nil {
		errMssg := fmt.Sprintf ("Unable to generate ID for interface. [%s]", errX.Error ())
 		return errors.New (errMssg)
	}
	i.disconnected = false

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

func (i *Interface) getStoreP () {}

func (i *Interface) disconnect () {
	i.disconnected = true
}

// ---------- Section B ---------- //

func (i *Interface) Open () {}

func (i *Interface) Opened () () {}

func (i *Interface) Send () {}

func (i *Interface) Read () {}

func (i *Interface) Check () {}

func (i *Interface) Close () {}

func (i *Interface) ReleaseAddr () {}

func (i *Interface) NewIntf () {}

func (i *Interface) GetUser () {}

func (i *Interface) GetAddr () {}

func (i *Interface) Destroy () {}
