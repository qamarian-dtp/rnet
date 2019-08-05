type store struct {
	id string
	state int32 /* 0: not in use; 1: about to-be manipulated; 2: about to-be
		harvested */
	racks list.List
	wakeupSignal struct {
		waiting bool
		signalChan sync.Cond
	}
}

func (s *store) getID () (string) {
	return s.id
}

func (s *store) lockRacks () {
	for {
		ok := atomic.CompareAndSwapInt32 (*s.racks.state, 0, 1)
		if ok == true {
			break
		}
	}
}

func (s *store) unlockRacks () {
	s.racks.state = 0
}

func (s *store) addRack (senderRack *rack) (error) {
	if s.racks.state == 2 {
		return StrErrToBeHarvested
	}
	if s.racks.Len () == 0 {
		s.racks.PushFront (senderRack)
	} else {
		s.racks.PushBack (senderRack)
	}
	s.racks.state = 0
}

var (
	StrErrToBeHarvested error = errors.New ("The store is about to be harvested.")
)
