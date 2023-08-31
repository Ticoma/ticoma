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
	SpaceTaken                     // How much space is already taken in this panel (includes vertical paddings)
	// ^ NOTE: This is very useful for avoiding ugly math when drawing inside panels, but
	// be mindful of when and how this changes and keep track of the drawing order
}

// Initializes a new side panel struct, no drawing
// Width, height, x, y will be used for drawing the skeleton and all content inside it
func New(rt2d *rl.RenderTexture2D, width float32, height float32, x float32, y float32, bgColor *rl.Color, tabs Tabs) SidePanel {
	return SidePanel{
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
	skeletonRec := rl.Rectangle{X: 0, Y: 0, Width: float32(sp.Txt.Texture.Width), Height: float32(sp.Txt.Texture.Height)}
	rl.DrawRectangleRec(skeletonRec, *sp.BgColor)
	rl.EndTextureMode()
}

// Draw panel title (@top)
func (sp *SidePanel) DrawPanelTitle() {

	titleCtnH := float32(sp.Txt.Texture.Height / 12) // Container full height (padding inside)
	titleSize := rl.MeasureTextEx(c.DEFAULT_FONT, sp.Tabs[sp.ActiveTab][0], c.DEFAULT_FONT_SIZE, 0)

	rl.BeginTextureMode(*sp.Txt)
	rl.DrawRectangleRec(rl.Rectangle{X: c.SIDE_PANEL_PADDING, Y: c.SIDE_PANEL_PADDING, Width: float32(sp.Txt.Texture.Width) - 2*c.SIDE_PANEL_PADDING, Height: titleCtnH - float32(2*c.SIDE_PANEL_PADDING)}, c.COLOR_PANEL_CONTENT)
	rl.DrawTextEx(c.DEFAULT_FONT, sp.Tabs[sp.ActiveTab][0], rl.Vector2{X: float32(sp.Txt.Texture.Width/2) - titleSize.X/2, Y: titleCtnH/2 - titleSize.Y/2}, c.DEFAULT_FONT_SIZE, 0, c.COLOR_PANEL_TEXT)
	rl.EndTextureMode()

	// Update taken space
	sp.SpaceTaken.Top = titleCtnH - c.SIDE_PANEL_PADDING
}

// Draw panel tab switcher (@bottom)
func (sp *SidePanel) DrawPanelTabs() {

	tabsCtnH := float32(sp.Txt.Texture.Height / 12)
	tabsCtnY := float32(sp.Txt.Texture.Height) - tabsCtnH + c.SIDE_PANEL_PADDING
	tabsContentW := float32(sp.Txt.Texture.Width) - 2*c.SIDE_PANEL_PADDING
	singleTabCtnW := tabsContentW / float32(len(sp.Tabs))

	rl.BeginTextureMode(*sp.Txt)
	for i, tab := range sp.Tabs {
		// Draw tab containers
		tabX := float32(i) * singleTabCtnW
		tabRec := rl.Rectangle{
			X:      tabX + c.SIDE_PANEL_PADDING,
			Y:      tabsCtnY,
			Width:  singleTabCtnW,
			Height: tabsCtnH - 2*c.SIDE_PANEL_PADDING,
		}
		rl.DrawRectangleRec(tabRec, c.COLOR_PANEL_CONTENT)
		// Calculate shortcut text size
		shSize := rl.MeasureTextEx(c.DEFAULT_FONT, tab[1], c.DEFAULT_FONT_SIZE, 0)
		// Draw tab shortcuts
		rl.DrawTextEx(c.DEFAULT_FONT, tab[1], rl.Vector2{
			X: tabX + singleTabCtnW/2,
			Y: tabsCtnY + ((tabsCtnH - 2*c.SIDE_PANEL_PADDING) / 2) - shSize.Y/2,
		}, c.DEFAULT_FONT_SIZE, 0, c.COLOR_PANEL_TEXT)

		// Get bounds of tab panel on the screen, not on the texture itself
		screenTabRecPos := rl.Rectangle{
			X:      sp.Pos.X + tabX + c.SIDE_PANEL_PADDING,
			Y:      tabsCtnY,
			Width:  singleTabCtnW,
			Height: tabsCtnH - 2*c.SIDE_PANEL_PADDING,
		}
		isHoveringTab := mouse.IsMouseHoveringRec(&screenTabRecPos)

		// Draw outline for currently active tab
		if sp.ActiveTab == i {
			rl.DrawRectangleLinesEx(tabRec, 2, c.COLOR_PANEL_OUTLINE)
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
	sp.SpaceTaken.Bottom = tabsCtnH
}
