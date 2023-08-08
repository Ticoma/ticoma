package nodecache

import (
	// "fmt"

	"fmt"
	intf "ticoma/internal/packages/player/interfaces"
	"ticoma/internal/packages/player/nodecache/verifier"
)

type PrevAndCurrentADPT [2]intf.ActionDataPackageTimestamped // [ADPT, ADPT]
type Store map[int]PrevAndCurrentADPT                        // NodeCache internal memory model => [Prev ADPT, Current ADPT]

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

func (nc *NodeCache) GetPrevious(id int) intf.ActionDataPackageTimestamped {
	return nc.CacheStore.Store[id][0]
}

func (nc *NodeCache) GetCurrent(id int) intf.ActionDataPackageTimestamped {
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
	if nc.CacheStore.Store[pkg.PlayerId][0] == (intf.ActionDataPackageTimestamped{}) || nc.CacheStore.Store[pkg.PlayerId][1] == (intf.ActionDataPackageTimestamped{}) {
		var store PrevAndCurrentADPT
		store[0], store[1] = pkg, pkg
		nc.CacheStore.Store[pkg.PlayerId] = store
		return nil
	}

	validPos := nc.EngineVerifier.VerifyLastMovePos(nc.CacheStore.Store[pkg.PlayerId][1].DestPosition, pkg.Position)

	if !validPos {
		return fmt.Errorf("[NODE CACHE] - Coulnd't verify move direction or position. err: %w", err)
	}

	curr := nc.CacheStore.Store[pkg.PlayerId][1]
	validMove := nc.EngineVerifier.VerifyPlayerMovement(&curr, &pkg)

	if !validMove {
		return fmt.Errorf("[NODE CACHE] - Engine couldn't verify move. %w", err)
	}

	// push stack to the left
	cache := PrevAndCurrentADPT{
		curr, pkg,
	}
	nc.CacheStore.Store[pkg.PlayerId] = cache
	return nil

}
