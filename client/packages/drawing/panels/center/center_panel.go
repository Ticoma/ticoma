package center

import (
	c "ticoma/client/packages/constants"
	intf "ticoma/client/packages/interfaces"

	// gamemap "ticoma/client/packages/game_map"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Draw a complete map
// func DrawMap(panel *rl.RenderTexture2D, gameMap *gamemap.GameMap, blockSprite *rl.Texture2D) {
// 	rl.BeginTextureMode(*panel)
// 	for i := 0; i < len(gameMap.FrontLayer); i++ {
// 		for j := 0; j < len(gameMap.FrontLayer[0]); j++ {
// 			blockId := gameMap.FrontLayer[i][j].BlockId
// 			spritePos := rl.Rectangle{X: float32(blockId * c.BLOCK_SIZE), Y: 0, Width: c.BLOCK_SIZE, Height: c.BLOCK_SIZE}
// 			gameMapPos := rl.Vector2{X: float32(c.BLOCK_SIZE * j), Y: float32(c.BLOCK_SIZE * i)}
// 			rl.DrawTextureRec(*blockSprite, spritePos, gameMapPos, rl.White)
// 		}
// 	}
// 	rl.EndTextureMode()
// }

func TestMap(panel *rl.RenderTexture2D, img *rl.Image, screenConf *intf.ScreenInfo) {
	txt := rl.LoadTextureFromImage(img)
	rl.BeginTextureMode(*panel)
	rl.DrawTextureRec(txt, rl.Rectangle{X: 0, Y: 0, Width: float32(txt.Width), Height: float32(txt.Height)}, rl.Vector2{X: 0, Y: 0}, rl.White)
	rl.EndTextureMode()
}

// Draw own, local player instance
func DrawSelfPlayer(panel *rl.RenderTexture2D, id int, screenWidth int, screenHeight int) {
	rl.BeginTextureMode(*panel)
	rl.DrawRectangleRec(rl.Rectangle{X: float32(screenWidth/2) - (c.BLOCK_SIZE / 2), Y: float32(screenHeight/2) - (c.BLOCK_SIZE / 2), Width: c.BLOCK_SIZE, Height: c.BLOCK_SIZE}, rl.White)
	rl.EndTextureMode()
}

// Draw all other players based on their current positions (taken from local cache)
func DrawPlayer() {}
