package rnet

// newMDICache () creates a new MDI cache.
func newMDICache () (*mdiCache) {
	return &mdiCache {make (map[string]*mdi)}
}

type mdiCache struct {
	mdi map[string]*mdi  /*
		The MDIs in the cache.
		key: id of recipient PPO; value: MDI of recipient PPO
		*/
}

// Put () puts an MDI in the cache.
func (c *mdiCache) Put (mdi *mdi, ppoID string) {
	if ppoID == "" {
		return
	}
	c.mdi[ppoID] = mdi
}

// Get () gets an MDI from the cache.
func (c *mdiCache) Get (ppoID string) (*mdi) {
	return c.mdi[ppoID]
}
