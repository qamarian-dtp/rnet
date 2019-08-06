package rnet

import (
	"container/list"
	"errors"
	"fmt"
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
	closed bool
	deliveryStore *store
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

func (i *Interface) getStore () (*store) {
	return i.deliveryStore
}

func (i *Interface) Open () {
	i.closed = true
}

func (i *Interface) Opened () (bool) {
	return !i.closed
}

func (i *Interface) Send (mssg interface {}, recipient string) (error) {

}

func (i *Interface) Read () (interface {}, error) {
	readBeginning:

	mssg := i.harvest.Front ()
	if mssg == nil {
		harvest := i.deliveryStore
		newStre, errX := newStore ()
		if errX != nil {
			errMssg := fmt.Sprintf ("Delivery store could not be " +
				"harvested. [%s]", errX.Error ())
			return nil, errors.New (errMssg)
		}
		i.deliveryStore = newStre
		okX := harvest.setState (StrStateToBeHarvested)
		if okX == false {
			errMssg := fmt.Sprintf ("Delivery store could not be " +
				"harvested.")
			return nil, errors.New (errMssg)
		}
		mssgs, errY := harvest.harvest ()
		if errY != nil {
			errMssg := fmt.Sprintf ("Delivery store could not be " +
				"harvested. [%s]", errY.Error ())
			return nil, errors.New (errMssg)
		}
		if mssgs.Len () == 0 {
			return nil, nil
		} else {
			i.harvest = mssgs
			goto readBeginning
		}
	}
	return mssg.Value, nil
}

func (i *Interface) Close () {
	i.closed = false
}

func (i *Interface) Disconnect () {}

func (i *Interface) releaseAddr () {
	i.netAddr = ""
}

func (i *Interface) getDInfo () (*dInfo, error) {
	if i.netAddr == "" {
		return nil, IntErrNotConnected
	}
	return newDInfo (i), nil
}

var (
	IntErrNotConnected error = errors.New ("Interface is not connected.")
)
