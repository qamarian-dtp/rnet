package rnet

func New () (n *Network, p *UserPanel) {
	return n, p
}

type Network struct {
	freezed bool

	allocatedAddrs map[string]interface {}
	reservedAddrs d
}

func (n *Network) Freeze () {}

func (n *Network) Unfreeze () {}

func (n *Network) Reserve () {}

func (n *Network) Release () {}

func (n *Network) Reclaim () {}

func (n *Network) NewNI () {}

func (n *Network) GetNI (addr string) (*Interface) {} // get network interface of addr

func (n *Network) GetNA () {}

func (n *Network) GetID () {}

func (n *Network) Freezed () {}
