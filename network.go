package rnet

func New () (n *Network, p *UserPanel) {
	return n, p
}

type Network struct {}

func (n *Network) Freeze () {}

func (n *Network) Unfreeze () {}

func (n *Network) Reserve () {}

func (n *Network) Release () {}

func (n *Network) Reclaim () {}
