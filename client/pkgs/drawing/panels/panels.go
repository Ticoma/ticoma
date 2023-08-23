package panels

import (
	c "ticoma/client/pkgs/constants"
	"ticoma/client/pkgs/input/mouse"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//
// Side panels
//

type SpaceTaken struct {
	Top    float32
	Bottom float32
}

type Tabs map[int][2]string // key=Index , val=[Full name, Shortcut]

type SidePanel struct {
	Txt        *rl.RenderTexture2D // Render Texture
	Pos        rl.Vector2          // Where to render the panel at last
	BgColor    *rl.Color           // Skeleton Bg Color
	Tabs       Tabs                // All the tabs available
	ActiveTab  int                 // Index of currently active tab
	SpaceTaken                     // How much space is already taken in this panel
	// ^ NOTE: This is very useful for avoiding ugly math when drawing inside panels, but
	// be mindful of when and how this changes, especially while drawing panel elements in bulk
}

// Initializes a new side panel struct, no drawing
// Width, height, x, y will be used for drawing the skeleton and all content inside it
func New(rt2d *rl.RenderTexture2D, width float32, height float32, x float32, y float32, bgColor *rl.Color, tabs Tabs) *SidePanel {
	return &SidePanel{
		Txt: rt2d,
		Pos: rl.Vector2{
			X: x,
			Y: y,
		},
		BgColor:    bgColor,
		Tabs:       tabs,
		ActiveTab:  0,
		SpaceTaken: SpaceTaken{},
	}
}

// NOTE: I dont think this is going to be needed, but keeping it just in case
// // Clears the panel texture with panel's bgColor
// func (sp *SidePanel) Clear() {
// 	rl.BeginTextureMode(*sp.Txt)
// 	rl.ClearBackground(*sp.BgColor)
// 	rl.EndTextureMode()
// }

// Draw the side panel skeleton (plain bg rectangle)
func (sp *SidePanel) DrawSkeleton() {
	rl.BeginTextureMode(*sp.Txt)
	rl.DrawRectangleRec(rl.Rectangle{X: sp.Pos.X, Y: sp.Pos.Y, Width: float32(sp.Txt.Texture.Width), Height: float32(sp.Txt.Texture.Height)}, *sp.BgColor)
	rl.EndTextureMode()
}

// Draw panel title (@top)
func (sp *SidePanel) DrawPanelTitle(font *rl.Font, fontSize float32) {

	titleCtnH := float32(sp.Txt.Texture.Height / 12)
	titleSize := rl.MeasureTextEx(*font, sp.Tabs[sp.ActiveTab][0], c.DEFAULT_FONT_SIZE, 0)

	rl.BeginTextureMode(*sp.Txt)
	rl.DrawRectangleRec(rl.Rectangle{X: c.SIDE_PANEL_PADDING, Y: c.SIDE_PANEL_PADDING, Width: float32(sp.Txt.Texture.Width) - 2*c.SIDE_PANEL_PADDING, Height: titleCtnH - float32(2*c.SIDE_PANEL_PADDING)}, rl.Gray)
	rl.DrawTextEx(*font, sp.Tabs[sp.ActiveTab][0], rl.Vector2{X: float32(sp.Txt.Texture.Width/2) - titleSize.X/2, Y: titleCtnH/2 - titleSize.Y/2}, fontSize, 0, rl.White)
	rl.EndTextureMode()

	// Update taken space
	sp.SpaceTaken.Top += titleCtnH + 2*c.SIDE_PANEL_PADDING
}

// Draw panel tab switcher (@bottom)
func (sp *SidePanel) DrawPanelTabs(font *rl.Font, fontSize float32) {

	tabsCtnH := float32(sp.Txt.Texture.Height / 12)
	tabsCtnY := float32(sp.Txt.Texture.Height) - tabsCtnH + c.SIDE_PANEL_PADDING
	tabsContentW := float32(sp.Txt.Texture.Width) - 2*c.SIDE_PANEL_PADDING
	singleTabCtnW := tabsContentW / float32(len(sp.Tabs))

	rl.BeginTextureMode(*sp.Txt)
	// Draw container // Off for now
	// rl.DrawRectangleRec(rl.Rectangle{X: c.SIDE_PANEL_PADDING, Y: tabsCtnY, Width: float32(sp.Txt.Texture.Width) - 2*c.SIDE_PANEL_PADDING, Height: tabsCtnH - 2*c.SIDE_PANEL_PADDING}, rl.Gray)
	for i, elem := range sp.Tabs {
		// Draw tab containers
		tabX := c.SIDE_PANEL_PADDING + float32(i)*singleTabCtnW
		tabRec := rl.Rectangle{
			X:      tabX,
			Y:      tabsCtnY,
			Width:  singleTabCtnW,
			Height: tabsCtnH - 2*c.SIDE_PANEL_PADDING,
		}
		rl.DrawRectangleRec(tabRec, rl.Gray)
		// Calculate shortcut text size
		shSize := rl.MeasureTextEx(*font, elem[1], c.DEFAULT_FONT_SIZE, 0)
		// Draw tab shortcuts
		rl.DrawTextEx(*font, elem[1], rl.Vector2{
			X: tabX + singleTabCtnW/2 - shSize.X/2,
			Y: tabsCtnY + ((tabsCtnH - 2*c.SIDE_PANEL_PADDING) / 2) - shSize.Y/2,
		}, c.DEFAULT_FONT_SIZE, 0, rl.White)

		isHoveringTab := mouse.IsMouseHoveringRec(&tabRec)

		// Draw outline for currently active tab
		if sp.ActiveTab == i {
			rl.DrawRectangleLinesEx(tabRec, 2, rl.SkyBlue)
		} else if sp.ActiveTab != i && isHoveringTab {
			rl.DrawRectangleLinesEx(tabRec, 2, rl.Black)
		}

		// Handle tab switching logic
		if isHoveringTab {
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				sp.ActiveTab = i
			}
		}
	}
	rl.EndTextureMode()

	// Update space
	sp.SpaceTaken.Bottom += tabsCtnH + 2*c.SIDE_PANEL_PADDING
}
