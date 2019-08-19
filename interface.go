package rnet

import (
	"container/list"
	"errors"
	"fmt"
	"runtime"
	"sync/atomic"
)

/* newIntf () creates a new interface.

	Inputs
	input 0: The network addrees allocated to the interface, by the network.

	input 1: A pointer of the network which the interface should operate on.

	Outpts
	outpt 0: The network interface created. Value would be nil, if an interface
		could not be succssfully created.

	outpt 1: If an error should occur during the creation of the network
		interface, value would be an error. Otherwise, value would be nil.
*/
func newIntf (netAddr string, underlyingNet *Network) (*Interface, error) {
	// Validating inputs
	if netAddr == "" || underlyingNet == nil {
		return nil, errors.New ("One or more of the inputs are invalid.")
	}
	// Creating the interface. { ...
	i := Interface {}
	i.underlyingNet = underlyingNet
	i.netAddr = netAddr
	i.netState = IntStateIdle
	dStore, errY := newStore ()
	if errY != nil {
		errMssg := fmt.Sprintf ("Unable to create new store. [%s]",
			errY.Error ())
		return nil, errors.New (errMssg)
	}
	i.deliveryStore = dStore
	i.stash = []*store {}
	i.harvest = list.New ()
	i.cache = newMDICache ()
	// ... }
	return &i, nil
}

type Interface struct {
	netAddr string         // The network address assigned to the interface.
	underlyingNet *Network // The network the interface is attached to.
	netState int32         /* The state of the interface. See the variable
		section for possible values of this data */
	deliveryStore *store   // The place where new messages are added.
	stash []*store         /* Delivery stores that could not be harvested
		successfully. */
	harvest *list.List     // Messages that have been harvested from store.
	cache *mDICache        // The MDI cache of the interface.
}

// NetAddr () provides the network address assigned to the interface.
func (i *Interface) NetAddr () (string) {
	return i.netAddr
}

func (i *Interface) getUNet () (*Network) {
	return i.underlyingNet
}

// NetState () provides the state of the interface. See the variable section for the
// possible states of an interface.
func (i *Interface) NetState () (int32) {
	return i.netState
}

func (i *Interface) getStore () (*store) {
	return i.deliveryStore
}

// Send () helps send a message to another user on the network.
//
// Inputs
//
// input 0: The message to be sent.
//
// input 1: The network address of the recipient.
//
// Outpts
//
// outpt 0: If an error should occur, value of this would be an error. Possible
// values include: IntErrNotConnected, and IntErrAddrNotInUse.
func (i *Interface) Send (mssg interface {}, recipient string) (error) {
	// Validating input data.
	if recipient == "" {
		return errors.New ("No recipient network address was specified.")
	}
	lockingBeginning:
	// Tries to lock the interface. { ...
	okX := atomic.CompareAndSwapInt32 (&i.netState, IntStateIdle,
		IntStateSendingTo)
	if okX == false {
		switch i.netState {
			case IntStateIdle:
				runtime.Gosched ()
				goto lockingBeginning
			case IntStateSendingTo:
				runtime.Gosched ()
				goto lockingBeginning
			case IntStateDestroyed:
				return IntErrNotConnected
			default:
				return errors.New ("Interface is in an invalid state.")
		}
	}
	defer func () {
		i.netState = IntStateIdle
	} ()
	// ... }
	// Getting the MDI of the recipient. { ...
	mdi := i.cache.Get (recipient)
	if mdi == nil {
		var errG error
		mdi, errG = i.underlyingNet.provideMDInfo (recipient)
		if errG == NetErrNotInUse {
			return IntErrAddrNotInUse
		} else if errG != nil {
			errMssg := fmt.Sprintf ("Unable to retrieve a MDI for the " +
				"recipient network address.")
			return errors.New (errMssg)
		}
		i.cache.Put (mdi, recipient)
	}
	// ... }
	// Sending the message. { ...
	errE := mdi.sendMssg (mssg)
	if errE != nil {
		errMssg := fmt.Sprintf ("Unable to send message. [%s]", errE.Error ())
		return errors.New (errMssg)
	}
	// ... }
	return nil
}

// Read () checks if there is any new message. If there is a new message, it returns
// the message. Otherwise, it returns nil.
//
// Outpts
//
// outpt 0: If there is a new message, value would be a message. Messages are read on
// first-come first-served basis. If there are no new message, value would be nil.
//
// outpt 1: If an error was encoutered during operation, value would be an error. A
// possible value of this is IntErrNoStoreAvail. If no error was encountered during
// the operation, value would be nil.
func (i *Interface) Read () (interface {}, error) {
	readBeginning:
	mssg := i.harvest.Front ()
	if mssg == nil {
		if (len (i.stash) > 0 || i.deliveryStore.checkNewMssg () == true) {
			errM := i._harvest_ (true)
			if ((errM == nil) && (i.harvest.Len () == 0)) ||
				errM == IntErrNoStoreAvail {
				return nil, nil
			} else if errM != nil {
				errMssg := fmt.Sprintf ("Unable to harvest store. " +
					"[%s]", errM.Error ())
				return nil, errors.New (errMssg)
			}
			goto readBeginning
		} else {
			return nil, nil
		}
	}
	i.harvest.Remove (mssg)
	return mssg.Value, nil
}

