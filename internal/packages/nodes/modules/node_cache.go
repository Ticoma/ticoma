package modules

import (
	// "fmt"

	"fmt"
	. "ticoma/internal/packages/nodes/interfaces"
	. "ticoma/internal/packages/nodes/modules/verifier"
)

type PrevAndCurrentADPT [2]ActionDataPackageTimestamped // [ADPT, ADPT]
type Store map[int]PrevAndCurrentADPT                   // NodeCache internal memory model => [Prev ADPT, Current ADPT]

type NodeCache struct {
	*CacheStore
	*NodeVerifier
}

type CacheStore struct {
	Store
}

// NodeCache functions

func NewNodeCache(v *NodeVerifier) *NodeCache {
	return &NodeCache{
		CacheStore:   &CacheStore{},
		NodeVerifier: v,
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
//
// (TMP HERE) Pkg verification process:
// 1. verify fieldNames, types of json string pkg
// 2. extract json pkg values
// 3. initialize ADP struct and try assign values
// 4. engine verifies ADP
// 5. If all 4 pass, put ADP in cache
func (nc *NodeCache) Put(pkgBytes []byte) {

	// Verify types of incoming pkg
	valid := nc.NodeVerifier.SecurityVerifier.VerifyADPTypes(pkgBytes)
	if !valid {
		fmt.Println("[NODE CACHE] - Couldn't verify pkg types.")
	}

	// Construct adpt
	pkg, err := nc.NodeVerifier.SecurityVerifier.ConstructADPT(pkgBytes)
	if err != nil {
		fmt.Println("[NODE CACHE] - ADPT Construct err: ", err)
	}

	// init cache map if first pkg
	if len(nc.CacheStore.Store) == 0 {
		nc.CacheStore.Store = make(Store)
	}

	// if there's no cache and a package arrives, it means it's the first package
	if nc.CacheStore.Store[pkg.PlayerId][0] == (ActionDataPackageTimestamped{}) || nc.CacheStore.Store[pkg.PlayerId][1] == (ActionDataPackageTimestamped{}) {
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
