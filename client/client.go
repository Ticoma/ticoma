package client

import (
	"ticoma/client/packages/actions"
	"ticoma/client/packages/camera"
	c "ticoma/client/packages/constants"
	"ticoma/client/packages/drawing/panels/left"
	"ticoma/client/packages/drawing/panels/right"
	"ticoma/client/packages/input/keyboard"
	"ticoma/client/packages/input/mouse"
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
	SIDE_PANEL_WIDTH := int32((screenConf.Width / 4))

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
	gameCam := camera.New()
	gameCam.Offset = rl.Vector2{X: float32(screenConf.Width / 2), Y: float32(screenConf.Height / 2)}
	playerMoved := false
	spawnMapSize := 25 // In blocks

	// Setup panels
	world := rl.LoadRenderTexture(spawnImg.Width, spawnImg.Height) // Fullscreen panel
	// game := rl.LoadRenderTexture(int32(screenConf.Width), int32(screenConf.Height))
	// player := rl.LoadRenderTexture(c.BLOCK_SIZE, c.BLOCK_SIZE)
	rightPanel := rl.LoadRenderTexture(SIDE_PANEL_WIDTH, int32(screenConf.Height)) // Side panels
	leftPanel := rl.LoadRenderTexture(SIDE_PANEL_WIDTH, int32(screenConf.Height))

	// tmp, Draw map on world
	spawnTxt := rl.LoadTextureFromImage(spawnImg)
	rl.BeginTextureMode(world)
	rl.ClearBackground(rl.NewColor(0, 0, 0, 0))
	rl.DrawTextureRec(spawnTxt, rl.Rectangle{
		X:      0,
		Y:      0,
		Width:  float32(spawnTxt.Width) * gameCam.Zoom,
		Height: float32(spawnTxt.Height) * gameCam.Zoom},
		rl.Vector2{
			X: 0,
			Y: 0,
		},
		rl.White)
	rl.EndTextureMode()

	// Wait for player conn
	p := <-pc

	// Draw textures
	left.DrawLeftPanelSkeleton(&leftPanel, SIDE_PANEL_WIDTH, int32(screenConf.Height))
	right.DrawRightPanelSkeleton(&rightPanel, SIDE_PANEL_WIDTH, int32(screenConf.Height))

	actions.InitPlayer(p, &playerMoved, spawnMapSize/2, spawnMapSize/2)

	for !rl.WindowShouldClose() {

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		// Draw game
		rl.BeginMode2D(gameCam.Camera2D)
		rl.DrawTextureRec(world.Texture, rl.Rectangle{X: 0, Y: 0, Width: float32(world.Texture.Width), Height: float32(-world.Texture.Height)}, rl.Vector2{X: 0, Y: 0}, rl.White)
		rl.DrawRectangle(int32(p.GetPos().X)*c.BLOCK_SIZE, int32(p.GetPos().Y)*c.BLOCK_SIZE, c.BLOCK_SIZE, c.BLOCK_SIZE, rl.Red)
		rl.EndMode2D()
		// rl.DrawTextureRec(leftPanel.Texture, rl.Rectangle{X: 0, Y: 0, Width: float32(leftPanel.Texture.Width), Height: float32(-leftPanel.Texture.Height)}, rl.Vector2{X: 0, Y: 0}, rl.White)
		// rl.DrawTextureRec(rightPanel.Texture, rl.Rectangle{X: 0, Y: 0, Width: float32(rightPanel.Texture.Width), Height: float32(-rightPanel.Texture.Height)}, rl.Vector2{X: float32(int32(screenConf.Width) - SIDE_PANEL_WIDTH), Y: 0}, rl.White)

		keyboard.HandleKeyboardMoveInput(p, gameCam, &playerMoved)
		mouse.HandleMouseInputs(gameCam)

		rl.EndDrawing()

	}

}
