package panels

import (
	c "ticoma/client/packages/constants"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//
// Package responsible for drawing side panels and any
// items that belong inside them
//

// Draws a panel skeleton (solid color) at pos {x, y}
func DrawSidePanelSkeleton(panel *rl.RenderTexture2D, width float32, height float32, col *rl.Color) {
	rl.BeginTextureMode(*panel)

	rl.DrawRectangleRec(rl.Rectangle{X: 0, Y: 0, Width: width, Height: height}, *col)

	rl.EndTextureMode()
}

// Draw title at specified Y pos on panel
func DrawTitleBlock(panel *rl.RenderTexture2D, yPos float32, title string, font *rl.Font) {

	titleBlockHeight := float32(panel.Texture.Height / 12)
	titleSize := rl.MeasureTextEx(*font, title, c.DEFAULT_FONT_SIZE, 0)

	rl.BeginTextureMode(*panel)

	rl.DrawRectangleRec(rl.Rectangle{X: c.SIDE_PANEL_PADDING, Y: yPos + c.SIDE_PANEL_PADDING, Width: float32(panel.Texture.Width - 2*c.SIDE_PANEL_PADDING), Height: titleBlockHeight - float32(2*c.SIDE_PANEL_PADDING)}, rl.Gray)
	rl.DrawTextEx(*font, title, rl.Vector2{X: float32(panel.Texture.Width/2) - titleSize.X/2, Y: titleBlockHeight/2 - titleSize.Y/2}, c.DEFAULT_FONT_SIZE, 0, rl.Black)

	rl.EndTextureMode()
}
