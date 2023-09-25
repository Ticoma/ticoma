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
var playerNick = "player"

func TestGetAll(t *testing.T) {

	// Get all on just initialized cache
	c := cache.New()
	all := c.GetAll()
	allVal := *all
	assert.Equal(t, &cache.Memory{}, all)
	assert.Equal(t, cache.PlayerStates{}, allVal["test"])
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
	assert.Equal(t, &types.PlayerPosition{}, prev)
	assert.Equal(t, &types.PlayerPosition{}, curr)

	// Register should set player's pos and destPos to {13, 13}
	_, pfx, err := c.Put(playerID, []byte(fmt.Sprintf(`REGISTER_{"nickname":"%s"}`, playerNick)))

	assert.Equal(t, "REGISTER_", pfx)
	assert.NoError(t, err)

	expectedPos := &types.PlayerPosition{
		Timestamp:    0,
		Position:     types.Position{X: cache.SPAWN_POS_X, Y: cache.SPAWN_POS_Y},
		DestPosition: types.DestPosition{X: cache.SPAWN_POS_X, Y: cache.SPAWN_POS_Y},
	}

	prev = c.GetPrevPlayerPos(playerID)
	curr = c.GetCurrPlayerPos(playerID)
	prev.Timestamp, curr.Timestamp = 0, 0 // Ignore the timestamp

	assert.Equal(t, expectedPos.Position.X, prev.Position.X)
	assert.Equal(t, expectedPos.Position.X, curr.Position.X)
	assert.Equal(t, expectedPos.Position.Y, prev.Position.Y)
	assert.Equal(t, expectedPos.Position.Y, curr.Position.Y)
	assert.Equal(t, expectedPos.DestPosition.X, prev.DestPosition.X)
	assert.Equal(t, expectedPos.DestPosition.X, curr.DestPosition.X)
	assert.Equal(t, expectedPos.DestPosition.Y, prev.DestPosition.Y)
	assert.Equal(t, expectedPos.DestPosition.Y, curr.DestPosition.Y)

}

func TestMove(t *testing.T) {

	c := cache.New()

	// First make sure cache is empty on init
	assert.Equal(t, &cache.Memory{}, c.GetAll())
	assert.Equal(t, &cache.PlayerStates{}, c.GetPlayer(playerID))

	// The player doesn't exist, so we need to register first.
	_, _, err := c.Put(playerID, []byte(fmt.Sprintf(`REGISTER_{"nickname":"%s"}`, playerNick)))

	p := c.GetPlayer(playerID)

	// Check if registration OK
	assert.NoError(t, err)
	assert.Equal(t, true, p.Curr.IsOnline)

	// Should be at spawn pos
	assert.Equal(t, cache.SPAWN_POS_X, p.Curr.PlayerGameData.Position.X)
	assert.Equal(t, cache.SPAWN_POS_Y, p.Curr.PlayerGameData.Position.Y)

	// Timestamp in those positions should not be 0 if cache initialized correctly
	assert.NotEqual(t, 0, p.Prev.PlayerGameData.Timestamp)
	assert.NotEqual(t, 0, p.Curr.PlayerGameData.Timestamp)

	// Put valid move request to cache
	moveReqPos := cache.SPAWN_POS_X
	moveReqDestPos := cache.SPAWN_POS_X + 1
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
	_, _, err = c.Put(playerID, []byte(invalidMoveReqData))

	// Should throw
	assert.Error(t, err)

	// Logout and try to move
	logoutReqData := []byte("LOGOUT_")
	_, _, err = c.Put(playerID, logoutReqData)

	assert.NoError(t, err)

	// Try to move when logged out

	anotherMoveReqData := fmt.Sprintf(`MOVE_{"pos":{"posX":%d,"posY":%d},"destPos":{"destPosX":%d,"destPosY":%d}}`, 14, 14, 14, 14)
	_, _, err = c.Put(playerID, []byte(anotherMoveReqData))

	// The mv req above should get rejected (not logged in)
	assert.Error(t, err)
	assert.Equal(t, cache.SPAWN_POS_X, c.GetPlayer(playerID).Curr.Position.X)
	assert.Equal(t, cache.SPAWN_POS_Y, c.GetPlayer(playerID).Curr.Position.Y)
	assert.NotEqual(t, &cache.Memory{}, c.GetAll())

}

func TestGetNickname(t *testing.T) {

	c := cache.New()
	localId := "123000123"
	localNickname := "tester0"

	nick := c.GetNickname(localId)

	assert.Equal(t, "", *nick)

	_, _, err := c.Put(localId, []byte(fmt.Sprintf(`REGISTER_{"nickname":"%s"}`, localNickname)))

	nick = c.GetNickname(localId)

	assert.NoError(t, err)
	assert.Equal(t, localNickname, *nick)
}
