package constants

import rl "github.com/gen2brain/raylib-go/raylib"

// All constants here assume @1.0 UI scaling factor and 1080p res,
// Should be multiplied by scale when drawing on different resolutions

// Colors
var COLOR_PANEL_BG rl.Color = rl.NewColor(40, 40, 40, 255)
var COLOR_PANEL_CONTENT rl.Color = rl.NewColor(60, 60, 60, 255)
var COLOR_PANEL_OUTLINE rl.Color = rl.NewColor(102, 191, 255, 255) // Skyblue
var COLOR_PANEL_OUTLINE_TRANSPARENT rl.Color = rl.NewColor(102, 191, 255, 150)
var COLOR_PANEL_TEXT rl.Color = rl.NewColor(220, 220, 220, 255)

const (
	// ===============================================
	//Game

	BLOCK_SIZE = 64

	// ===============================================
	// UI

	// Spacing
	DEFAULT_FONT_SIZE          = 16
	DEFAULT_PADDING    float32 = 16
	SIDE_PANEL_PADDING float32 = DEFAULT_PADDING // padding around side panels
)
