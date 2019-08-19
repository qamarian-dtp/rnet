package rnet

import (
	"container/list"
	"errors"
	"fmt"
	"gopkg.in/qamarian-dtp/cart.v1"
)

// newRack () helps create a new rack.
func newRack () (*rack) {
	rk := &rack {}
	rk.mssgs, rk.mssgsManager = cart.New ()
	return rk
}

type rack struct {
	mssgs *cart.Cart               // The messages in the rack.
	mssgsManager *cart.AdminPanel  // A data you could use to harvest the messages in the rack.
}

// addMssg () adds a message to a rack.
func (r *rack) addMssg (mssg interface {}) (error) {
	errX := r.mssgs.Put (mssg)
	if errX == cart.ErrBeenHarvested {
		return rckErrBeenHarvested
	} else if errX != nil {
		errMssg := fmt.Sprintf ("Unable to add message to rack. [%s]", errX.Error ())
		return errors.New (errMssg)
	}
	return nil
}

// harvest () helps harvest the messages in the rack.
func (r *rack) harvest () (*list.List, error) {
	mssgs, errX := r.mssgsManager.Harvest ()
	if errX == cart.ErrBeenHarvested {
		return nil, rckErrBeenHarvested
	} else if errX != nil {
		errMssg := fmt.Sprintf ("Unable to harvest rack. [%s]", errX.Error ())
		return nil, errors.New (errMssg)
	}
	return mssgs, nil
}

var (
	rckErrBeenHarvested error = errors.New ("This rack has already been harvested.")
)
