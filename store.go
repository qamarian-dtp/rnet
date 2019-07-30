package rnet

import (
	"container/list"
	"sync"
)

type store struct {
	id string
	racks list.List
	wakeupSignal struct {
		waiting bool
		signalChan sync.Cond
	}
}

func (s *store) AddToRack () {}
