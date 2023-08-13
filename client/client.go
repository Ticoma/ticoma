package client

import (
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
	player "ticoma/internal/packages/player"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WINDOW_WIDTH       = 1280
	WINDOW_HEIGHT      = 720
	BLOCK_SIZE         = 50
	MAP_SIZE_IN_BLOCKS = 10                              // tmp name - how many blocks inside viewport
	VIEWPORT_SIZE      = BLOCK_SIZE * MAP_SIZE_IN_BLOCKS // main in-game window
	VIEWPORT_START_X   = (WINDOW_WIDTH - VIEWPORT_SIZE) / 2
	VIEWPORT_START_Y   = (WINDOW_HEIGHT - VIEWPORT_SIZE) / 2
)

var playerMoved = false

// get commit id
func getVersion() string {
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		fmt.Println(err)
	}
	commitHash := string(out)
	return commitHash
}

func Main(c chan player.Player) {

	// Setup

	// Misc
	ver := getVersion()[0:6]

	// Raylib
	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "Ticoma Client")
	defer rl.CloseWindow()
	rl.SetTraceLog(4) // Disable unnecessary raylib logs
	rl.SetTargetFPS(60)

	// Fonts
	jetbrains := rl.LoadFont("../client/assets/fonts/JetBrainsMono-Medium.ttf")

	// Block execution till we get player instance from internal
	p := <-c

	// init player pos @ random pos
	randX := RandRange(0, 10)
	randY := RandRange(0, 10)
	MovePlayer(p, randX, randY, randX, randY)

	// game loop
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		DrawBg()

		KeyPressHandler(p)

		for id, pos := range *p.GetPlayersPos() {
			DrawPlayer(id, pos.X, pos.Y)
		}

		rl.ClearBackground(rl.RayWhite)

		// info
		leftTop := &rl.Vector2{X: 2, Y: 3}
		leftTop2 := &rl.Vector2{X: 2, Y: 20}
		rl.DrawTextEx(jetbrains, "ticoma git-"+ver, *leftTop, 20, 0, rl.DarkGray)
		rl.DrawTextEx(jetbrains, "peerid-"+p.GetPeerID(), *leftTop2, 20, 0, rl.DarkGray)

		rl.EndDrawing()
	}

}

func KeyPressHandler(p player.Player) {
	if playerMoved {
		return
	}
	pos := p.GetPos()
	x, y := pos.X, pos.Y

	if rl.IsKeyDown(rl.KeyA) {
		MovePlayer(p, x, y, x-1, y)
	}
	if rl.IsKeyDown(rl.KeyD) {
		MovePlayer(p, x, y, x+1, y)
	}
	if rl.IsKeyDown(rl.KeyS) {
		MovePlayer(p, x, y, x, y+1)
	}
	if rl.IsKeyDown(rl.KeyW) {
		MovePlayer(p, x, y, x, y-1)
	}
}

func MovePlayer(p player.Player, posX int, posY int, destX int, destY int) {
	playerMoved = true
	fmt.Println("\nMOVE MOVE MOVE")
	err := p.Move(posX, posY, destX, destY)
	if err != nil {
		fmt.Println("[CLIENT] - Failed to request move, err: ", err)
	}
	time.Sleep(time.Millisecond * 300)
	err = p.Move(destX, destY, destX, destY)
	if err != nil {
		fmt.Println("[CLIENT] - Failed to fulfill move, err: ", err)
	}
	time.Sleep(time.Millisecond * 300)
	playerMoved = false
}

func RandRange(min int, max int) int {
	return rand.Intn(max-min) + min
}

func DrawPlayer(id int, posX int, posY int) {
	rl.DrawRectangle(int32(VIEWPORT_START_X+(BLOCK_SIZE*posX)), int32(VIEWPORT_START_Y+(BLOCK_SIZE*posY)), BLOCK_SIZE, BLOCK_SIZE, rl.Black)
	ids := strconv.Itoa(id)
	rl.DrawText(ids, int32(VIEWPORT_START_X+(BLOCK_SIZE*posX)+5), int32(VIEWPORT_START_Y+(BLOCK_SIZE*posY)+5), 18, rl.Red)
}

func DrawBg() {
	for i := 0; i <= MAP_SIZE_IN_BLOCKS; i++ {
		rl.DrawLine(int32(VIEWPORT_START_X+i*BLOCK_SIZE), VIEWPORT_START_Y, int32(VIEWPORT_START_X+i*BLOCK_SIZE), VIEWPORT_SIZE+VIEWPORT_START_Y, rl.Black)
		for j := 0; j <= MAP_SIZE_IN_BLOCKS; j++ {
			rl.DrawLine(VIEWPORT_START_X, int32(VIEWPORT_START_Y+j*BLOCK_SIZE), (VIEWPORT_START_X + VIEWPORT_SIZE), int32(VIEWPORT_START_Y+j*BLOCK_SIZE), rl.Black)
		}
	}
}
