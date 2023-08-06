package modules

import (
	// "fmt"
	"encoding/json"
	"fmt"
	. "ticoma/packages/nodes/interfaces"
	. "ticoma/packages/nodes/modules/verifier"
)

type PrevAndCurrentADPT [2]ActionDataPackageTimestamped // [ADPT, ADPT]
type Store map[int]PrevAndCurrentADPT                   // NodeCache internal memory model => [Prev ADPT, Current ADPT]

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
}

func (nc *NodeCache) GetPrevious(id int) ActionDataPackageTimestamped {
	return nc.CacheStore.Store[id][0]
}

func (nc *NodeCache) GetCurrent(id int) ActionDataPackageTimestamped {
	return nc.CacheStore.Store[id][1]
}

// Put new package to NodeCache
//
// (move stack to the left and delete oldest package from cache)
func (nc *NodeCache) Put(pkg ActionDataPackageTimestamped) {

	// Conv json pkg to byte array
	jsonBytes, err := json.Marshal(pkg)

	if err != nil {
		fmt.Println("[NODE CACHE] - Couldn't serialize package. ", err)
		return
	}

	// Type check
	validPkgTypes := nc.SecurityVerifier.VerifyADPTypes(jsonBytes, true)
	if !validPkgTypes {
		fmt.Println("[NODE CACHE] - Couldn't verify package types.")
		return
	}

	// init cache map if needed
	if len(nc.CacheStore.Store) == 0 {
		fmt.Println("INIT CACHE MAP")
		nc.CacheStore.Store = make(Store)
	}

	// if there's no cache and a package arrives, it means it's the first package
	if nc.CacheStore.Store[pkg.PlayerId][0] == (ActionDataPackageTimestamped{}) || nc.CacheStore.Store[pkg.PlayerId][1] == (ActionDataPackageTimestamped{}) {
		fmt.Println("FIRST PKG")
		var store PrevAndCurrentADPT
		store[0], store[1] = pkg, pkg
		nc.CacheStore.Store[pkg.PlayerId] = store
		return
	}

	validPos := nc.EngineVerifier.VerifyLastMovePos(nc.CacheStore.Store[pkg.PlayerId][1].DestPosition, pkg.Position)

	if !validPos {
		fmt.Println("[NODE CACHE] - Coulnd't verify move direction or position. ", err)
		return
	}

	curr := nc.CacheStore.Store[pkg.PlayerId][1]
	validMove := nc.EngineVerifier.VerifyPlayerMovement(&curr, &pkg)

	if !validMove {
		fmt.Println("[NODE CACHE] - Engine couldn't verify move. ", err)
		return
	}

	// push stack to the left
	cache := PrevAndCurrentADPT{
		curr, pkg,
	}
	nc.CacheStore.Store[pkg.PlayerId] = cache

}
