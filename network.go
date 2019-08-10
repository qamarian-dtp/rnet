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
		errMssg := fmt.Sprintf ("Unable to generate ID for network. [%s]",
			errX.Error ())
		return nil, errors.New (errMssg)
	}
	net := Network {id: netID}
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
	allocations struct {
		locker sync.Mutex
		alloc sync.Map // KEY: net-addr; VAL: *interface
	}
}

func (n *Network) NewIntf (netAddr string) (*Interface, error) {
	if netAddr == "" {
		return nil, errors.New ("Network address can not be an empty string.")
	}
	_, ok := n.allocations.alloc.Load (netAddr)
	if ok == true {
		return nil, NetErrInUse
	}
	i, errX := newIntf (netAddr, n)
	if errX != nil {
		errMssg := fmt.Sprintf ("Unable to create new interface. [%s]",
			errX.Error ())
		return nil, errors.New (errMssg)
	}
	n.allocations.alloc.Store (netAddr, i)
	return i, nil
}

func (n *Network) Disconnect (netAddr string) (error) {
	alloc, okX := n.allocations.alloc.Load (netAddr)
	if okX == false {
		return NetErrNotInUse
	}
	n.allocations.alloc.Delete (netAddr)
	addrUser, okY := alloc.(*Interface)
	if okY == false {
		return errors.New ("Address-allocation-data value could not be " +
			"treated as an interface.")
	}
	errX := addrUser.destroy ()
	if errX != nil {
		errMssg := fmt.Sprintf ("The interface using the address could " +
			"not be destroyed. [%s]", errX.Error ())
		return errors.New (errMssg)
	}
	return nil
}

func (n *Network) provideMDInfo (netAddr string) (*mDInfo, error) {
	alloc, ok := n.allocations.alloc.Load (netAddr)
	if ok == false {
		return nil, NetErrNotInUse
	}
	intf, _ := alloc.(*Interface)
	di, errX := intf.getMDInfo ()
	if errX == IntErrNotConnected {
		return nil, NetErrNotInUse
	}
	if errX != nil {
		errMssg := fmt.Sprintf ("Unable to get message-delivery info " +
		 	"from recipient. [%s]", errX.Error ())
		return nil, errors.New (errMssg)
	}
	return di, nil
}

var (
	NetErrLocked error = errors.New ("Interface creation not allowed: network " +
		"currently locked.")
	NetErrInUse error = errors.New ("Network address already in use.")
	NetErrNotInUse error = errors.New ("Network address not in use.")
)
