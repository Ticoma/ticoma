package client

import (
	"fmt"
	c "ticoma/client/pkgs/constants"
	scene_handler "ticoma/client/pkgs/drawing/scenes"
	"ticoma/client/pkgs/player"
	"ticoma/client/pkgs/utils"
	internal_player "ticoma/internal/pkgs/player"

	"github.com/fstanis/screenresolution"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var icon *rl.Image
var programState *scene_handler.ProgramState

func Main(pc chan internal_player.Player, rc chan interface{}, fullscreen *bool) {

	// Load icon
	icon = rl.LoadImage("../client/assets/logo/ticoma-logo-64.png")

	// Load screen conf
	res := screenresolution.GetPrimary()
	c.SCREEN = utils.ConfLaunchGameRes(res.Width, res.Height, fullscreen)

	// Setup window, resolution, raylib
	rl.InitWindow(c.SCREEN.Width, c.SCREEN.Height, "Ticoma Client")
	if *fullscreen {
		rl.ToggleFullscreen()
	}

	defer rl.CloseWindow()
	rl.SetTraceLogLevel(rl.LogError)
	rl.SetTargetFPS(60)
	rl.SetWindowIcon(*icon)

	// Game state / scene handler
	programState = scene_handler.New()
	programState.Running = true

	// Load fonts, imgs
	c.DEFAULT_FONT = rl.LoadFontEx("../client/assets/fonts/clacon2.ttf", int32(c.DEFAULT_FONT_SIZE)*6, nil)
	rl.SetTextureFilter(c.DEFAULT_FONT.Texture, rl.TextureFilterMode(rl.RL_TEXTURE_FILTER_BILINEAR))

	// Wait until internal connects with network
	p := <-pc
	// Init client side player instance
	cp := player.New(&p)

	for programState.Running {

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		programState.HandleScene(cp)
		updateState()

		rl.EndDrawing()
	}

	exit()
}

func updateState() {
	programState.Running = !rl.WindowShouldClose()
}

// Cleanup before exiting
func exit() {
	rl.CloseWindow()
	rl.UnloadFont(c.DEFAULT_FONT)
	rl.UnloadImage(icon)
}

func requestListener(rq *chan interface{}) {
	for {
		req := <-*rq
		fmt.Printf("Client received request: %v\n", req)
	}
}
