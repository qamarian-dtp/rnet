package rnet

import (
	"sync"
)

func New () (*Network) {
	// id
	net := Network {id, false, false, sync.Map {},

		struct {
			locked int8
			addr map[string][]string
		},

		struct {
			locked int8
			addr map[string]bool
		},
	}
	return net
}

type Network struct {
	id string
	freezed bool
	locked bool

	allocations sync.Map // KEY: net-addr; VAL: user, net-interface, expirey-signal
	userAddr struct {
		locked int8
		addr map[string][]string
	}
	reservedAddr struct {
		locked int8
		addr map[string]bool
	}
}

func (n *Network) NewIntf () {}

func (n *Network) GetUser () {}

func (n *Network) GetAddr () {}

func (n *Network) GetStoreP () {}

// --------------

func (n *Network) Reserve () {}

func (n *Network) ScanReserver () {}

func (n *Network) Release () {}

func (n *Network) Reclaim () {}

func (n *Network) Lock () {}

func (n *Network) Unlock () {}

func (n *Network) Freeze () {}

func (n *Network) Freezed () {}

func (n *Network) Unfreeze () {}
