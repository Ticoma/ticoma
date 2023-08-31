package keyboard

import (
	"ticoma/client/packages/camera"
	"ticoma/client/packages/player"
	internal_player "ticoma/internal/pkgs/player"

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

// Handle key presses when chat is active (appends to value of provided input byte arr)
func HandleChatInput(input []byte) []byte {

	key := rl.GetCharPressed()

	if len(input) == 0 && rl.IsKeyPressed(rl.KeySpace) {
		return input
	}

	if key != 0 {
		if key >= 32 && key <= 125 {
			input = append(input, byte(key))
		}
	}

	if rl.IsKeyPressed(rl.KeyBackspace) {
		if len(input) > 0 {
			input = input[:len(input)-1] // remove last char
		}
	}

	return input
}
