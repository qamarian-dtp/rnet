package rnet

import (
	"sync"
)

func newMDICache () {
	return &mdiCache {sync.Mutex {}, make (map[string]*dInfo
type diCache struct {
	locker sync.RWMutex
	info map[string]*dInfo
}
