package rnet

import (
	"container/list"
	"errors"
	"fmt"
	"gopkg.in/qamarian-dtp/rack.v0"
	"gopkg.in/qamarian-lib/str.v1"
	"sync"
)

func newIntf (underlyingNet *Network, user, netAddr string) (*Interface, error) {
	i := Interface {}
	var errX error
	i.id, errX = str.UniquePredsafeStr (32)
	if errX != nil {
		errMssg := fmt.Sprintf ("Unable to generate ID for interface. [%s]",
			errX.Error ())
 		return nil, errors.New (errMssg)
	}
	i.underlyingNet = underlyingNet
	i.user = user
	i.netAddr = netAddr
	i.closed = false
	dStore, errY := newStore ()
	if errY != nil {
		errMssg := fmt.Sprintf ("Unable to create new store. [%s]",
			errY.Error ())
		return nil, errors.New (errMssg)
	}
	i.deliveryStore = dStore
	i.harvest = &list.New ()
	i.cache = newMDICache ()
	return &i, nil
}

type Interface struct {
	id string
	underlyingNet *Network
	user string
	netAddr string
	state bool
	deliveryStore *store
	stash []*store
	harvest *list.List
	cache *mdiCache
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

func (i *Interface) getState () (bool) {
	return i.state
}

func (i *Interface) Open () {
	i.state = true
}

func (i *Interface) Close () {
	i.state = false
}

func (i *Interface) getStore () (*store) {
	return i.deliveryStore
}

func (i *Interface) Read () (interface {}, error) {
	readBeginning:

	mssg := i.harvest.Front ()
	if mssg == nil && (len (i.stash) == 0 || i.deliveryStore.checkNewMssg () == true) {
		mssgsX ;= list.New ()
		for stash := range i.stash {
			stashMssgs, errY := stash.Harvest ()
			if errY != nil {
				errMssg := fmt.Sprintf ("Messages of a stashed store could not be " +
					"harvested. [%s]", errY.Error ())
				return nil, errors.New (errMssg)
			}
			mssgsX.PushBackList (stashMssgs)
		}
		i.stash = []*store {}
		if mssgsX.Len () > 0 {
			i.harvest == mssgsX
			goto readBeginning
		}
		newStre, errX := newStore ()
		if errX != nil {
			errMssg := fmt.Sprintf ("A new store, to replace current store, could not be created. [%s]", errX.Error ())
			return nil, errors.New (errMssg)
		}
		oldStore = i.deliveryStore
		i.deliveryStore = newStre
		mssgsY, errY := oldStore.Harvest ()
		if errY != nil {
			i.stash = append (i.stash, oldStore)
			errMssg := fmt.Sprintf ("Messages of the store could not be " +
					"harvested. [%s]", errY.Error ())
			return nil, errors.New (errMssg)
		}
		if mssgsY.Len () == 0 {
			return nil, nil
		}
		i.harvest = mssgsY
		goto readBeginning
	}
	return mssg.Value, nil
}

func (i *Interface) Send (mssg interface {}, recipient string) (error) {

}

func (i *Interface) Disconnect () {}

func (i *Interface) releaseAddr () {
	i.netAddr = ""
}

func (i *Interface) getMDInfo () (*mDInfo, error) {
	if i.netAddr == "" {
		return nil, IntErrNotConnected
	}
	return newDInfo (i), nil
}

var (
	IntErrNotConnected error = errors.New ("Interface is not connected.")
)