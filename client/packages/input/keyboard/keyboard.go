package keyboard

import (
	"ticoma/client/packages/camera"
	"ticoma/client/packages/player"
	internal_player "ticoma/internal/packages/player"

	c "ticoma/client/packages/constants"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func HandleKeyboardMoveInput(p internal_player.Player, cam *camera.GameCamera, playerMoveState *bool) {
	if *playerMoveState {
		return
	}
	pos := p.GetPos()
	x, y := pos.X, pos.Y

	if rl.IsKeyDown(rl.KeyA) {
		player.MovePlayer(p, cam, camera.LEFT, playerMoveState, x, y, x-1, y)
		cam.Target.X -= c.BLOCK_SIZE
	}
	if rl.IsKeyDown(rl.KeyD) {
		player.MovePlayer(p, cam, camera.RIGHT, playerMoveState, x, y, x+1, y)
		cam.Target.X += c.BLOCK_SIZE
	}
	if rl.IsKeyDown(rl.KeyS) {
		player.MovePlayer(p, cam, camera.DOWN, playerMoveState, x, y, x, y+1)
		cam.Target.Y += c.BLOCK_SIZE
	}
	if rl.IsKeyDown(rl.KeyW) {
		player.MovePlayer(p, cam, camera.UP, playerMoveState, x, y, x, y-1)
		cam.Target.Y -= c.BLOCK_SIZE
	}
}
