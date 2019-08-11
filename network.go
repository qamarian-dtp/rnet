package rnet

import (
	"errors"
	"fmt"
	"sync"
)

func New () (*Network) {
	return &Network {sync.Map {}}
}

type Network struct {
	allocations sync.Map /* This map keeps track of net addresses that are in
		use, and the interfaces they are assigned to.
		key: net-addr; val: *interface */
}

// NewIntf () helps create a new network interface.
//
// Inputs
//
// input 0: A desired network address for the interface.
//
// Outpts
//
// outpt 0: The network interface created. If an error ahould occur during
// the creation the interface, value would be nil.
//
// outpt 1: If interface creation should fail, possible values include: NetErrInUse.
func (n *Network) NewIntf (netAddr string) (*Interface, error) {
	// Input data validation.
	if netAddr == "" {
		return nil, errors.New ("Network address can not be an empty string.")
	}
	// Checking if address is already in use. { ...
	_, ok := n.allocations.Load (netAddr)
	if ok == true {
		return nil, NetErrInUse
	}
	// ...}
	// Creating a new network interface. { ...
	i, errX := newIntf (netAddr, n)
	if errX != nil {
		errMssg := fmt.Sprintf ("Unable to create new interface. [%s]",
			errX.Error ())
		return nil, errors.New (errMssg)
	}
	n.allocations.Store (netAddr, i)
	// ... }
	return i, nil
}

// Disconnect () disconnects a network interface from the network.
//
// Inputs
//
// input 0: The network address of the interface to disconnect.
//
// Outpts
//
// outpt 0: If the interface could not be successfully disconnected, value would not
// be nil. Also, if no interface is using the network address provided as input 0,
// NetErrNotInUse would be returned.
func (n *Network) Disconnect (netAddr string) (error) {
	// Trying to get the interface the net address is allocated to. { ...
	alloc, okX := n.allocations.Load (netAddr)
	if okX == false {
		return NetErrNotInUse
	}
	// ... }
	// Disconnecting the network interface. { ...
	n.allocations.Delete (netAddr)
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
	// ... }
	return nil
}

/* provideMDInfo () provides an MDI that could be used to send messages to another
	interface.

	Inputs
	input 0: The net addr of the interface whose MDI should be provided.

	Outpts
	outpt 0: An MDI that could be used to message the interface specified
		as input 0.
	outpt 1: On success, value would be nil. On failure value would be an error.
		If the net addr provided as input 0 is not in use, value would be
		NetErrNotInUse.
*/
func (n *Network) provideMDInfo (netAddr string) (*mDInfo, error) {
	// Fetching the interface whose net addr was provided. { ...
	alloc, ok := n.allocations.Load (netAddr)
	if ok == false {
		return nil, NetErrNotInUse
	}
	// ... }
	// Requesting an MDI from the interface. { ...
	intf, _ := alloc.(*Interface)
	di, errX := intf.getMDInfo ()
	if errX == IntErrNotConnected {
		return nil, NetErrNotInUse
	} else if errX != nil {
		errMssg := fmt.Sprintf ("Unable to get message-delivery info " +
		 	"from recipient. [%s]", errX.Error ())
		return nil, errors.New (errMssg)
	}
	// ... }
	return di, nil
}

var (
	NetErrInUse error = errors.New ("Network address already in use.")
	NetErrNotInUse error = errors.New ("Network address not in use.")
)
