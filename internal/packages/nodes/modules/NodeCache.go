package modules

import (
	// "fmt"
	"encoding/json"
	"fmt"
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

func (nc *NodeCache) GetAll() Store {
	return nc.CacheStore.Store
}

func (nc *NodeCache) GetCache(id int) PrevAndCurrentADPT {
	cache := nc.CacheStore.Store[id]
	// fmt.Println("CACHE", cache)
	return cache
	// cache := cs.Store[id]
	// return cache
}

func (nc *NodeCache) GetPrevious(id int) ActionDataPackageTimestamped {
	return nc.CacheStore.Store[id][0]
	// return cs.Store[id][0]
}

func (nc *NodeCache) GetCurrent(id int) ActionDataPackageTimestamped {
	return nc.CacheStore.Store[id][1]
	// return cs.Store[id][1]
}

func (nc *NodeCache) Put(pkg ActionDataPackageTimestamped) {

	jsonBytes, err := json.Marshal(pkg)
	if err != nil {
		fmt.Println("[NODE CACHE] - Couldn't serialize package. ", err)
	}

	verified := nc.SecurityVerifier.VerifyADPTypes(jsonBytes, true)
	if !verified {
		fmt.Println("[NODE CACHE] - Couldn't verify package types. ", err)
	}

	// init cache map if needed
	if len(nc.CacheStore.Store) == 0 {
		fmt.Println("INIT CACHE MAP")
		nc.CacheStore.Store = make(Store)
	}

	// if there's no cache and a package arrives, it means it's the first package
	if nc.CacheStore.Store[pkg.PlayerId][0] == (ActionDataPackageTimestamped{}) || nc.CacheStore.Store[pkg.PlayerId][1] == (ActionDataPackageTimestamped{}) {
		fmt.Println("FIRST PKG")
		// nc.CacheStore.Store = make(Store)

		var test PrevAndCurrentADPT
		test[0], test[1] = pkg, pkg
		nc.CacheStore.Store[pkg.PlayerId] = test
	}

}
