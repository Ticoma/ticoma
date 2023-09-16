package tests

import (
	"fmt"
	"testing"
	"ticoma/internal/pkgs/gamenode/cache"
	"ticoma/types"
	"time"

	assert "github.com/stretchr/testify/assert"
)

var playerId = "foo"

func TestGetAll(t *testing.T) {

	// Get all on just initialized cache
	c := cache.New()
	all := c.GetAll()
	assert.Equal(t, cache.Memory{}, all)
	assert.Equal(t, cache.PlayerStates{}, all["test"])
}

func TestGetPlayer(t *testing.T) {

	c := cache.New()
	p := c.GetPlayer(playerID)
	// Not initialized yet => ptr to empty struct
	assert.Equal(t, cache.PlayerStates{}, *p)
}

func TestGetPlayerPositions(t *testing.T) {

	c := cache.New()
	prev := c.GetPrevPlayerPos(playerID)
	curr := c.GetCurrPlayerPos(playerID)
	// Player not yet initialized, so cache will return an empty struct
	assert.Equal(t, types.Position{}, prev)
	assert.Equal(t, types.Position{}, curr)

}

func TestMove(t *testing.T) {

	c := cache.New()

	// First make sure cache is empty on init
	assert.Equal(t, cache.Memory{}, c.GetAll())
	assert.Equal(t, &cache.PlayerStates{}, c.GetPlayer(playerID))

	// The player doesn't exist, so we need to register first.
	_, err := c.Put(playerID, []byte("REGISTER_"))

	p := c.GetPlayer(playerID)

	// Check if registration OK
	assert.NoError(t, err)
	assert.Equal(t, true, p.Curr.IsOnline)

	// Should spawn at {13, 13} or whatever is set in Cache
	assert.Equal(t, 13, p.Curr.PlayerGameData.Position.X)
	assert.Equal(t, 13, p.Curr.PlayerGameData.Position.Y)

	// Timestamp in those positions should not be 0 if cache initialized correctly
	assert.NotEqual(t, 0, p.Prev.PlayerGameData.Timestamp)
	assert.NotEqual(t, 0, p.Curr.PlayerGameData.Timestamp)

	// Put valid move request to cache (pos 1, 1 -> destPos 2,2)
	moveReqPos := 13
	moveReqDestPos := 14
	moveReqPrefix := "MOVE_"
	moveReqData := fmt.Sprintf(`{"pos":{"posX":%d,"posY":%d},"destPos":{"destPosX":%d,"destPosY":%d}}`, moveReqPos, moveReqPos, moveReqDestPos, moveReqDestPos)

	time.Sleep(time.Millisecond * 500)
	c.Put(playerID, []byte(moveReqPrefix+moveReqData))

	// Check if cache processed move req
	player := c.GetPlayer(playerID)
	assert.Equal(t, moveReqPos, player.Curr.Position.X)
	assert.Equal(t, moveReqPos, player.Curr.Position.Y)
	assert.Equal(t, moveReqDestPos, player.Curr.DestPosition.X)
	assert.Equal(t, moveReqDestPos, player.Curr.DestPosition.Y)

	// Send a silly bonkers request
	invalidMoveReqData := fmt.Sprintf(`MOVE_{"pos":{"posX":%d,"posY":%d},"destPos":{"destPosX":%d,"destPosY":%d}}`, 44, 41, -423, 12)
	_, err = c.Put(playerID, []byte(invalidMoveReqData))

	// Should throw
	assert.Error(t, err)

	// Logout and try to move
	logoutReqData := []byte("LOGOUT_")
	_, err = c.Put(playerID, logoutReqData)

	assert.NoError(t, err)

	// Try to move when logged out

	anotherMoveReqData := fmt.Sprintf(`MOVE_{"pos":{"posX":%d,"posY":%d},"destPos":{"destPosX":%d,"destPosY":%d}}`, 14, 14, 14, 14)
	_, err = c.Put(playerID, []byte(anotherMoveReqData))

	// The mv req above should get rejected (not logged in)
	assert.Error(t, err)
	assert.Equal(t, 13, c.GetPlayer(playerID).Curr.Position.X)
	assert.Equal(t, 13, c.GetPlayer(playerID).Curr.Position.Y)

}
