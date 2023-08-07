package tests

import (
	"fmt"
	"testing"
	. "ticoma/internal/packages/nodes/interfaces"
	nodecache "ticoma/internal/packages/nodes/modules"
	"ticoma/internal/packages/nodes/utils"

	assert "github.com/stretchr/testify/assert"
)

func TestNodeCacheInit(t *testing.T) {

	// Initialize a cache and put first pkg
	nc := nodecache.New()

	pkg := ActionDataPackage{
		PlayerId:     1,
		PubKey:       "PUBKEY",
		Position:     &Position{X: 1, Y: 1},
		DestPosition: &DestPosition{X: 2, Y: 2},
	}

	// should be empty -> map[]
	fmt.Println(nc.GetAll())

	// put first pkg in
	pkgStr := utils.StringifyPkg(pkg, false)
	nc.Put([]byte(pkgStr))

	prev := nc.GetPrevious(1)
	curr := nc.GetCurrent(1)

	// Remove timestamp from cache
	prevADP := utils.StripPkgFromTimestamp(&prev)
	currADP := utils.StripPkgFromTimestamp(&curr)

	assert.Equal(t, prevADP, &pkg)
	assert.Equal(t, currADP, &pkg)

}

func TestNodeCacheGetters(t *testing.T) {

	// Initialize cache, verifier for player 0
	nc := nodecache.New()
	nc2 := nodecache.New()

	// init player 0
	pkg := ActionDataPackage{
		PlayerId:     0,
		PubKey:       "PUBKEY0",
		Position:     &Position{X: 1, Y: 1},
		DestPosition: &DestPosition{X: 1, Y: 1},
	}

	pkgStr := utils.StringifyPkg(pkg, false)
	nc.Put([]byte(pkgStr))

	// Init player 1
	pkg2 := ActionDataPackage{
		PlayerId:     1,
		PubKey:       "PUBKEY1",
		Position:     &Position{X: 2, Y: 2},
		DestPosition: &DestPosition{X: 2, Y: 2},
	}

	pkgStr2 := utils.StringifyPkg(pkg2, false)

	// Imitate PlayerNode pubsub listener
	nc.Put([]byte(pkgStr2))
	nc2.Put([]byte(pkgStr2))

	// Get player1 pos from p0's cache
	p0c := nc.GetCache(1)
	p0cPrev, p0cCurr := p0c[0], p0c[1]
	p0cPrevADP := utils.StripPkgFromTimestamp(&p0cPrev)
	p0cCurrADP := utils.StripPkgFromTimestamp(&p0cCurr)

	// Both prev, curr of p0's cache for p1 should be p1's init pkg
	assert.Equal(t, p0cPrevADP, &pkg2)
	assert.Equal(t, p0cCurrADP, &pkg2)

	// Assert self
	p0self := nc.GetCurrent(0)
	p0selfADP := utils.StripPkgFromTimestamp(&p0self)
	assert.Equal(t, p0selfADP, &pkg)

}

func TestNodeCachePut(t *testing.T) {

	// Init two players
	nc := nodecache.New()

	// init player 0
	pkg := ActionDataPackage{
		PlayerId:     0,
		PubKey:       "PUBKEY0",
		Position:     &Position{X: 1, Y: 1},
		DestPosition: &DestPosition{X: 1, Y: 1},
	}

	pkgStr := utils.StringifyPkg(pkg, false)
	nc.Put([]byte(pkgStr))

	cache := nc.GetCache(0)
	prev := cache[0]
	curr := cache[1]
	prevADP := utils.StripPkgFromTimestamp(&prev)
	currADP := utils.StripPkgFromTimestamp(&curr)

	// init pkg should fill both prev, curr of p0's node cache
	assert.Equal(t, prevADP, &pkg)
	assert.Equal(t, currADP, &pkg)

	// init second node later
	nc2 := nodecache.New()

	cache2 := nc2.GetAll()

	// fmt.Println(cache)
	assert.Equal(t, len(cache2), 0) // should return empty map

	// Player 0 move
	pkg2 := ActionDataPackage{
		PlayerId:     0,
		PubKey:       "PUBKEY0",
		Position:     &Position{X: 1, Y: 1},
		DestPosition: &DestPosition{X: -3, Y: -3},
	}
	pkg2Str := utils.StringifyPkg(pkg2, false)

	nc.Put([]byte(pkg2Str))
	nc2.Put([]byte(pkg2Str))

	p0prev := nc.GetPrevious(0)
	p0prev2 := nc2.GetPrevious(0) // get p0's previous package from NodeCache2
	p0prevADP := utils.StripPkgFromTimestamp(&p0prev)
	p0prev2ADP := utils.StripPkgFromTimestamp(&p0prev2)

	assert.Equal(t, p0prevADP, &pkg)
	assert.Equal(t, p0prev2ADP, &pkg2) // NodeCache2 should have no record of p0's init pkg

}

