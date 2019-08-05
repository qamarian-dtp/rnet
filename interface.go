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

// ---------- Section A ---------- //
// Ambassadors

func (i *Interface) init (underlyingNet *Network, user, netAddr string) (error) {
	var errX error
	i.id, errX = str.UniquePredsafeStr (32)
	if errX != nil {
		errMssg := fmt.Sprintf ("Unable to generate ID for interface. [%s]",
			errX.Error ())
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