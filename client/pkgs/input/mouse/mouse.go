package mouse

import (
	"fmt"
	"math"
	"ticoma/client/pkgs/camera"
	"ticoma/client/pkgs/interfaces"
	"ticoma/client/pkgs/utils"

	c "ticoma/client/pkgs/constants"
	"ticoma/client/pkgs/player"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type PANEL_TYPE int

const (
	CENTER PANEL_TYPE = iota
	LEFT
	RIGHT
)

// Checks if user is currently hovering over rect
func IsMouseHoveringRec(objRec *rl.Rectangle) bool {
	return rl.CheckCollisionPointRec(rl.GetMousePosition(), *objRec)
}

// Handles mouse inputs on specific panel / rect inside panel
func HandleMouseInputs(cp *player.ClientPlayer, cam *camera.GameCamera, targetRec *rl.Rectangle, targetPanel PANEL_TYPE) {

	// Check if mouse is on panel
	mousePos := rl.GetMousePosition()
	isMouseOnGamePanel := rl.CheckCollisionPointRec(mousePos, *targetRec)

	if isMouseOnGamePanel {

		// empty := [2]int {0, 0}
		// Handle mouse input logic for specific panel
		switch targetPanel {
		case CENTER:
			centerHandleZoom(cam)
			centerHandleClick(cp, &mousePos, targetRec, cam)
		case LEFT:
		case RIGHT:
		default:
		}
	}
}

// Input handler for zooming the game view
func centerHandleZoom(cam *camera.GameCamera) {
	scroll := rl.GetMouseWheelMoveV().Y
	if scroll != 0 {
		if scroll > 0 {
			camera.HandleScrollZoom(camera.UP, cam)
		} else {
			camera.HandleScrollZoom(camera.DOWN, cam)
		}
	}
}

// Input handler for all tile-related stuff
func centerHandleClick(cp *player.ClientPlayer, mousePos *rl.Vector2, target *rl.Rectangle, cam *camera.GameCamera) {

	gameViewCenterX, gameViewCenterY := target.Width/2-c.BLOCK_SIZE/2, target.Height/2-c.BLOCK_SIZE/2 // center of visible game viewport
	centerPos := cp.InternalPlayer.GetPos()                                                           // player is always in the center, so we'll use that as a reference point
	mousePosGameX, mousePosGameY := mousePos.X-target.X, mousePos.Y                                   // mousePos but with offset calculated from left-top corner of game viewport
	blockX, blockY := math.Ceil(float64((mousePosGameX-gameViewCenterX)/c.BLOCK_SIZE))-1, math.Ceil(float64((mousePosGameY-gameViewCenterY)/c.BLOCK_SIZE))-1
	tileX, tileY := centerPos.X+int(blockX), centerPos.Y+int(blockY)

	bX := utils.FloatRound(mousePos.X, 64)
	bY := utils.FloatRound(mousePosGameY-c.BLOCK_SIZE/2, 64)

	// If there's an active tile, draw an outline around it
	// if cp.IsActiveTile {
	// 	xTileDiff := math.Abs(float64(centerPos.X) - float64(cp.ActiveTile.X))
	// 	yTileDiff := float64(centerPos.Y) - float64(cp.ActiveTile.Y)
	// 	fmt.Println(xTileDiff, yTileDiff)
	// 	activeTile := rl.Rectangle{
	// 		X:      gameViewCenterX + (c.BLOCK_SIZE * float32(xTileDiff)),
	// 		Y:      gameViewCenterY - (c.BLOCK_SIZE * float32(yTileDiff)),
	// 		Width:  c.BLOCK_SIZE,
	// 		Height: c.BLOCK_SIZE,
	// 	}
	// 	rl.DrawRectangleLinesEx(activeTile, 2, rl.Maroon)
	// }

	cursorSnapTile := rl.Rectangle{
		X:      bX - c.BLOCK_SIZE/2,
		Y:      bY, // Weird offset, check math later
		Width:  c.BLOCK_SIZE,
		Height: c.BLOCK_SIZE,
	}
	rl.DrawRectangleLinesEx(cursorSnapTile, 2, rl.Magenta)

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		fmt.Println("MOUSE ", bX, bY)
		fmt.Println("TILE : ", tileX, tileY)
		cp.ActiveTile = &interfaces.Tile{Vector2: rl.Vector2{X: float32(tileX), Y: float32(tileY)}}
		fmt.Println(cp.ActiveTile)
		cp.IsActiveTile = true
	}
}
