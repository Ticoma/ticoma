package right

import (
	"ticoma/client/pkgs/drawing/scenes/game/panels"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type RightPanel struct {
	*panels.SidePanel
}

func New(width float32, height float32, renderPosX float32, renderPosY float32, bgColor *rl.Color, tabs panels.Tabs) *RightPanel {
	rt2d := rl.LoadRenderTexture(int32(width), int32(height)) // Create render texture for panel
	sp := panels.New(&rt2d, width, height, renderPosX, renderPosY, bgColor, tabs)
	return &RightPanel{
		SidePanel: &sp,
	}
}

// Renders the panel to screen
func (rp *RightPanel) RenderPanel() {
	rl.DrawTextureRec(rp.SidePanel.Txt.Texture, rl.Rectangle{X: 0, Y: 0, Width: float32(rp.Txt.Texture.Width), Height: float32(-rp.Txt.Texture.Height)}, rl.Vector2{X: rp.Pos.X, Y: rp.Pos.Y}, rl.White)
}

func (rp *RightPanel) DrawContent() {

	rp.SidePanel.DrawSkeleton()
	rp.SidePanel.DrawPanelTabs()

	// Content
	switch rp.Tabs[rp.ActiveTab] {
	// Inventory, map
	case rp.Tabs[0]:
		rp.DrawPanelTitle()
	// Settings
	case rp.Tabs[1]:
		rp.DrawPanelTitle()
	}

}
