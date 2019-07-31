package rnet

import (
	"errors"
	"fmt"
	"gopkg.in/qamarian-lib/str.v1"
	"sync"
)

func New () (*Network, error) {
	netID, errX := str.UniquePredsafeStr (32)
	if errX != nil {
		errMssg := fmt.Sprintf ("Unable to generate ID for network. [%s]", errX.Error ())
		return nil, errors.New (errMssg)
	}
	net := Network {id: netID, locked: false, freezed: false}
	net.allocations = struct {
		locker sync.Mutex
		alloc sync.Map
	}{
		sync.Mutex {},
		sync.Map {},
	}
	return &net, nil
}

type Network struct {
	id string
	locked bool
	freezed bool
	allocations struct {
		locker sync.Mutex
		alloc sync.Map // KEY: net-addr; VAL: *interface
	}
}

// ---------- Section A ---------- //
// Originals

func (n *Network) NewIntf (userID, netAddr string) (*Interface, error) {
	if userID == "" {
		return nil, errors.New ("User ID can not be an empty string.")
	}
	if netAddr == "" {
		return nil, errors.New ("Network address can not be an empty string.")
	}
	if n.locked == true {
		return nil, NetErrLocked
	}
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
	return allok.getUser ()
}

func (n *Network) Disconnect (netAddr string) {
	alloc, ok := n.allocations.alloc.Load (netAddr)
	if ok == false {
		return
	}
	allok, _ := alloc.(*Interface)
	allok.releaseAddr ()
	n.allocations.alloc.Delete (netAddr)
}

func (n *Network) Lock () {
	n.locked = true
}

func (n *Network) Locked () (bool) {
	return n.locked
}

func (n *Network) Unlock () {
	n.locked = false
}

func (n *Network) Freeze () {
	n.freezed = true
}

func (n *Network) Freezed () (bool) {
	return n.freezed
}

func (n *Network) Unfreeze () {
	n.freezed = false
}

// ---------- Section B ---------- //
// Ambassadors

func (n *Network) getSPOfInt (netAddr, myAddr string) (*storeProtected, error) {
	alloc, ok := n.allocations.alloc.Load (netAddr)
	if ok == false {
		return nil, NetErrNotInUse
	}
	intf, _ := alloc.(*Interface)
	sp, errX := intf.provideSP (myAddr)
	if errX == IntErrNotConnected {
		return nil, NetErrNotInUse
	}
	if errX != nil {
		errMssg := fmt.Sprintf ("Unable to get store protected. [%s]", errX.Error ())
		return nil, errors.New (errMssg)
	}
	return sp, nil
}

var (
	NetErrNotInUse error = errors.New ("Network address not in use.")
)
