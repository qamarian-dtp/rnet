package rnet

import (
	"container/list"
)

type Interface struct { // Some comment
	underlyingNetwork *Network

	interfaceID   string
	networkAddr   string
	interfaceOpen bool

	harvestBasket list.List
	deliveryStore *store

	spCache map[string]struct {
		storeP *storeProtected
		expired *bool
	}
}

func (i *Interface) Open () {}

func (i *Interface) Opened () () {}

func (i *Interface) Send () {}

func (i *Interface) Read () {}

func (i *Interface) Check () {}

func (i *Interface) DLink () {}

func (i *Interface) Close () {}

func (i *Interface) ReleaseAddr () {}

func (i *Interface) NewIntf () {}

func (i *Interface) GetUser () {}

func (i *Interface) GetAddr () {}

func (i *Interface) Destroy () {}
