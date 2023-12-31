package cache

import (
	// "fmt"

	"fmt"
	"ticoma/internal/debug"
	"ticoma/internal/packages/gamenode/cache/interfaces"
	"ticoma/internal/packages/gamenode/cache/verifier"
)

type PrevAndCurrentADPT [2]interfaces.ActionDataPackageTimestamped // [ADPT, ADPT]
type Store map[int]PrevAndCurrentADPT                              // NodeCache internal memory model => [Prev ADPT, Current ADPT]

type NodeCache struct {
	*CacheStore
	*verifier.NodeVerifier
}

type CacheStore struct {
	Store
}

// NodeCache functions

func New() *NodeCache {
	v := verifier.New()
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

func (nc *NodeCache) GetPrevious(id int) interfaces.ActionDataPackageTimestamped {
	return nc.CacheStore.Store[id][0]
}

func (nc *NodeCache) GetCurrent(id int) interfaces.ActionDataPackageTimestamped {
	return nc.CacheStore.Store[id][1]
}

// Put new package to NodeCache
//
// (move stack to the left and delete oldest package from cache)
func (nc *NodeCache) Put(pkgBytes []byte) error {

	// Try construct adpt
	pkg, err := nc.NodeVerifier.SecurityVerifier.ConstructADPT(pkgBytes)
	if err != nil {
		return fmt.Errorf("[NODE CACHE] - ADPT Construct err: %w", err)
	}

	// init cache map if first pkg
	if len(nc.CacheStore.Store) == 0 {
		nc.CacheStore.Store = make(Store)
	}

	// if there's no cache and a package arrives, it means it's the first package
	if nc.CacheStore.Store[pkg.PlayerId][0] == (interfaces.ActionDataPackageTimestamped{}) || nc.CacheStore.Store[pkg.PlayerId][1] == (interfaces.ActionDataPackageTimestamped{}) {
		var store PrevAndCurrentADPT
		store[0], store[1] = pkg, pkg
		nc.CacheStore.Store[pkg.PlayerId] = store
		debug.DebugLog("[NODE CACHE] - First pkg from: "+string(rune(pkg.PlayerId)), debug.PLAYER)
		return nil
	}

	validPos := nc.EngineVerifier.VerifyLastMovePos(nc.CacheStore.Store[pkg.PlayerId][1].DestPosition, pkg.Position)
	if !validPos {
		return fmt.Errorf("[NODE CACHE] - Coulnd't verify move direction or position.%s", "")
	}

	curr := nc.CacheStore.Store[pkg.PlayerId][1]
	validMove := nc.EngineVerifier.VerifyPlayerMovement(&curr, &pkg)

	if !validMove {
		return fmt.Errorf("[NODE CACHE] - Engine couldn't verify move.%s", "")
	}

	// push stack to the left
	cache := PrevAndCurrentADPT{
		curr, pkg,
	}
	nc.CacheStore.Store[pkg.PlayerId] = cache
	return nil

}
