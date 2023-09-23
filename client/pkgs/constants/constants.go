package constants

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type ScreenInfo struct {
	MonitorId   int
	Width       int32
	Height      int32
	RefreshRate int
}

// All constants here assume @1.0 UI scaling factor and 1080p res,
// Should be multiplied by scale when drawing on different resolutions

// Screen stuff
var SCREEN *ScreenInfo

// Fonts
var DEFAULT_FONT rl.Font

// Game
var SCALE float32
var BLOCK_SIZE float32 = 64
var MOVE_COOLDOWN_IN_MS int = 300

// UI , Scaling, Padding
var DEFAULT_FONT_SIZE float32 = 16
var DEFAULT_PADDING float32 = 12
var SIDE_PANEL_PADDING float32 = DEFAULT_PADDING

// Colors
var COLOR_PANEL_BG rl.Color = rl.NewColor(40, 40, 40, 255)
var COLOR_PANEL_CONTENT rl.Color = rl.NewColor(60, 60, 60, 255)
var COLOR_PANEL_OUTLINE rl.Color = rl.NewColor(102, 191, 255, 255) // Skyblue
var COLOR_PANEL_OUTLINE_TRANSPARENT = rl.NewColor(102, 191, 255, 125)
var COLOR_PANEL_TEXT rl.Color = rl.NewColor(240, 240, 240, 255)

// Block runes and pathfinding stuff
var COLLISION_BLOCK_RUNE = 'x'
var WALKABLE_BLOCK_RUNE = ' '
