package rnet

import (
	"sync"
)

type diCache struct {
	locker sync.RWMutex
	info map[string]*dInfo
}
