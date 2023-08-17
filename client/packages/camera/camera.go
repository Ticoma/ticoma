package camera

import (
	"fmt"
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
		fmt.Println("UP")
		gc.Camera2D.Offset.Y += c.BLOCK_SIZE
	}
	if dir == DOWN {
		fmt.Println("DOWN")
		gc.Camera2D.Offset.Y -= c.BLOCK_SIZE
	}
	if dir == LEFT {
		fmt.Println("LEFT")
		gc.Camera2D.Offset.X += c.BLOCK_SIZE
	}
	if dir == RIGHT {
		fmt.Println("RIGHT")
		gc.Camera2D.Offset.X -= c.BLOCK_SIZE
	}
}

func (gc *GameCamera) ChangeCameraZoom(val float32) {}
