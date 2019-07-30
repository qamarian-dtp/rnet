package rnet

import (
	"sync"
)

func New () (n *Network) {
	return
}

type Network struct {
	freezed bool

	allocatedAddrs sync.Map // KEY: net-addr; VAL: delivery-store-rack, reclaimed signal
	reservedAddrs d
}

func (n *Network) NewIntf () {}

func (n *Network) GetUser () {}

func (n *Network) GetAddr () {}

func (n *Network) GetStoreP () {}

// --------------

func (n *Network) Reserve () {}

func (n *Network) Release () {}

func (n *Network) ReleaseAll () {}

func (n *Network) Reclaim () {}

func (n *Network) Lock () {}

func (n *Network) Unlock () {}

func (n *Network) Freeze () {}

func (n *Network) Freezed () {}

func (n *Network) Unfreeze () {}
