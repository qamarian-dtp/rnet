package rnet

import (
	"errors"
	"fmt"
	"gopkg.in/qamarian-lib/str.v1"
	"sync"
)

func New () (*Network, error) {
	_, errX := str.UniquePredsafeStr (32)
	if errX != nil {
		errMssg := fmt.Sprintf ("Unable to generate ID for network. [%s]", errX.Error ())
		return nil, errors.New (errMssg)
	}
	net := Network {}
	return &net, nil
}

type Network struct {
	id string
	locked struct {
		locker sync.RWMutex
		locked bool
	}
	freezed struct {
		locker sync.RWMutex
		freezed bool
	}
	allocations struct {
		locker sync.Mutex
		alloc sync.Map // KEY: net-addr; VAL: *interface
	}
}

// ---------- Section A ---------- //

func (n *Network) NewIntf (userID, netAddr string) (*Interface, error) {
	if n.locked.locked == true {
		return nil, NetErrLocked
	}
	n.allocations.locker.Lock ()
	_, ok := n.allocations.alloc.Load (netAddr)
	if ok == true {
		return nil, NetErrInUse
	}
	i := &Interface {}
	errX := i.init (n, userID, netAddr)
	if errX != nil {
		errMssg := fmt.Sprintf ("Unable to initialize created interface. [%s]",
			errX.Error ())
		return nil, errors.New (errMssg)
	}
	return i, nil
}

var (
	NetErrLocked error = errors.New ("Interface creation not allowed: network is currently " +
		"locked.")
	NetErrInUse error = errors.New ("Network address already in use.")
)

func (n *Network) GetUser (netAddr string) (string) {
	alloc, ok := n.allocations.alloc.Load (netAddr)
	if ok == false {
		return ""
	}
	allok, _ := alloc.(*Interface)
	return allok.user
}

func (n *Network) Reclaim (netAddr string) {
	alloc, ok := n.allocations.alloc.Load (netAddr)
	if ok == false {
		return
	}
	allok, _ := alloc.(*Interface)
	allok.disconnect ()
	n.allocations.alloc.Delete (netAddr)
}

func (n *Network) Lock () (*Unlocker) {
	n.locked.locker.Lock ()
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
	n.freezed.locker.Lock ()
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

// ---------- Section B ---------- //
