package client

import (
	"ticoma/client/packages/actions"
	"ticoma/client/packages/camera"
	"ticoma/client/packages/drawing/panels/center"
	"ticoma/client/packages/drawing/panels/left"
	"ticoma/client/packages/drawing/panels/right"

	// gamemap "ticoma/client/packages/game_map"
	"ticoma/client/packages/input/keyboard"
	intf "ticoma/client/packages/interfaces"
	internal_player "ticoma/internal/packages/player"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func Main(pc chan internal_player.Player) {

	// Setup resolution, scaling
	// screenConf = utils.GetScreenInfo()
	// Tmp, testing > fixed res
	screenConf := intf.ScreenInfo{
		Width:       1920,
		Height:      1080,
		RefreshRate: 60,
	}
	SIDE_PANEL_WIDTH := int32((screenConf.Width / 5))

	// Load assets
	// font := rl.LoadFont("../client/assets/fonts/Consolas.ttf") // TODO: Fix font
	icon := rl.LoadImage("../client/assets/logo/ticoma-logo-64.png")
	spawnImg := rl.LoadImage("../client/assets/textures/map/spawn.png")

	// Setup window, raylib
	rl.InitWindow(int32(screenConf.Width), int32(screenConf.Height), "Ticoma Client")
	// rl.ToggleFullscreen()
	defer rl.CloseWindow()
	rl.SetTraceLog(4) // Disable unnecessary raylib logs
	rl.SetTargetFPS(60)
	rl.SetWindowIcon(*icon)

	// Setup sprites
	// blockSprite := rl.LoadTextureFromImage(blocks)

	// Setup game
	cam := camera.New()
	playerMoved := false
	spawnMapSize := 25 // In blocks
	// testMap := gamemap.NewMap(spawnMapSize, spawnMapSize)
	// testMap.GenerateRandomMap()

	// Setup panels
	centerPanel := rl.LoadRenderTexture(spawnImg.Width, spawnImg.Height)
	playerPanel := rl.LoadRenderTexture(int32(screenConf.Width), int32(screenConf.Height))
	rightPanel := rl.LoadRenderTexture(SIDE_PANEL_WIDTH, int32(screenConf.Height))
	leftPanel := rl.LoadRenderTexture(SIDE_PANEL_WIDTH, int32(screenConf.Height))

	left.DrawLeftPanelSkeleton(&leftPanel, SIDE_PANEL_WIDTH, int32(screenConf.Height))
	right.DrawRightPanelSkeleton(&rightPanel, SIDE_PANEL_WIDTH, int32(screenConf.Height))
	center.TestMap(&centerPanel, spawnImg, &screenConf)
	center.DrawSelfPlayer(&playerPanel, 0, screenConf.Width, screenConf.Height)

	// Display loading scene, wait for internal
	p := <-pc

	actions.InitPlayer(p, &playerMoved, spawnMapSize/2, spawnMapSize/2)

	for !rl.WindowShouldClose() {

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		rl.BeginMode2D(cam.Camera2D)
		xOffset := 2.5 * 64
		yOffset := 2.5 * 64
		rl.DrawTextureRec(centerPanel.Texture, rl.Rectangle{X: 0, Y: 0, Width: float32(spawnImg.Width), Height: float32(-spawnImg.Height)}, rl.Vector2{X: float32(xOffset), Y: float32(yOffset)}, rl.White)
		rl.EndMode2D()

		rl.DrawTextureRec(playerPanel.Texture, rl.Rectangle{X: 0, Y: 0, Width: float32(screenConf.Width), Height: float32(-screenConf.Height)}, rl.Vector2{X: 0, Y: 0}, rl.White)
		rl.DrawTextureRec(leftPanel.Texture, rl.Rectangle{X: 0, Y: 0, Width: float32(leftPanel.Texture.Width), Height: float32(-leftPanel.Texture.Height)}, rl.Vector2{X: 0, Y: 0}, rl.White)
		rl.DrawTextureRec(rightPanel.Texture, rl.Rectangle{X: 0, Y: 0, Width: float32(rightPanel.Texture.Width), Height: float32(-rightPanel.Texture.Height)}, rl.Vector2{X: float32(int32(screenConf.Width) - SIDE_PANEL_WIDTH), Y: 0}, rl.White)

		keyboard.HandleKeyboardMoveInput(p, cam, &playerMoved)

		rl.EndDrawing()

	}

}
