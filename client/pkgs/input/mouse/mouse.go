package mouse

import (
	"fmt"
	"ticoma/client/pkgs/camera"
	"ticoma/client/pkgs/utils"

	"ticoma/client/pkgs/player"

	c "ticoma/client/pkgs/constants"
	gamemap "ticoma/client/pkgs/game_map"

	"github.com/solarlune/paths"

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
func HandleMouseInputs(cp *player.ClientPlayer, cMap *paths.Grid, cam *camera.GameCamera, targetRec *rl.Rectangle, targetPanel PANEL_TYPE) {

	// Check if mouse is on panel
	mousePos := rl.GetMousePosition()
	isMouseOnGamePanel := rl.CheckCollisionPointRec(mousePos, *targetRec)

	if isMouseOnGamePanel {
		// Handle mouse input logic for specific panel
		switch targetPanel {
		case GAME:
			gameHandleZoom(cam)
			gameHandleMouse(cp, cMap, &mousePos, cam)
		case LEFT:
		case RIGHT:
		default:
		}
	} else {
		cp.IsHoveringTile = false
		return
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
// This includes: Hovering over a tile, Left-clicking a tile (move request), Right-clicking a tile (display tile actions, OSRS style)
func gameHandleMouse(cp *player.ClientPlayer, cMap *paths.Grid, mousePos *rl.Vector2, cam *camera.GameCamera) {

	nTileWorld := utils.GetNearestCursorTile(mousePos, cam)
	nTileScreen := rl.GetWorldToScreen2D(rl.Vector2{X: nTileWorld.X, Y: nTileWorld.Y}, cam.Camera2D)
	rl.DrawRectangleLinesEx(rl.Rectangle{X: nTileScreen.X, Y: nTileScreen.Y, Width: c.BLOCK_SIZE * cam.Zoom, Height: c.BLOCK_SIZE * cam.Zoom}, 2, rl.Magenta)

	hoverTile := rl.Vector2{X: (nTileWorld.X / c.BLOCK_SIZE), Y: (nTileWorld.Y / c.BLOCK_SIZE)}
	hoverTileWalkable := gamemap.IsTileWalkable(cMap, int(hoverTile.X), int(hoverTile.Y))
	// fmt.Println("HOVER TILE: ", hoverTile.X, hoverTile.Y)

	if hoverTileWalkable {

		// Calculate path
		path := cMap.GetPath(float64(cp.InternalPlayer.GetPos().Position.X-1), float64(cp.InternalPlayer.GetPos().Position.Y-1), float64(hoverTile.X), float64(hoverTile.Y), false, false)
		if path != nil {
			fmt.Println("PATH START ===")
			for i, tile := range path.Cells {
				fmt.Println("MOVE ", i, ": ", tile.X, tile.Y)
				tileScreenPos := utils.TileToScreenPos(tile.X, tile.Y, cam)
				tileRec := utils.TileScreenToRec(tileScreenPos, cam)
				rl.DrawRectangleLinesEx(*tileRec, 2, rl.Yellow)
			}
			fmt.Println("PATH END ===")
		}

		// Click tile to set it to active destination
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			cp.ActiveTile.X, cp.ActiveTile.Y = int(hoverTile.X), int(hoverTile.Y)
			cp.IsActiveTile = true
			// fmt.Println("ACTIVE TILE: ", cp.ActiveTile.X, cp.ActiveTile.Y)
		}
	}

	// Draw current active tile (destination)
	if cp.IsActiveTile {
		activeTileScreen := utils.TileToScreenPos(cp.ActiveTile.X, cp.ActiveTile.Y, cam)
		rl.DrawRectangleLinesEx(rl.Rectangle{X: activeTileScreen.X, Y: activeTileScreen.Y, Width: c.BLOCK_SIZE * cam.Zoom, Height: c.BLOCK_SIZE * cam.Zoom}, 2, c.COLOR_PANEL_OUTLINE)
	}
}
