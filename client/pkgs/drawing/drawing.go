package center

import (
	"fmt"
	c "ticoma/client/pkgs/constants"
	"ticoma/client/pkgs/player"

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
func DrawPlayers(world *rl.RenderTexture2D, p player.ClientPlayer) {
	cheMap := p.InternalPlayer.GetCache()
	rl.BeginTextureMode(*world)
	for _, player := range *cheMap {
		pos := player.Curr.Position
		rl.DrawRectangleRec(rl.Rectangle{X: float32(pos.X) * c.BLOCK_SIZE, Y: float32(pos.Y) * c.BLOCK_SIZE, Width: c.BLOCK_SIZE, Height: c.BLOCK_SIZE}, rl.Purple)
	}
	rl.EndTextureMode()
}

// (Tmp) draws current coordinates on the map
func DrawCoordinates(p player.ClientPlayer, x float32, y float32) {
	pPos := p.InternalPlayer.GetPos().Position
	rl.DrawTextEx(c.DEFAULT_FONT, fmt.Sprintf("<%d, %d>", pPos.X, pPos.Y), rl.Vector2{X: x, Y: y}, c.DEFAULT_FONT_SIZE, 0, rl.Blue)
}

// (Tmp) draws a block from block sprite
func DrawBlock(blockTxt *rl.Texture2D, id int, mapX float32, mapY float32) {
	blockRec := rl.Rectangle{X: float32(id) * c.BLOCK_SIZE, Y: 0, Width: c.BLOCK_SIZE, Height: c.BLOCK_SIZE}
	rl.DrawTextureRec(*blockTxt, blockRec, rl.Vector2{X: mapX * c.BLOCK_SIZE, Y: mapY * c.BLOCK_SIZE}, rl.White)
}
