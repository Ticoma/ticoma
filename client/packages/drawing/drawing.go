package center

import (
	c "ticoma/client/packages/constants"
	internal_player "ticoma/internal/packages/player"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Draw empty world map from texture
//
// TODO: Draw map from json-like object imported from texture editor instead of img
func DrawMap(world *rl.RenderTexture2D, txt *rl.Texture2D, zoom float32) {
	rl.BeginTextureMode(*world)
	rl.DrawTextureRec(*txt, rl.Rectangle{X: 0, Y: 0, Width: float32(txt.Width) * zoom, Height: float32(txt.Height) * zoom}, rl.Vector2{X: 0, Y: 0}, rl.White)
	rl.EndTextureMode()
}

// Draw all online players on world texture
func DrawPlayers(world *rl.RenderTexture2D, p internal_player.Player, zoom float32) {
	cMap := p.GetCache().CacheStore.Store
	rl.BeginTextureMode(*world)
	for _, elem := range cMap {
		if p.GetId() == elem[0].PlayerId {
			continue // Ignore self
		}
		pos := elem[1].Position // get current pos
		rl.DrawRectangleRec(rl.Rectangle{X: float32(pos.X) * c.BLOCK_SIZE, Y: float32(pos.Y) * c.BLOCK_SIZE, Width: c.BLOCK_SIZE * zoom, Height: c.BLOCK_SIZE * zoom}, rl.Purple)
	}
	rl.EndTextureMode()
}
