package center

import (
	c "ticoma/client/packages/constants"
	gamemap "ticoma/client/packages/game_map"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//
// NOTE: Everything that is drawn on the center panel
// should be precisely centered to avoid visual glitches
//

// Draw a complete map
func DrawMap(panel *rl.RenderTexture2D, gameMap *gamemap.GameMap, blockSprite *rl.Texture2D) {
	rl.BeginTextureMode(*panel)
	for i := 0; i < len(gameMap.FrontLayer); i++ {
		for j := 0; j < len(gameMap.FrontLayer[0]); j++ {
			blockId := gameMap.FrontLayer[i][j].BlockId
			spritePos := rl.Rectangle{X: float32(blockId * c.BLOCK_SIZE), Y: 0, Width: c.BLOCK_SIZE, Height: c.BLOCK_SIZE}
			gameMapPos := rl.Vector2{X: float32(c.BLOCK_SIZE * j), Y: float32(c.BLOCK_SIZE * i)}
			rl.DrawTextureRec(*blockSprite, spritePos, gameMapPos, rl.White)
		}
	}
	rl.EndTextureMode()
}

// Draw own, local player instance
func DrawSelfPlayer(panel *rl.RenderTexture2D, id int, screenWidth int, screenHeight int) {
	rl.BeginTextureMode(*panel)
	rl.ClearBackground(rl.Black)
	rl.DrawRectangleRec(rl.Rectangle{X: float32(screenWidth/2) - (c.BLOCK_SIZE / 2), Y: float32(screenHeight/2) - (c.BLOCK_SIZE / 2), Width: c.BLOCK_SIZE, Height: c.BLOCK_SIZE}, rl.White)
	// rl.DrawRectanglePro(rl.Rectangle{X: float32(screenWidth / 2), Y: float32(screenHeight / 2)}, rl.Vector2{X: 50, Y: 5}, 0, rl.Green)
	rl.EndTextureMode()
}

// Draw all other players based on their current positions (taken from local cache)
func DrawPlayer() {}
