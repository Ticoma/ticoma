package mouse

import (
	"fmt"
	"ticoma/client/pkgs/camera"
	"ticoma/client/pkgs/utils"

	"ticoma/client/pkgs/player"

	c "ticoma/client/pkgs/constants"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type PANEL_TYPE int

const (
	GAME PANEL_TYPE = iota
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
		// Handle mouse input logic for specific panel
		switch targetPanel {
		case GAME:
			gameHandleZoom(cam)
			gameHandleMouse(cp, &mousePos, cam)
		case LEFT:
		case RIGHT:
		default:
		}
	}
}

// Input handler for zooming the game view
func gameHandleZoom(cam *camera.GameCamera) {
	scroll := rl.GetMouseWheelMoveV().Y
	if scroll != 0 {
		if scroll > 0 {
			camera.HandleScrollZoom(camera.UP, cam)
		} else {
			camera.HandleScrollZoom(camera.DOWN, cam)
		}
	}
}

// Input handler for the mouse
//
// This includes:
// Hovering over a tile, Left-clicking a tile (move request), Right-clicking a tile (display tile actions, OSRS style)
func gameHandleMouse(cp *player.ClientPlayer, mousePos *rl.Vector2, cam *camera.GameCamera) {

	nTileWorld := utils.GetNearestCursorTile(mousePos, cam)
	nTileScreen := rl.GetWorldToScreen2D(rl.Vector2{X: nTileWorld.X, Y: nTileWorld.Y}, cam.Camera2D)
	rl.DrawRectangleLinesEx(rl.Rectangle{X: nTileScreen.X, Y: nTileScreen.Y, Width: c.BLOCK_SIZE * cam.Zoom, Height: c.BLOCK_SIZE * cam.Zoom}, 2, rl.Magenta)

	// Draw current active tile
	if cp.IsActiveTile {
		activeTileScreen := utils.TileToScreenPos(cp.ActiveTile.X, cp.ActiveTile.Y, cam)
		fmt.Println(activeTileScreen)
		activeTileColor := rl.Yellow
		rl.DrawRectangleLinesEx(rl.Rectangle{X: activeTileScreen.X, Y: activeTileScreen.Y, Width: c.BLOCK_SIZE * cam.Zoom, Height: c.BLOCK_SIZE * cam.Zoom}, 2, activeTileColor)
	}

	// Click tile to set it to active destination
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		tileNumX, tileNumY := (nTileWorld.X/c.BLOCK_SIZE)+1, (nTileWorld.Y/c.BLOCK_SIZE)+1
		cp.ActiveTile.X, cp.ActiveTile.Y = int(tileNumX), int(tileNumY)
		cp.IsActiveTile = true
		// fmt.Println("ACTIVE TILE: ", cp.ActiveTile.X, cp.ActiveTile.Y)

		// TODO: Add path finding, mv request loop to dest
	}
}
