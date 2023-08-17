package left

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func DrawLeftPanelSkeleton(panel *rl.RenderTexture2D, width int32, height int32) {
	rl.BeginTextureMode(*panel)
	rl.ClearBackground(rl.Black)
	rl.DrawRectangleRec(rl.Rectangle{X: 0, Y: 0, Width: float32(width), Height: float32(height)}, rl.DarkGray)
	rl.DrawText("Chat", 10, 10, 32, rl.White)
	rl.EndTextureMode()
}
