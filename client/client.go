package client

import (
	"fmt"
	"ticoma/client/pkgs/camera"
	c "ticoma/client/pkgs/constants"
	dr "ticoma/client/pkgs/drawing"
	"ticoma/client/pkgs/drawing/panels/left"
	"ticoma/client/pkgs/input/keyboard"
	"ticoma/client/pkgs/input/mouse"
	"ticoma/client/pkgs/player"

	"ticoma/client/pkgs/utils"
	internal_player "ticoma/internal/packages/player"
	"ticoma/types"

	"github.com/fstanis/screenresolution"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var chatMsgs []string
var chatInput []byte
var hold int
var mousePos rl.Vector2

func Main(pc chan internal_player.Player, cc chan types.ChatMessage, fullscreen *bool) {

	// Load icon
	icon := rl.LoadImage("../client/assets/logo/ticoma-logo-64.png")

	// Load screen conf
	res := screenresolution.GetPrimary()
	screenC := utils.GetScreenConf(res.Width, res.Height, fullscreen)

	// Setup window, resolution, raylib
	rl.InitWindow(screenC.Width, screenC.Height, "Ticoma Client")
	if *fullscreen {
		rl.ToggleFullscreen()
	}

	defer rl.CloseWindow()
	rl.SetTraceLog(4) // Disable unnecessary raylib logs
	rl.SetTargetFPS(60)
	rl.SetWindowIcon(*icon)

	// Load fonts, imgs
	font := rl.LoadFontEx("../client/assets/fonts/ponderosa.regular.ttf", c.DEFAULT_FONT_SIZE*4, nil)
	spawnImg := rl.LoadImage("../client/assets/textures/map/spawn.png")

	// Setup res, scaling
	SIDE_PANEL_WIDTH := int32((screenC.Width / 4))
	// SCALING := screenC.Height / 1080

	// Setup game
	gameCam := camera.New(float32(spawnImg.Width/2), float32(spawnImg.Height/2), float32(screenC.Width/2), float32(screenC.Height/2))
	playerMoved := false
	spawnMapSize := 25 // In blocks (tmp)

	// Setup textures, panels
	world := rl.LoadRenderTexture(spawnImg.Width, spawnImg.Height) // Full sized map
	// rightPanel := rl.LoadRenderTexture(SIDE_PANEL_WIDTH, int32(screenC.Height)) // Side panels

	// tmp, Draw map on world from texture
	spawnTxt := rl.LoadTextureFromImage(spawnImg)

	// Wait for player conn
	p := <-pc

	// Activate chat listener
	go ChatMsgListener(p, cc)

	// Init player
	player.InitPlayer(p, &playerMoved, spawnMapSize/2, spawnMapSize/2)

	// Init side panels
	sidePanelColor := rl.DarkGray
	leftPanelTabs := map[int][2]string{
		0: {"Chat", "C"},
		1: {"Test", "T"},
		2: {"Tabssss bro", "Tb"},
	}
	leftPanelRt2D := rl.LoadRenderTexture(SIDE_PANEL_WIDTH, screenC.Height)
	leftPanel := left.New(&leftPanelRt2D, float32(SIDE_PANEL_WIDTH), float32(screenC.Height), 0, 0, &sidePanelColor, leftPanelTabs)

	for !rl.WindowShouldClose() {

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		mousePos = rl.GetMousePosition()

		// Draw players
		dr.DrawMap(&world, &spawnTxt, gameCam.Zoom)
		dr.DrawPlayers(&world, p)

		// Draw game
		rl.BeginMode2D(gameCam.Camera2D)
		rl.DrawTextureRec(world.Texture, rl.Rectangle{X: 0, Y: 0, Width: float32(world.Texture.Width), Height: float32(-world.Texture.Height)}, rl.Vector2{X: 0, Y: 0}, rl.White)
		rl.DrawRectangleRec(rl.Rectangle{X: float32(p.GetPos().X * c.BLOCK_SIZE), Y: float32(p.GetPos().Y * c.BLOCK_SIZE), Width: c.BLOCK_SIZE, Height: c.BLOCK_SIZE}, rl.Black)
		rl.EndMode2D()

		// Render panels
		leftPanel.DrawContent(&font)
		leftPanel.RenderPanel()

		// Handle inputs
		chatInput, hold = keyboard.HandleChatInput(chatInput, hold)
		mouse.HandleMouseInputs(gameCam)

		rl.EndDrawing()
	}

}

func ChatMsgListener(p internal_player.Player, cc chan types.ChatMessage) {
	for {
		chat := <-cc
		msg := fmt.Sprintf("[player %d]: %s", chat.PlayerId, chat.Message)
		chatMsgs = append(chatMsgs, msg)
		if p.GetId() == chat.PlayerId { // if the msg arrived from us - clear input
			chatInput = nil
		}
	}
}
