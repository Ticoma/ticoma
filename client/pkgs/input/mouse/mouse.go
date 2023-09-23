package mouse

import (
	"ticoma/client/pkgs/camera"
	"ticoma/client/pkgs/utils"
	"ticoma/types"

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

	// Check if mouse is on game panel
	mousePos := rl.GetMousePosition()
	isMouseOnGamePanel := rl.CheckCollisionPointRec(mousePos, *targetRec)

	// Draw destination tiles independently
	drawDestTile(cp, cam)

	if isMouseOnGamePanel {
		// Handle mouse input logic for specific panel
		switch targetPanel {
		case GAME:
			gameHandleMouse(cp, cMap, &mousePos, cam)
			gameHandleZoom(cam)
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
// This includes: Hovering over a tile, Left-clicking a tile (move request), Right-clicking a tile (display tile actions, OSRS style)
func gameHandleMouse(cp *player.ClientPlayer, cMap *paths.Grid, mousePos *rl.Vector2, cam *camera.GameCamera) {

	nTileWorld := utils.GetNearestCursorTile(mousePos, cam)
	nTileScreen := rl.GetWorldToScreen2D(rl.Vector2{X: nTileWorld.X, Y: nTileWorld.Y}, cam.Camera2D)
	rl.DrawRectangleLinesEx(rl.Rectangle{X: nTileScreen.X, Y: nTileScreen.Y, Width: c.BLOCK_SIZE * cam.Zoom, Height: c.BLOCK_SIZE * cam.Zoom}, 2, rl.Maroon)

	hoverTileX, hoverTileY := int(nTileWorld.X/c.BLOCK_SIZE), int(nTileWorld.Y/c.BLOCK_SIZE)
	hoverTileWalkable := gamemap.IsTileWalkable(cMap, hoverTileX, hoverTileY)
	cp.HoverTile = &types.Position{
		X: hoverTileX + 1, Y: hoverTileY + 1, // +1 = Pos offset
	}

	if hoverTileWalkable {

		rl.DrawRectangleLinesEx(rl.Rectangle{X: nTileScreen.X, Y: nTileScreen.Y, Width: c.BLOCK_SIZE * cam.Zoom, Height: c.BLOCK_SIZE * cam.Zoom}, 2, rl.Green)

		// Click tile to set it to active destination and Move towards it
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			cp.DestTile = &types.Position{
				X: hoverTileX + 1, Y: hoverTileY + 1, // +1 = Pos offset
			}
			cp.MouseMove(cMap)
		}
	}

}

func drawDestTile(cp *player.ClientPlayer, cam *camera.GameCamera) {
	pPos := cp.InternalPlayer.GetPos().Position
	if cp.DestTile != nil {
		if cp.DestTile.X != pPos.X || cp.DestTile.Y != pPos.Y {
			activeTileScreen := utils.TileToScreenPos(cp.DestTile.X-1, cp.DestTile.Y-1, cam)
			rl.DrawRectangleLinesEx(rl.Rectangle{X: activeTileScreen.X, Y: activeTileScreen.Y, Width: c.BLOCK_SIZE * cam.Zoom, Height: c.BLOCK_SIZE * cam.Zoom}, 2, c.COLOR_PANEL_OUTLINE)
		}
	}
}
