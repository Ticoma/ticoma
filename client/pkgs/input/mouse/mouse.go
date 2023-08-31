package mouse

import (
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
			gameHandleMouse(cp, &mousePos, targetRec, cam)
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
// Hovering over a tile, Left-clicking a tile (move request), Right-clicking a tile (actions)
func gameHandleMouse(cp *player.ClientPlayer, mousePos *rl.Vector2, target *rl.Rectangle, cam *camera.GameCamera) {

	nTileWorld := utils.GetNearestCursorTile(mousePos, cam)
	nTileScreen := rl.GetWorldToScreen2D(rl.Vector2{X: nTileWorld.X, Y: nTileWorld.Y}, cam.Camera2D)
	rl.DrawRectangleLinesEx(rl.Rectangle{X: nTileScreen.X, Y: nTileScreen.Y, Width: c.BLOCK_SIZE * cam.Zoom, Height: c.BLOCK_SIZE * cam.Zoom}, 2, rl.Magenta)

	// if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
	// 	fmt.Println("MOUSE ", bX, bY)
	// 	fmt.Println("TILE : ", tileX, tileY)
	// 	cp.ActiveTile = &interfaces.Tile{Vector2: rl.Vector2{X: float32(tileX), Y: float32(tileY)}}
	// 	fmt.Println(cp.ActiveTile)
	// 	cp.IsActiveTile = true
	// }
}
