package client

import (
	// "ticoma/client/packages/actions"
	// "ticoma/client/packages/drawing/player"
	// "ticoma/client/packages/input/keyboard"
	// "ticoma/client/packages/utils"
	// gamemap "ticoma/client/packages/game_map"
	"ticoma/client/packages/drawing/panels/left"
	"ticoma/client/packages/drawing/panels/right"
	intf "ticoma/client/packages/interfaces"

	internal_player "ticoma/internal/packages/player"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Conf
var screenConf intf.ScreenInfo

// Misc
var playerMoved = false

// Panels
var leftPanel rl.RenderTexture2D
var centerPanel rl.RenderTexture2D
var rightPanel rl.RenderTexture2D

func Main(pc chan internal_player.Player) {

	// Load assets
	// font := rl.LoadFont("../client/assets/fonts/Consolas.ttf") // TODO: Fix font
	icon := rl.LoadImage("../client/assets/logo/ticoma-logo-64.png")

	// Setup resolution, scaling
	// screenConf = utils.GetScreenInfo()
	// Tmp, testing > fixed res
	screenConf := intf.ScreenInfo{
		Width:       1920,
		Height:      1080,
		RefreshRate: 60,
	}

	// Panel scaling config
	SIDE_PANEL_WIDTH := int32((screenConf.Width / 4))

	// Setup window, raylib
	rl.InitWindow(int32(screenConf.Width), int32(screenConf.Height), "Ticoma Client")
	// rl.ToggleFullscreen()
	defer rl.CloseWindow()
	rl.SetTraceLog(4) // Disable unnecessary raylib logs
	rl.SetTargetFPS(60)
	rl.SetWindowIcon(*icon)

	// Setup panels
	centerPanel = rl.LoadRenderTexture(int32(screenConf.Width), int32(screenConf.Height))
	rightPanel = rl.LoadRenderTexture(SIDE_PANEL_WIDTH, int32(screenConf.Height))
	leftPanel = rl.LoadRenderTexture(SIDE_PANEL_WIDTH, int32(screenConf.Height))

	left.DrawLeftPanelSkeleton(&leftPanel, 0, 0, SIDE_PANEL_WIDTH, int32(screenConf.Height))
	right.DrawRightPanelSkeleton(&rightPanel, SIDE_PANEL_WIDTH, int32(screenConf.Height))

	// Display loading scene, wait for internal
	// p := <-pc

	// init player pos @ random pos
	// randX := utils.RandRange(0, 10)
	// randY := utils.RandRange(0, 10)
	// actions.MovePlayer(p, &playerMoved, randX, randY, randX, randY)

	// Misc
	// ver := utils.GetCommitHash()[0:6]

	for !rl.WindowShouldClose() {

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		rl.DrawTexture(centerPanel.Texture, 0, 0, rl.White)
		rl.DrawTexture(leftPanel.Texture, 0, 0, rl.White)
		rl.DrawTexture(rightPanel.Texture, (int32(screenConf.Width) - SIDE_PANEL_WIDTH), 0, rl.White)

		rl.EndDrawing()

	}

}

// rl.DrawText("center", 10, 10, 30, rl.White)
// rl.EndTextureMode()

// rl.DrawTexture(centerPanelTxt, 0, 0, rl.White)
// rl.BeginTextureMode(centerPanel)

// rl.EndTextureMode()
// Draw main map
// gamemap.DrawGameMap(c.GAME_MAP_START_X, c.GAME_MAP_START_Y, c.GAME_MAP_BLOCK_SIZE, c.GAME_MAP_SIZE_IN_BLOCKS)

// keyboard.HandleKeyboardMoveInput(p, &playerMoved)

// for id, pos := range *p.GetPlayersPos() {
// 	player.DrawPlayer(id, pos.X, pos.Y)
// }

// rl.ClearBackground(rl.RayWhite)

// Draw top-left info
// rl.DrawText("ticoma git-"+ver, 3, 0, 24, rl.Black)
// rl.DrawText("peerid-"+p.GetPeerID(), 3, 30, 24, rl.Black)

// rl.DrawTextEx(font, "Test", rl.Vector2{X: 100, Y: 100}, 32, 0, rl.Black)
