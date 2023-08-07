package client

import (
	"fmt"
	"os/exec"

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

// get commit id
func getVersion() string {
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		fmt.Println(err)
	}
	commitHash := string(out)
	return commitHash
}

func Main() {
	ver := getVersion()[0:6]
	rl.InitWindow(1280, 720, "raylib [core] example - basic window")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		DrawBg()
		rl.ClearBackground(rl.RayWhite)
		rl.DrawText("ticoma git-"+ver, 2, 3, 20, rl.DarkGray)

		rl.EndDrawing()
	}
}

func DrawBg() {
	for i := 0; i <= MAP_SIZE_IN_BLOCKS; i++ {
		rl.DrawLine(int32(VIEWPORT_START_X+i*BLOCK_SIZE), VIEWPORT_START_Y, int32(VIEWPORT_START_X+i*BLOCK_SIZE), VIEWPORT_SIZE+VIEWPORT_START_Y, rl.Black)
		for j := 0; j <= MAP_SIZE_IN_BLOCKS; j++ {
			rl.DrawLine(VIEWPORT_START_X, int32(VIEWPORT_START_Y+j*BLOCK_SIZE), (VIEWPORT_START_X + VIEWPORT_SIZE), int32(VIEWPORT_START_Y+j*BLOCK_SIZE), rl.Black)
		}
	}
}
