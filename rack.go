package rner

import (
	"container/list"
	"errors"
	"fmt"
	"gopkg.in/qamarian-dtp/cart.v0"
)

func newRack () (*rack) {
	rk := &rack {}
	rk.mssgCart, rk.adminPanel = cart.New ()
	return rk
}

type rack struct {
	mssgs *cart.Cart
	mssgsManager *cart.AdminPanel
}

func (r *rack) addMssg (mssg interface {}) (error) {
	errX := r.mssgs.Put (mssg)
	if errX == cart.ErrBeenHarvested {
		return RckErrBeenHarvested
	} else if errX != nil {
		errMssg := fmt.Sprintf ("Unable to add message to rack. [%s]", errX.Error ())
		return errors.New (errMssg)
	}
	return nil

func (r *rack) harvest () (*list.List, error) {
	mssgs, errX := r.mssgsManager.Harvest ()
	if errX == cart.ErrBeenHarvested {
		return nil, RckErrBeenHarvested
	} else if errX != nil {
		errMssg := fmt.Sprintf ("Unable to harvest rack. [%s]", errX.Error ())
		return nil, errors.New (errMssg)
	}
	return mssgs, nil
}

var (
	RckErrBeenHarvested error = errors.New ("This rack has already been harvested.")
)