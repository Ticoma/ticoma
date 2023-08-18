package camera

import (
	c "ticoma/client/packages/constants"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type DIRECTION int

const (
	UP DIRECTION = iota
	DOWN
	LEFT
	RIGHT
)

type GameCamera struct {
	rl.Camera2D
}

func New() *GameCamera {
	offset := rl.Vector2{X: 0, Y: 0}
	target := rl.Vector2{X: 0, Y: 0}
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
