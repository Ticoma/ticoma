package utils

import (
	"fmt"
	"math"
	"math/rand"
	"os/exec"
	"ticoma/client/pkgs/camera"
	c "ticoma/client/pkgs/constants"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Get current hash as string
func GetCommitHash() string {
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		fmt.Println(err)
	}
	commitHash := string(out)
	return commitHash
}

// Configure game window resolution based on screen, flags
func ConfLaunchGameRes(width int, height int, fullscreen *bool) *c.ScreenInfo {
	screenConf := &c.ScreenInfo{}
	if *fullscreen { // if not fullscreen, make a quarter window
		screenConf.Width = int32(width)
		screenConf.Height = int32(height)
	} else {
		screenConf.Width = int32(width / 2)
		screenConf.Height = int32(height / 2)
	}

	// Ignore this for now
	screenConf.RefreshRate = 60
	// TODO: Add scale factor, UI scaling for different screens
	return screenConf
}

// Get substring of first N chars in string (if N > len, returns original string)
func FirstN(s string, n int) string {
	i := 0
	for j := range s {
		if i == n {
			return s[:j]
		}
		i++
	}
	return s
}

// Gen random number in range - (inclusive, exclusive)
func RandRange(min int, max int) int {
	return rand.Intn(max-min) + min
}

// Round float to specified unit - (upwards i think?)
func FloatRound(x, unit float32) float32 {
	return float32(math.Round(float64(x/unit)) * float64(unit))
}

// Returns the nearest to mouse cursor tile position (in px, world)
func GetNearestCursorTile(mousePos *rl.Vector2, cam *camera.GameCamera) *rl.Rectangle {
	worldMousePos := rl.GetScreenToWorld2D(*mousePos, cam.Camera2D)
	nearestTile := &rl.Rectangle{
		X:      FloatRound(worldMousePos.X-c.BLOCK_SIZE/2, c.BLOCK_SIZE),
		Y:      FloatRound(worldMousePos.Y-c.BLOCK_SIZE/2, c.BLOCK_SIZE),
		Width:  c.BLOCK_SIZE,
		Height: c.BLOCK_SIZE,
	}
	return nearestTile
}

// Convert tile position, e.g {5, 5} to top-left pixel pos on Screen
func TileToScreenPos(tileX int, tileY int, cam *camera.GameCamera) *rl.Vector2 {
	tWorldPos := rl.Vector2{X: float32(tileX) * c.BLOCK_SIZE, Y: float32(tileY) * c.BLOCK_SIZE}
	tScreenPos := rl.GetWorldToScreen2D(tWorldPos, cam.Camera2D)
	return &tScreenPos
}

// Convert Tile screen position to rl.Rectangle ready to render on screen
func TileScreenToRec(tileScreenPos *rl.Vector2, cam *camera.GameCamera) *rl.Rectangle {
	return &rl.Rectangle{
		X:      tileScreenPos.X,
		Y:      tileScreenPos.Y,
		Width:  c.BLOCK_SIZE * cam.Zoom,
		Height: c.BLOCK_SIZE * cam.Zoom,
	}
}
