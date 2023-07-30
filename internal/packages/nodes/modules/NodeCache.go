package modules

import (
	// "fmt"
	. "ticoma/packages/nodes/interfaces"
	. "ticoma/packages/nodes/modules/verifier"
)

// NodeCache internal memory model => [Prev, Current]
type PrevAndCurrentADPT [2]ActionDataPackageTimestamped
type Store map[int]PrevAndCurrentADPT

type NodeCache struct {
	*CacheStore
	*Verifier
}

type CacheStore struct {
	Store
}

// NodeCache functions

func NewNodeCache(v *Verifier) *NodeCache {
	return &NodeCache{
		CacheStore: &CacheStore{},
		Verifier:   v,
	}
}

func (cs *CacheStore) GetAll() Store {
	return cs.Store
}

func (cs *CacheStore) GetCache(id int) PrevAndCurrentADPT {
	cache := cs.Store[id]
	return cache
}

func (cs *CacheStore) GetPrevious(id int) ActionDataPackageTimestamped {
	return cs.Store[id][0]
}

func (cs *CacheStore) GetCurrent(id int) ActionDataPackageTimestamped {
	return cs.Store[id][1]
}

func (cs *CacheStore) Put(pkg ActionDataPackageTimestamped) {
	// soon
}
