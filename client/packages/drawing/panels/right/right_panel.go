package right

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func DrawRightPanelSkeleton(panel *rl.RenderTexture2D, width int32, height int32) {
	rl.BeginTextureMode(*panel)
	rl.ClearBackground(rl.Black)
	rl.DrawRectangle(0, 0, width, height, rl.DarkGray)
	rl.EndTextureMode()
}
