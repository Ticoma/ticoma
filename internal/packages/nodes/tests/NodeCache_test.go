package tests

// import (
// 	"fmt"
// 	"testing"
// 	. "ticoma/packages/nodes/interfaces"
// 	. "ticoma/packages/nodes/modules"
// 	. "ticoma/packages/nodes/modules/verifier"

// 	assert "github.com/stretchr/testify/assert"
// )

// func TestNodeCacheInit(t *testing.T) {

// 	// Initialize a cache and put first pkg
// 	v := NewVerifier()
// 	nc := NewNodeCache(v)

// 	pkg := ActionDataPackageTimestamped{
// 		ActionDataPackage: &ActionDataPackage{
// 			PlayerId:     1,
// 			PubKey:       "PUBKEY",
// 			Position:     &Position{X: 1, Y: 1},
// 			DestPosition: &DestPosition{X: 2, Y: 2},
// 		},
// 		Timestamp: 1337,
// 	}

// 	// should be empty -> map[]
// 	fmt.Println(nc.GetAll())

// 	// put first pkg in
// 	nc.Put(pkg)

// 	prev := nc.GetPrevious(1)
// 	curr := nc.GetCurrent(1)

// 	assert.Equal(t, prev, pkg)
// 	assert.Equal(t, curr, pkg)

// }

// func TestNodeCacheGetters(t *testing.T) {

// 	// Initialize cache, verifier for player 0
// 	v := NewVerifier()
// 	nc := NewNodeCache(v)

// 	v2 := NewVerifier()
// 	nc2 := NewNodeCache(v2)

// 	// init player 0
// 	pkg := ActionDataPackageTimestamped{
// 		ActionDataPackage: &ActionDataPackage{
// 			PlayerId:     0,
// 			PubKey:       "PUBKEY",
// 			Position:     &Position{X: 1, Y: 1},
// 			DestPosition: &DestPosition{X: 1, Y: 1},
// 		},
// 		Timestamp: 1,
// 	}

// 	nc.Put(pkg)

// 	// Init player 1
// 	pkg2 := ActionDataPackageTimestamped{
// 		ActionDataPackage: &ActionDataPackage{
// 			PlayerId:     1,
// 			PubKey:       "PUBKEY",
// 			Position:     &Position{X: 1, Y: 1},
// 			DestPosition: &DestPosition{X: 1, Y: 1},
// 		},
// 		Timestamp: 223,
// 	}

// 	// Imitate evt listener (temp)
// 	nc.Put(pkg2)
// 	nc2.Put(pkg2)

// 	// Get player1 pos from p0's cache
// 	p0c := nc.GetCache(1)

// 	assert.Equal(t, p0c[0], pkg2)

// 	// Assert self
// 	p0self := nc.GetCurrent(0)

// 	assert.Equal(t, p0self, pkg)

// }

// func TestNodeCachePut(t *testing.T) {

// 	// Init two players
// 	v := NewVerifier()
// 	nc := NewNodeCache(v)

// 	// init player 0
// 	pkg := ActionDataPackageTimestamped{
// 		ActionDataPackage: &ActionDataPackage{
// 			PlayerId:     0,
// 			PubKey:       "PUBKEY",
// 			Position:     &Position{X: 1, Y: 1},
// 			DestPosition: &DestPosition{X: 1, Y: 1},
// 		},
// 		Timestamp: 1,
// 	}

// 	nc.Put(pkg)

// 	// init pkg should fill both prev, curr of p0's node cache
// 	cache := nc.GetCache(0)
// 	prev := cache[0]
// 	curr := cache[1]

// 	assert.Equal(t, prev, pkg)
// 	assert.Equal(t, curr, pkg)

// 	// init second node later
// 	v2 := NewVerifier()
// 	nc2 := NewNodeCache(v2)

// 	cache2 := nc2.GetAll()

// 	// fmt.Println(cache)
// 	assert.Equal(t, len(cache2), 0) // should return empty map

