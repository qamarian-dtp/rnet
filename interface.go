package rnet

import (
	"container/list"
)

type Interface struct { // Some comment
	interfaceID   string
	networkAddr   string
	addrReclaimed bool
	interfaceOpen bool

	harvestBasket list.List
	deliveryStore struct {
		id    string
		racks list.List
	}

	cache
}

func (i *Interface) Open () {}

func (i *Interface) Opened () (bool) {
	return i.interfaceOpen
}

func (i *Interface) Send () {}

func (i *Interface) Read () {}

func (i *Interface) DLink () {}

func (i *Interface) Close () {}

func (i *Interface) Release () {}

// ------------

func (i *Interface) addrReclaimed ()