/* This function is a shared-sub-function, and should not be called called by any
	function or method not belonging to this data type. This function harvests
	the delivery store (and stashes, if applicable), the place the harvested
	messages in the harvest (basket).

	Inputs
	input 0: This value specifies whether a new store should be given to the
		interface or not. If value is true, a new store would replace the
		current one. If the value if false, no store would replace the
		current one.

	Outpts
	outpt 0: If an error ahould occur during the harvest, value would be an
		error. A possible error is IntErrNoStoreAvail. If no error should
		occur during the harvest, value would be nil.
*/
func (i *Interface) _harvest_ (replaceStore bool) (error) {
	// If interface has no delivery store, operation halts.
	if i.deliveryStore == nil {
		return IntErrNoStoreAvail
	}
	// Harvesting the stores in the stash. { ...
	mssgsX := list.New ()
	for _, stash := range i.stash {
		stashMssgs, errY := stash.racksManager.Harvest ()
		if errY != nil {
			errMssg := fmt.Sprintf ("Messages of a stashed store could" + 
				"not be harvested. [%s]", errY.Error ())
			return errors.New (errMssg)
		}
		mssgsX.PushBackList (stashMssgs)
	}
	// ... }
	// Replacing the current store with a new one. { ...
	oldStore := i.deliveryStore
	if replaceStore == true {
		newStre, errX := newStore ()
		if errX != nil {
			errMssg := fmt.Sprintf ("A new store, to replace current " +
				"store, could not be created. [%s]", errX.Error ())
			return errors.New (errMssg)
		}
		i.deliveryStore = newStre
	} else {
		i.deliveryStore = nil
	}
	// ... }
	// Harvesting the current store. { ...
	mssgsY, errY := oldStore.Harvest ()
	if errY != nil {
		i.stash = append (i.stash, oldStore)
		errMssg := fmt.Sprintf ("Messages of the current store could not be " +
			"harvested. [%s]", errY.Error ())
		return errors.New (errMssg)
	}
	// ... }
	// Putting harvested messages in the harvest (basket). { ...
	i.harvest.PushBackList (mssgsX)
	i.harvest.PushBackList (mssgsY)
	// ... }
	/* Since harvest of the stash succeeded and the current store was not added
		to the stash, the stash is emptied. { ... */
	i.stash = []*store {}
	// ... }
	return nil
}

// NewMssg () simply checks if there is any message that could be read. If there is,
// true would be returned. Otherwise, false would be returned.
func (i *Interface) NewMssg () (bool) {
	if i.deliveryStore == nil {
		return false
	}
	if i.deliveryStore.checkNewMssg () == true {
		return true
	}
	for _, stash := range i.stash {
		if stash.checkNewMssg () == true {
			return true
		}
	}
	return false
}

/* This function destroys the interface. In other words, it prevents further sending
	and receiving of messages, although reading of messages that are in the
	harvest would not be prohibited. */
func (i *Interface) destroy () (error) {
	// Harvesting the messages in the store. { ...
	errX := i._harvest_ (false)
	if errX != nil && errX == IntErrNoStoreAvail {
		errMssg := fmt.Sprintf ("Store could not be harvested. [%s]",
			errX.Error ())
		return errors.New (errMssg)
	}
	// ... }
	changeoverBeginning:
	// Changing the state of the interface. { ...
	okX := atomic.CompareAndSwapInt32 (&i.netState, IntStateIdle,
		IntStateDestroyed)
	if okX == false {
		switch i.netState {
			case IntStateIdle:
				runtime.Gosched ()
				goto changeoverBeginning
			case IntStateSendingTo:
				runtime.Gosched ()
				goto changeoverBeginning
			case IntStateDestroyed:
				return nil
			default:
				return errors.New ("Interface is in an invalid state.")
		}
	}
	// ... }
	return nil
}

/* This method provides an MDI that can be used to send messages to the interface.

	Outpts
	outpt 0: On success, value would be an MDI that can be used to send messages
		to the interface upon which this method is called upon. On failure,
		value would be nil.

	outpt 1: On success, value would be an error. On failure, value would be an
		error. A possible error is IntErrNotConnected.
*/
func (i *Interface) getMDInfo () (*mDInfo, error) {
	if i.netState == IntStateDestroyed {
		return nil, IntErrNotConnected
	}
	mdi, errX := newMDInfo (i)
	if errX != nil {
		errMssg := fmt.Sprintf ("Could not create an MD info. [%s]",
			errX.Error ())
		return nil, errors.New (errMssg)
	}
	return mdi, nil
}

var (
	IntStateIdle      int32 = 0
	IntStateSendingTo int32 = 1
	IntStateDestroyed int32 = 2

	IntErrNoStoreAvail error = errors.New ("This interface has no store.")
	IntErrNotConnected error = errors.New ("Interface is not connected.")
	IntErrAddrNotInUse error = errors.New ("The recipient network address " +
		"provided is not in use.")
)
