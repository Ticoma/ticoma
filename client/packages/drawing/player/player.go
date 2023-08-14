package player

import (
	"strconv"
	c "ticoma/client/packages/constants"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Draws a player with specified at position on game map
//
// id - playerId, posX, posY - map position in blocks
func DrawPlayer(id int, posX int, posY int) {
	rl.DrawRectangle(int32(c.GAME_MAP_START_X+(c.GAME_MAP_BLOCK_SIZE*posX)), int32(c.GAME_MAP_START_Y+(c.GAME_MAP_BLOCK_SIZE*posY)), c.GAME_MAP_BLOCK_SIZE, c.GAME_MAP_BLOCK_SIZE, rl.Black)
	ids := strconv.Itoa(id)
	rl.DrawText(ids, int32(c.GAME_MAP_START_X+(c.GAME_MAP_BLOCK_SIZE*posX)+5), int32(c.GAME_MAP_START_Y+(c.GAME_MAP_BLOCK_SIZE*posY)+5), 18, rl.Red)
}
