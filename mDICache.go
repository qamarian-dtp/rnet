package rnet

import (
	"sync"
)

func newMDICache () {
	return &mDICache {sync.Mutex {}, make (map[string]*mDInfo}
}

type mDICache struct {
	locker sync.Mutex
	info map[string]*mDInfo
}

func 