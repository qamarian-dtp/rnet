package rnet

import (
	"errors"
	"gopkg.in/qamarian-lib/str.v1"
	"sync"
)

func New () (*Network, error) {
	id, errX := str.UniquePredsafeStr (32)
	if errX != nil {
		errMssg := fmt.Sprintf ("Unable to generate ID for network. [%s]", errX.Error ())
		return nil, errors.New (errMssg)
	}
	net := Network {}
	return &net, nil
}

type Network struct {
	id string
	locked struct
		locker sync.RWMutex
		locked bool
	}
	freezed struct {
		locker sync.RWMutex
		freezed bool
	}
	allocations sync.Map // KEY: net-addr; VAL: user, net-interface, expirey-signal
	userAddr struct {
		locker sync.RWMutex
		addr map[string][]string
	}
	reservedAddr struct {
		locker sync.Mutex
		addr map[string]bool
	}
}

func (n *Network) NewIntf (userID, netAddr) (*Interface, error) {
	id, errX := str.UniquePredsafeStr (32)
	if errX != nil {
		errMssg := fmt.Sprintf ("Unable to generate ID for interface. [%s]", errX.Error ())
		return nil, errors.New (errMssg)
	}
	intf := Interface {
		underlyingNetwork: n,
		interfaceID: id
		user: userID
}

func (n *Network) GetUser () {}

func (n *Network) GetAddr () {}

func (n *Network) GetStoreP () {}

// --------------

func (n *Network) CheckUsage (netAddr string) (string) {}

func (n *Network) Reserve () {}

func (n *Network) ScanReserve () {}

func (n *Network) Release (netAddr string) (error) {
	n.reservedAddr.Lock ()
	delete (n.reservedAddr.addr, netAddr)
	n.reservedAddr.Unlock ()
}

func (n *Network) Reclaim () {}

// -------- completed -------------

func (n *Network) Lock () (*Unlocker) {
	n.locked.Lock ()
	n.locked.locked = true
	return &Unlocker {locker: sync.RWMutex {}, underlyingNet: n, used: false}
}

func (n *Network) Locked () (bool) {
	return n.locked.locked
}

type Unlocker struct {
	locker sync.RWMutex
	underlyingNet *Network
	used bool
}

func (u *Unlocker) Unlock () {
	u.locker.Lock ()
	defer u.locker.Unlock ()
	if u.used == true {
		return
	}
	u.underlyingNet.locked.locked = false
	u.underlyingNet.locked.locker.Unlock ()
	u.used = true
}

func (n *Network) Freeze () (*Unfreezer) {
	n.freezed.Lock ()
	n.freezed.freezed = true
	return &Unfreezer {locker: sync.RWMutex {}, underlyingNet: n, used: false}
}

func (n *Network) Freezed () (bool) {
	return n.freezed.freezed
}

type Unfreezer struct {
	locker sync.RWMutex
	underlyingNet *Network
	used bool
}

func (u *Unfreezer) Unfreeze () {
	u.locker.Lock ()
	defer u.locker.Unlock ()
	if u.used == true {
		return
	}
	u.underlyingNet.freezed.freezed = false
	u.underlyingNet.freezed.locker.Unlock ()
	u.used = true
}