func TestNodeCacheInvalidPutRequest(t *testing.T) {

	// init Nc, Nc2
	nc := nodecache.New()
	nc2 := nodecache.New()

	// init player 0
	pkg0 := ActionDataPackage{
		PlayerId:     0,
		PubKey:       "PUBKEY0",
		Position:     &Position{X: 1, Y: 1},
		DestPosition: &DestPosition{X: 1, Y: 1},
	}
	pkg0Str := utils.StringifyPkg(pkg0, false)

	// init p0 package for both Nc, Nc2
	nc.Put([]byte(pkg0Str))
	nc2.Put([]byte(pkg0Str))

	// p0 request invalid move (destPos of prev pkg does not equal current Pos)
	pkg0invalid := ActionDataPackage{
		PlayerId:     0,
		PubKey:       "PUBKEY0",
		Position:     &Position{X: 4, Y: 4},
		DestPosition: &DestPosition{X: 4, Y: 4},
	}
	pkg0invalidStr := utils.StringifyPkg(pkg0invalid, false)

	// this should get rejected
	nc.Put([]byte(pkg0invalidStr))
	nc2.Put([]byte(pkg0invalidStr))

	NcCache := nc.GetCache(0)
	Nc2Cache := nc2.GetCache(0)
	NcCurrADP := utils.StripPkgFromTimestamp(&NcCache[1])
	Nc2CurrADP := utils.StripPkgFromTimestamp(&Nc2Cache[1])

	// fmt.Println(NcCurrADP, Nc2CurrADP) // DEBUG

	assert.Equal(t, NcCurrADP, &pkg0) // pkg is still curr, pkg2 didn't pass
	assert.Equal(t, Nc2CurrADP, &pkg0)

	// p0 sends a start move pkg P{1, 1} to dP{12, 12}
	pkg0move := ActionDataPackage{
		PlayerId:     0,
		PubKey:       "PUBKEY0",
		Position:     &Position{X: 1, Y: 1},
		DestPosition: &DestPosition{X: 12, Y: 12},
	}
	pkg0moveStr := utils.StringifyPkg(pkg0move, false)

	// Send pkg
	nc.Put([]byte(pkg0moveStr))
	nc2.Put([]byte(pkg0moveStr))

	Ncp0curr := nc.GetCurrent(0)
	Nc2p0curr := nc2.GetCurrent(0)
	NcPlayer0CurrADP := utils.StripPkgFromTimestamp(&Ncp0curr)
	Nc2Player0CurrADP := utils.StripPkgFromTimestamp(&Nc2p0curr)

	// check if move request "pkg0move" got validated
	assert.Equal(t, NcPlayer0CurrADP, &pkg0move)
	assert.Equal(t, Nc2Player0CurrADP, &pkg0move)

	// p0 sends end move pkg with a direction shift (prev destPos =/= currPos), tp/backtrack attempt
	// This pkg below should get rejected
	cheatPkg := ActionDataPackage{
		PlayerId:     0,
		PubKey:       "PUBKEY0",
		Position:     &Position{X: 0, Y: 0},
		DestPosition: &DestPosition{X: 0, Y: 0},
	}
	cheatPkgStr := utils.StringifyPkg(cheatPkg, false)

	// Send it
	nc.Put([]byte(cheatPkgStr))
	nc2.Put([]byte(cheatPkgStr))

	// Check cache
	NcCurr := nc.GetCurrent(0)
	Nc2Curr := nc.GetCurrent(0)
	NcCurrADP2 := utils.StripPkgFromTimestamp(&NcCurr)
	Nc2CurrADP2 := utils.StripPkgFromTimestamp(&Nc2Curr)

	// t.Log("CURRENT ", nc.GetCurrent(0)) // DEBUG

	// Both Nc should reject cheatPkg and still hold last verified pkg - pkg0move
	assert.Equal(t, NcCurrADP2, &pkg0move)
	assert.Equal(t, Nc2CurrADP2, &pkg0move)
}