// 	// Player 0 move
// 	pkg2 := ActionDataPackageTimestamped{
// 		ActionDataPackage: &ActionDataPackage{
// 			PlayerId:     0,
// 			PubKey:       "PUBKEY",
// 			Position:     &Position{X: 1, Y: 1},
// 			DestPosition: &DestPosition{X: 3, Y: 3},
// 		},
// 		Timestamp: 223,
// 	}

// 	nc.Put(pkg2)
// 	nc2.Put(pkg2) // imitate evt listener

// 	p0prev := nc.GetPrevious(0)
// 	p0prev1 := nc2.GetPrevious(0)

// 	assert.Equal(t, p0prev, pkg)
// 	assert.Equal(t, p0prev1, pkg2) // p1's node cache should have no record of the first pkg

// }

// func TestNodeCacheInvalidPutRequest(t *testing.T) {

// 	// init p0 cache, p1 cache
// 	v := NewVerifier()
// 	nc := NewNodeCache(v)

// 	v2 := NewVerifier()
// 	nc2 := NewNodeCache(v2)

// 	// init player 0
// 	pkg := ActionDataPackageTimestamped{
// 		ActionDataPackage: &ActionDataPackage{
// 			PlayerId:     0,
// 			PubKey:       "PUBKEY",
// 			Position:     &Position{X: 1, Y: 1},
// 			DestPosition: &DestPosition{X: 1, Y: 1},
// 		},
// 		Timestamp: 1,
// 	}

// 	nc.Put(pkg) // init p0 pos for both
// 	nc2.Put(pkg)

// 	// p0 request invalid move (destPos of prev pkg does not equal current Pos)
// 	pkg2 := ActionDataPackageTimestamped{
// 		ActionDataPackage: &ActionDataPackage{
// 			PlayerId:     0,
// 			PubKey:       "PUBKEY",
// 			Position:     &Position{X: 4, Y: 4},
// 			DestPosition: &DestPosition{X: 4, Y: 4},
// 		},
// 		Timestamp: 2,
// 	}

// 	// this should get rejected
// 	nc.Put(pkg2)
// 	nc2.Put(pkg2)

// 	n0cache := nc.GetCache(0)
// 	n1cache := nc.GetCache(0)

// 	// fmt.Println(n0cache)

// 	assert.Equal(t, n0cache[1], pkg) // pkg is still curr, pkg2 didn't pass
// 	assert.Equal(t, n1cache[1], pkg)

// 	// p0 sends a start move pkg {1, 1} to {12, 12} at ts = 3
// 	pkg3 := ActionDataPackageTimestamped{
// 		ActionDataPackage: &ActionDataPackage{
// 			PlayerId:     0,
// 			PubKey:       "PUBKEY",
// 			Position:     &Position{X: 1, Y: 1},
// 			DestPosition: &DestPosition{X: 12, Y: 12},
// 		},
// 		Timestamp: 3,
// 	}

// 	nc.Put(pkg3)
// 	nc2.Put(pkg3)

// 	// check if start move req passed
// 	p0curr := nc.GetCurrent(0)
// 	p0curr1 := nc2.GetCurrent(0)

// 	assert.Equal(t, p0curr, pkg3)
// 	assert.Equal(t, p0curr1, pkg3)

// 	// p0 sends end move pkg which implies crazy move speed based on ts (24 blocks traveled in 1000ms)
// 	pkg4 := ActionDataPackageTimestamped{
// 		ActionDataPackage: &ActionDataPackage{
// 			PlayerId:     0,
// 			PubKey:       "PUBKEY",
// 			Position:     &Position{X: 12, Y: 12},
// 			DestPosition: &DestPosition{X: 12, Y: 12},
// 		},
// 		Timestamp: 1003,
// 	}

// 	nc.Put(pkg4)
// 	nc2.Put(pkg4)

// 	t.Log("CURRENT ", nc.GetCurrent(0))

// 	assert.Equal(t, nc.GetCurrent(0), pkg3)
// 	assert.Equal(t, nc2.GetCurrent(0), pkg3) // both nodes should reject pkg4
// }
