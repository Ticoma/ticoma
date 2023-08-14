package client

import (
	"ticoma/client/packages/actions"
	c "ticoma/client/packages/constants"
	gamemap "ticoma/client/packages/drawing/game_map"
	"ticoma/client/packages/drawing/player"
	"ticoma/client/packages/input/keyboard"
	"ticoma/client/packages/utils"
	internal_player "ticoma/internal/packages/player"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var playerMoved = false

func Main(pc chan internal_player.Player) {

	// Load assets
	// font := rl.LoadFont("../client/assets/fonts/Consolas.ttf") // TODO: Fix font
	icon := rl.LoadImage("../client/assets/logo/ticoma-logo-64.png")

	// Setup window, raylib
	rl.InitWindow(c.WINDOW_WIDTH, c.WINDOW_HEIGHT, "Ticoma Client")
	defer rl.CloseWindow()
	rl.SetTraceLog(4) // Disable unnecessary raylib logs
	rl.SetTargetFPS(60)
	rl.SetWindowIcon(*icon)

	// Display loading scene, wait for internal
	p := <-pc

	// init player pos @ random pos
	randX := utils.RandRange(0, 10)
	randY := utils.RandRange(0, 10)
	actions.MovePlayer(p, &playerMoved, randX, randY, randX, randY)

	// Misc
	ver := utils.GetCommitHash()[0:6]

	// game loop
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		// Draw main map
		gamemap.DrawGameMap(c.GAME_MAP_START_X, c.GAME_MAP_START_Y, c.GAME_MAP_BLOCK_SIZE, c.GAME_MAP_SIZE_IN_BLOCKS)

		keyboard.HandleKeyboardMoveInput(p, &playerMoved)

		for id, pos := range *p.GetPlayersPos() {
			player.DrawPlayer(id, pos.X, pos.Y)
		}

		rl.ClearBackground(rl.RayWhite)

		// Draw top-left info
		rl.DrawText("ticoma git-"+ver, 3, 0, 24, rl.Black)
		rl.DrawText("peerid-"+p.GetPeerID(), 3, 30, 24, rl.Black)

		rl.EndDrawing()
	}
}
