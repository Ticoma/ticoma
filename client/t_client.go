package client

import (
	"fmt"
	"math/rand"
	"os/exec"
	player "ticoma/internal/packages/player"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type PlayerClient struct {
	X int
	Y int
}

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

func Main(c chan player.PlayerInterface) {

	// setup
	ver := getVersion()[0:6]
	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "Ticoma Client")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)
	jetbrains := rl.LoadFont("client/assets/fonts/JetBrainsMono-Medium.ttf")

	// retrieve initialized and connected player from internal
	p := <-c

	// init player pos @ random pos
	randX := RandRange(0, 10)
	randY := RandRange(0, 10)
	pc := &PlayerClient{}
	MovePlayer(p, pc, randX, randY, randX, randY)
	go KeyPressHandler(p, pc)

	// game loop
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		DrawBg()

		if playerMoved {
			DrawPlayer(pc.X, pc.Y)
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

func KeyPressHandler(p player.PlayerInterface, pc *PlayerClient) {
	for {
		if rl.IsKeyDown(rl.KeyA) {
			MovePlayer(p, pc, pc.X, pc.Y, pc.X-1, pc.Y)
		}
		if rl.IsKeyDown(rl.KeyD) {
			MovePlayer(p, pc, pc.X, pc.Y, pc.X+1, pc.Y)
		}
		if rl.IsKeyDown(rl.KeyS) {
			MovePlayer(p, pc, pc.X, pc.Y, pc.X, pc.Y+1)
		}
		if rl.IsKeyDown(rl.KeyW) {
			MovePlayer(p, pc, pc.X, pc.Y, pc.X, pc.Y-1)
		}
	}
}

// {1, 1, 1, 1}
// {1, 1, 2, 2}
// {2, 2, 2, 2}

func MovePlayer(p player.PlayerInterface, pc *PlayerClient, posX int, posY int, destX int, destY int) {
	playerMoved = true
	err := p.Move(posX, posY, destX, destY)
	if err != nil {
		fmt.Println("[CLIENT] - Failed to move, err: ", err)
	}
	HandleKeyCooldown(p, pc, destX, destY, destX, destY)
}

func HandleKeyCooldown(p player.PlayerInterface, pc *PlayerClient, posX int, posY int, destX int, destY int) {
	time.Sleep(time.Millisecond * 250)
	err := p.Move(destX, destY, destX, destY)
	pc.X, pc.Y = destX, destY
	if err != nil {
		fmt.Println("[CLIENT][HANDLE KEY CD] - Failed to move, err: ", err)
	}
	fmt.Println("CAN PRESS AGAIN")
}

func RandRange(min int, max int) int {
	return rand.Intn(max-min) + min
}

// Tmp solution
func DrawPlayer(posX int, posY int) {
	rl.DrawRectangle(int32(VIEWPORT_START_X+(BLOCK_SIZE*posX)), int32(VIEWPORT_START_Y+(BLOCK_SIZE*posY)), BLOCK_SIZE, BLOCK_SIZE, rl.Black)
}

func DrawBg() {
	for i := 0; i <= MAP_SIZE_IN_BLOCKS; i++ {
		rl.DrawLine(int32(VIEWPORT_START_X+i*BLOCK_SIZE), VIEWPORT_START_Y, int32(VIEWPORT_START_X+i*BLOCK_SIZE), VIEWPORT_SIZE+VIEWPORT_START_Y, rl.Black)
		for j := 0; j <= MAP_SIZE_IN_BLOCKS; j++ {
			rl.DrawLine(VIEWPORT_START_X, int32(VIEWPORT_START_Y+j*BLOCK_SIZE), (VIEWPORT_START_X + VIEWPORT_SIZE), int32(VIEWPORT_START_Y+j*BLOCK_SIZE), rl.Black)
		}
	}
}
