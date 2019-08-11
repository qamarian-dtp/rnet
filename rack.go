package rnet

import (
	"container/list"
	"errors"
	"fmt"
	"gopkg.in/qamarian-dtp/cart.v1"
)

func newRack () (*rack) {
	rk := &rack {}
	rk.mssgs, rk.mssgsManager = cart.New ()
	return rk
}

type rack struct {
	mssgs *cart.Cart
	mssgsManager *cart.AdminPanel
}

func (r *rack) addMssg (mssg interface {}) (error) {
	errX := r.mssgs.Put (mssg)
	if errX == cart.ErrBeenHarvested {
		return rckErrBeenHarvested
	} else if errX != nil {
		errMssg := fmt.Sprintf ("Unable to add message to rack. [%s]",
			errX.Error ())
		return errors.New (errMssg)
	}
	return nil
}

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
