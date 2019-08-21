package rnet

import (
	"container/list"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

/* newPPO () creates a new PPO.

	Inputs
	input 0: The ID that should be used for the PPO.

	input 1: The network centre for the PPO.

	Outpts
	outpt 0: The PPO created. Value would be nil, if creation fails.

	outpt 1: On success, value would be nil. On failure, value would be an error.
*/
func newPPO (id string, netCentre *NetCentre) (*PPO, error) {
	// Validating inputs. { ...
	if id == "" || netCentre == nil {
		return nil, errors.New ("One or more of the inputs are invalid.")
	}
	// ... }
	// Creating the PPO. { ...
	ppo := PPO {
		id:         id,
		netCentre:  netCentre,
		state:      PpoStateIdle,
		stash:      []*store {},
		harvest:    list.New (),
		cache:      newMDICache (),
	}
	locker := &sync.Mutex {}
	ppo.wakeup = struct {
		wakeupChan       *sync.Cond
		wakeupChanLocker *sync.Mutex
		waiting          bool
	} {sync.NewCond (locker), locker, false}
	dStore, errY := newStore ()
	if errY != nil {
		errMssg := fmt.Sprintf ("Unable to create new store. [%s]", errY.Error ())
		return nil, errors.New (errMssg)
	}
	ppo.store = dStore
	// ... }
	return &ppo, nil
}

type PPO struct {
	id string             // The ID of the PPO.
	netCentre *NetCentre  // The network centre of the PPO.
	state int32           /* The state of the PPO. See the variable section for
		possible values of this data. */
	store *store          // The place where new messages are added.
	stash []*store        // Stores that could not be harvested successfully.
	harvest *list.List    // Messages that have been harvested from store.
	cache *mdiCache       // The MDI cache of the PPO.
	wakeup struct {
		wakeupChan       *sync.Cond
		wakeupChanLocker *sync.Mutex
		waiting          bool
	}
}

// ID () outputs the ID of the PPO.
func (ppo *PPO) ID () (string) {
	return ppo.id
}

func (ppo *PPO) getNC () (*NetCentre) {
	return ppo.netCentre
}

// State () provides the state of the PPO. See the variable section for the possible
// states of a PPO.
func (ppo *PPO) State () (int32) {
	return ppo.state
}

func (ppo *PPO) getStore () (*store) {
	return ppo.store
}

// Send () helps send a message to another PPO on the network.
//
// Inputs
//
// input 0: The message to be sent. Value can not be nil.
//
// input 1: The ID of the recipient PPO.
//
// Outpts
//
// outpt 0: Possible error values include: PpoErrNotConnected, PpoErrIdNotInUse, and
// PpoErrRecipientNotHere.
func (ppo *PPO) Send (mssg interface {}, recipient string) (error) {
	// Validating input data. { ...
	if mssg == nil {
		return errors.New ("Nil messages can not be sent over this network.")
	}
	if recipient == "" {
		return errors.New ("Recipient's ID can not be an empty string.")
	}
	// ... }
	lockingBeginning:
	// Tries to lock the PPO. { ...
	okX := atomic.CompareAndSwapInt32 (&ppo.state, PpoStateIdle, PpoStateSendingTo)
	if okX == false {
		switch ppo.state {
			case PpoStateIdle:
				runtime.Gosched ()
				goto lockingBeginning
			case PpoStateSendingTo:
				runtime.Gosched ()
				goto lockingBeginning
			case PpoStateDestroyed:
				return PpoErrNotConnected
			default:
				return errors.New ("PPO is in an invalid state.")
		}
	}
	defer func () {
		ppo.state = PpoStateIdle
	} ()
	// ... }
	// Getting an MDI to communicate with the recipient. { ...
	mdi := ppo.cache.Get (recipient)
	if mdi == nil {
		newMDI, errG := ppo.netCentre.provideMDI (recipient)
		if errG == NcErrIdNotInUse {
			return PpoErrIdNotInUse
		} else if errG != nil {
			errMssg := fmt.Sprintf ("Network centre unable to provide an " +
				"MDI for communicating with the recipient.")
			return errors.New (errMssg)
		}
		ppo.cache.Put (newMDI, recipient)
		mdi = newMDI
	}
	// ... }
	// Sending the message. { ...
		// Checking if recipient is on the network. { ...
		store := mdi.getRecipientPPO ().getStore ()
		if store == nil {
			return PpoErrRecipientNotHere
		}
		// ... }
		addBeginning:
		/* Adding message to the rack of the sender, in the recipient's store.
			{ ... */
		errT := mdi.getSenderRack ().addMssg (mssg)
		if errT == rckErrBeenHarvested { /* If message could not be added because
			the rack has been harvested, a new rack is added to the new store.
			*/

			errU := mdi.newRack ()
			if errU == mdiErrNotConnected {
				return PpoErrRecipientNotHere
			} else if errU != nil {
				errMssg := fmt.Sprintf ("Unable to create new rack in " +
					"recipient's  new store. [%s]", errU.Error ())
				return errors.New (errMssg)
			}
			goto addBeginning
		} else if errT != nil { /* If message could not be added for some other
			reasons, error is returned. */

			errMssg := fmt.Sprintf ("Unable to add message to the " +
				"recipient's store. [%s]", errT.Error ())
			return errors.New (errMssg)
		}
		store.mssgAdded ()
		for mdi.getRecipientPPO ().waiting () == true &&
			mdi.getRecipientPPO ().Check () == true {
			mdi.getRecipientPPO ().signalNewMssg ()
		}
		// ... }
	// ... }
	return nil
}

// Read () checks if there is any new message. If there is a new message, it returns the
// message. Otherwise, it returns nil.
//
// Outpts
//
// outpt 0: If there is a new message, value would be a message. Messages are read on
// first-in first-out basis (FIFO). If there are no new message, value would be nil.
//
// outpt 1: On success, value would be nil. On failure, value would the error that
// occured.
func (ppo *PPO) Read () (interface {}, error) {
	readBeginning:
	mssg := ppo.harvest.Front ()
	if mssg == nil {
		if ppo.store == nil {
			return nil, nil
		} else if len (ppo.stash) > 0 || ppo.store.checkNewMssg () == true {
			errM := ppo._harvest_ (true)
			if errM == nil && ppo.harvest.Len () == 0 {
				return nil, nil
			} else if errM != nil {
				errMssg := fmt.Sprintf ("Unable to harvest store and " +
					"stash. [%s]", errM.Error ())
				return nil, errors.New (errMssg)
			}
			goto readBeginning
		} else {
			return nil, nil
		}
	}
	ppo.harvest.Remove (mssg)
	return mssg.Value, nil
}

// This function is a shared-sub-function, and should not be called called by any function
// or method not belonging to this data type. This function harvests the delivery store
// (and stashes, if applicable), and place the harvested messages in the harvest (basket).
//
//	Inputs
//	input 0: This value specifies whether a new store should be given to the PPO or
//		not, after the harvest. If value is true, a new store would replace the
//		current one. If the value if false, no store would replace the current one.
//
//	Outpts
//	outpt 0: Possible error values include: ppoErrNoStoreAvail.
func (ppo *PPO) _harvest_ (replaceStore bool) (error) {
	// Harvesting the messages in the stash. { ...
	mssgsX := list.New ()
	for _, stash := range ppo.stash {
		stashMssgs, errY := stash.racksManager.Harvest ()
		if errY != nil {
			errMssg := fmt.Sprintf ("Messages of a stashed stores could " +
				"not be harvested. [%s]", errY.Error ())
			return errors.New (errMssg)
		}
		mssgsX.PushBackList (stashMssgs)
	}
	// ... }
	/* If PPO has no cuurent store, whatever have been harvested are placed in the
		harvest basket. { ... */
	if ppo.store == nil {
		ppo.stash = []*store {}
		ppo.harvest.PushBackList (mssgsX)
		return nil
	}
	// ... }
	// Replacing the current store with a new one. { ...
	oldStore := ppo.store
	if replaceStore == true {
		newStre, errX := newStore ()
		if errX != nil {
			errMssg := fmt.Sprintf ("A new store, to replace current " +
				"store, could not be created. [%s]", errX.Error ())
			return errors.New (errMssg)
		}
		ppo.store = newStre
	} else {
		ppo.store = nil
	}
	// ... }
	// Harvesting the current store. { ...
	mssgsY, errY := oldStore.Harvest ()
	if errY != nil {
		ppo.stash = append (ppo.stash, oldStore)
		errMssg := fmt.Sprintf ("Messages of the current store could not be " +
			"harvested. [%s]", errY.Error ())
		return errors.New (errMssg)
	}
	// ... }
	// Putting harvested messages in the harvest (basket). { ...
	ppo.harvest.PushBackList (mssgsX)
	ppo.harvest.PushBackList (mssgsY)
	// ... }
	// Since harvest of the stash succeeded, the stash is emptied. { ...
	ppo.stash = []*store {}
	// ... }
	return nil
}

// Check () simply checks if there is any new message that could be read. If there is,
// true would be returned. Otherwise, false would be returned.
func (ppo *PPO) Check () (bool) {
	if len (ppo.stash) == 0 && ppo.store == nil {
		return false
	}
	for _, stash := range ppo.stash {
		if stash.checkNewMssg () == true {
			return true
		}
		runtime.Gosched ()
	}
	if ppo.store != nil && ppo.store.checkNewMssg () == true {
		return true
	}
	return false
}

// Wait () suspends the goroutine that calls it, until a new message has been received or
// the PPO is disconnected from the network. In other words, it prevents unnecessary
// wastage of CPU cycles. Rather than using a for loop and Check () to wait for new
// message, this method should be used.
func (ppo *PPO) Wait () {
	if ppo.state == PpoStateDestroyed {
		return
	}
	ppo.wakeup.waiting = true
	defer func () {
		ppo.wakeup.waiting = false
	} ()
	if ppo.Check () == true {
		return
	}
	ppo.wakeup.wakeupChanLocker.Lock ()
	defer ppo.wakeup.wakeupChanLocker.Unlock ()
	ppo.wakeup.wakeupChan.Wait ()
}

// waiting () could be used to check the PPO user is waiting for a new message.
func (ppo *PPO) waiting () (bool) {
	return ppo.wakeup.waiting
}

// signalNewMessage () could be used to wakeup the waiting user of the PPO.
func (ppo *PPO) signalNewMssg () {
	ppo.wakeup.wakeupChan.Signal ()
}


// This function destroys the PPO. In other words, it prevents further sending and
// receiving of messages, although reading of messages that have already been received
// would be permitted.
func (ppo *PPO) destroy () (error) {
	// Harvesting the messages in the store. { ...
	errX := ppo._harvest_ (false)
	if errX != nil {
		errMssg := fmt.Sprintf ("PPO's store could not be harvested. [%s]",
			errX.Error ())
		return errors.New (errMssg)
	}
	// ... }
	changeoverBeginning:
	// Changing the state of the PPO. { ...
	okX := atomic.CompareAndSwapInt32 (&ppo.state, PpoStateIdle, PpoStateDestroyed)
	if okX == false {
		switch ppo.state {
			case PpoStateIdle:
				runtime.Gosched ()
				goto changeoverBeginning
			case PpoStateSendingTo:
				runtime.Gosched ()
				goto changeoverBeginning
			case PpoStateDestroyed:
				return nil
			default:
				return errors.New ("PPO is in an invalid state.")
		}
	}
	// ... }
	for ppo.wakeup.waiting == true {
		ppo.wakeup.wakeupChan.Signal ()
	}
	return nil
}

// getMDI () provides an MDI that can be used to send messages to the PPO.
//
//	Outpts
//	outpt 0: On success, value would be an MDI that can be used to send messages to
// 		the PPO. On failure, value would be nil.
//
//	outpt 1: Possible errors include: PpoErrNotConnected.
func (ppo *PPO) getMDI () (*mdi, error) {
	if ppo.state == PpoStateDestroyed {
		return nil, PpoErrNotConnected
	}
	someMDI, errX := newMDI (ppo)
	if errX != nil {
		errMssg := fmt.Sprintf ("Could not create an MD info. [%s]", errX.Error ())
		return nil, errors.New (errMssg)
	}
	return someMDI, nil
}

var (
	PpoStateIdle      int32 = 0
	PpoStateSendingTo int32 = 1
	PpoStateDestroyed int32 = 2

	PpoErrNotConnected     error = errors.New ("PPO is not connected.")
	PpoErrIdNotInUse       error = errors.New ("The recipient ID provided is not " +
		"in use.")
	PpoErrRecipientNotHere error = errors.New ("The recipient is not on the network.")
)
