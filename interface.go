package rnet

import (
	"container/list"
	"errors"
	"fmt"
	"runtime"
	"sync/atomic"
)

func newIntf (netAddr string, underlyingNet *Network) (*Interface, error) {
	if netAddr == "" || underlyingNet == nil {
		return nil, errors.New ("One or more of the inputs are invalid.")
	}
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
	return &i, nil
}

type Interface struct {
	netAddr string
	underlyingNet *Network
	netState int32
	deliveryStore *store
	stash []*store
	harvest *list.List
	cache *mDICache
}

func (i *Interface) NetAddr () (string) {
	return i.netAddr
}

func (i *Interface) getUNet () (*Network) {
	return i.underlyingNet
}

func (i *Interface) NetState () (int32) {
	return i.netState
}

func (i *Interface) getStore () (*store) {
	return i.deliveryStore
}

func (i *Interface) Send (mssg interface {}, recipient string) (error) {
	if recipient == "" {
		return errors.New ("No recipient network address was specified.")
	}
	lockingBeginning:
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
	errE := mdi.sendMssg (mssg)
	if errE != nil {
		errMssg := fmt.Sprintf ("Unable to send message. [%s]", errE.Error ())
		return errors.New (errMssg)
	}
	return nil
}

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

func (i *Interface) _harvest_ (replaceStore bool) (error) {
	if i.deliveryStore == nil {
		return IntErrNoStoreAvail
	}
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
	mssgsY, errY := oldStore.Harvest ()
	if errY != nil {
		i.stash = append (i.stash, oldStore)
		errMssg := fmt.Sprintf ("Messages of the current store could not be " +
			"harvested. [%s]", errY.Error ())
		return errors.New (errMssg)
	}
	i.harvest.PushBackList (mssgsX)
	i.harvest.PushBackList (mssgsY)
	i.stash = []*store {}
	return nil
}

func (i *Interface) destroy () (error) {
	errX := i._harvest_ (false)
	if errX != nil && errX == IntErrNoStoreAvail {
		errMssg := fmt.Sprintf ("Store could not be harvested. [%s]",
			errX.Error ())
		return errors.New (errMssg)
	}
	changeoverBeginning:
	okX := atomic.CompareAndSwapInt32 (&i.netState, IntStateIdle, IntStateDestroyed)
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
	return nil
}

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
