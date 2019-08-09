package rnet

import (
	"sync"
)

func newMDICache () {
	return &mDICache {make (map[string]*mDInfo}
}

type mDICache struct {
	mdi map[string]*mDInfo
}

func (c *mDICache) Get (netAddr string) (*mDInfo) {
	if netAddr == "" {
		return nil
	}
	return c.mdi[netAddr]
}

func (c *mDICache) Put (mdi *mDInfo, netAddr string) {
	if netAddr == "" {
		return
	}
	c.mdi[netAddr] = mdi
}
