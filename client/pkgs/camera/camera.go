package camera

import (
	c "ticoma/client/pkgs/constants"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type DIRECTION uint8

const (
	UP DIRECTION = iota
	DOWN
	LEFT
	RIGHT
)

type GameCamera struct {
	rl.Camera2D
}

func New(targetX float32, targetY float32, offsetX float32, offsetY float32) *GameCamera {
	offset := rl.Vector2{X: offsetX, Y: offsetY}
	target := rl.Vector2{X: targetX, Y: targetY}
	return &GameCamera{
		Camera2D: rl.NewCamera2D(offset, target, 0, 1),
	}
}

func (gc *GameCamera) MoveCamera(dir DIRECTION) {
	if dir == UP {
		gc.Camera2D.Offset.Y += c.BLOCK_SIZE
	}
	if dir == DOWN {
		gc.Camera2D.Offset.Y -= c.BLOCK_SIZE
	}
	if dir == LEFT {
		gc.Camera2D.Offset.X += c.BLOCK_SIZE
	}
	if dir == RIGHT {
		gc.Camera2D.Offset.X -= c.BLOCK_SIZE
	}
}

func (gc *GameCamera) ChangeCameraZoom(val float32) {}
