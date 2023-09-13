package tests

import (
	"fmt"
	"testing"
	"ticoma/internal/pkgs/gamenode/cache"
	"ticoma/types"

	// "ticoma/types"
	assert "github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {

	// Get all on just initialized cache
	c := cache.New()
	all := c.GetAll()
	assert.Equal(t, cache.Memory{}, all)
	assert.Equal(t, cache.PlayerStates{}, all["test"])
}

func TestGetPlayer(t *testing.T) {

	c := cache.New()
	p := c.GetPlayer("foo")
	// Not initialized yet => empty struct
	assert.Equal(t, cache.PlayerStates{}, p)
}

func TestGetPlayerPositions(t *testing.T) {

	c := cache.New()
	prev := c.GetPrevPlayerPos("foo")
	curr := c.GetCurrPlayerPos("foo")
	// Player not yet initialized, so cache will return an empty struct
	assert.Equal(t, types.Position{}, prev)
	assert.Equal(t, types.Position{}, curr)

}

func TestPut(t *testing.T) {

	c := cache.New()
	id := "foo"

	// First make sure cache is empty on init
	assert.Equal(t, cache.Memory{}, c.GetAll())
	assert.Equal(t, cache.PlayerStates{}, c.GetPlayer(id))

	// Put valid move request to cache
	moveReqPrefix := "MOVE_"
	moveReqData := fmt.Sprintf(`{"pos":{"posX":%d,"posY":%d},"destPos":{"destPosX":%d,"destPosY":%d}}`, 1, 1, 2, 2)
	c.Put(id, []byte(moveReqPrefix+moveReqData))

	p := c.GetPlayer(id)

	assert.NotEqual(t, cache.PlayerStates{}, p)

}
