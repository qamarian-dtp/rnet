package rnet

type storeProtected struct {
	underlyingStore *store
	sendersAddr string
	lastKnownStore string
	rack list.List
}

func (s *storeProtected) AddToRack () {}
