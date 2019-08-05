type rack struct {
	state int32
	rack list.List
}

func (r *rack) lock () (error) {
	for {
	ok := atomic.CompareAndSwapInt32 (&r.state, 0, 1)
	if ok == false && r.state == 2 {
		return RckErrToBeHarvested
	}
}

func (r *rack) unlock () {
	atomic.CompareAndSwapInt32 (&r.state, 1, 0)
}

var (
	RckErrToBeHarvested error = errors.New ("This rack is about to be " +
		"harvested.")
)

type spCache struct {
	locker sync.RWMutex
	storeP map[string]*storeProtected
}
