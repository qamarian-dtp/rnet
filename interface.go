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
	i.harvestBasket = &list.List {}
	dStore, errY := newStore ()
	if errY != nil {
		errMssg := fmt.Sprintf ("Unable to create new store. [%s]",
			errY.Error ())
		return nil, errors.New (errMssg)
	}
	i.deliveryStore = dStore
	i.cache = newMDICache ()
	return &i, nil
}

type Interface struct {
	id string
	underlyingNet *Network
	user string
	netAddr string
	closed bool
	harvestBasket *list.List
	deliveryStore *store
	cache *mdiCache
}

func (i *Interface) Open () {
	i.closed = true
}

func (i *Interface) Opened () (bool) {
	return !i.closed
}

func (i *Interface) Send (mssg interface {}, recipient string) (error) {

}

func (i *Interface) Read () (interface {}) {
	readBeginning:

	mssg := i.harvestBasket.Front ()
	if mssg == nil {
		harvest := i.deliveryStore
		newStre, errX := newStore ()
		if errX != nil {
			errMssg = fmt.Sprintf ("Delivery store could not be " +
				"harvested. [%s]", errX.Error ())
			return nil, errors.New (errMssg)
		}
		i.deliveryStore = newStre
		harvest.setState (StrStateToBeHarvested)
}

func (i *Interface) Wait () {}

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
