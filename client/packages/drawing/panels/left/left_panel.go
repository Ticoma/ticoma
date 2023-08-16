package left

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func DrawLeftPanelSkeleton(panel *rl.RenderTexture2D, x int, y int, width int32, height int32) {
	rl.BeginTextureMode(*panel)
	rl.ClearBackground(rl.Black)
	rl.DrawRectangle(int32(x), int32(y), width, height, rl.DarkGray)
	rl.DrawText("Chat", 10, 10, 32, rl.White)
	rl.EndTextureMode()
}
