package actions

import (
	"fmt"
	"ticoma/client/packages/camera"
	internal_player "ticoma/internal/packages/player"
	"time"
)

//
// Actions package
// High-level commands for performing client-side operations
//
// E.g. Move Action calls underlying Player interface (Move) and Updates
// the client-side player position
//

func MovePlayer(p internal_player.Player, cam *camera.GameCamera, dir camera.DIRECTION, moveState *bool, posX int, posY int, destX int, destY int) {
	*moveState = true
	err := p.Move(posX, posY, destX, destY)
	if err != nil {
		fmt.Println("[CLIENT] - Failed to move (1), err: ", err)
	}
	time.Sleep(time.Millisecond * 300)
	err = p.Move(destX, destY, destX, destY)
	if err != nil {
		fmt.Println("[CLIENT] - Failed to move (2), err: ", err)
	}
	time.Sleep(time.Millisecond * 300)
	*moveState = false
}

func InitPlayer(p internal_player.Player, moveState *bool, posX int, posY int) {
	*moveState = true
	err := p.Move(posX, posY, posX, posY)
	if err != nil {
		fmt.Println("[CLIENT] - Failed to init player! err: ", err)
	}
	time.Sleep(time.Millisecond * 300)
	*moveState = false
}
