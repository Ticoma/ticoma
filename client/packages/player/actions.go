package player

import (
	"fmt"
	"ticoma/client/packages/camera"
	internal_player "ticoma/internal/pkgs/player"
	"time"
)

//
// Player package
// High-level internal_player interface wrapper for Client
//
// E.g. Move Action calls underlying Player interface (Move) and Updates
// the client-side player position
//

// Like move, but single pkg -> it should fill both places in empty cache
func InitPlayer(p internal_player.Player, moveState *bool, posX int, posY int) {
	*moveState = true
	err := p.Move(posX, posY, posX, posY)
	if err != nil {
		fmt.Println("[CLIENT] - Failed to init player! err: ", err)
	}
	time.Sleep(time.Millisecond * 300)
	*moveState = false
}

// Move with engine-safe delay between pkgs
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

// Send a message to chat
func Chat(p internal_player.Player, msg []byte) bool {
	err := p.Chat(msg)
	if err != nil {
		fmt.Println("[CLIENT] - Failed to send chat msg, err: ", err)
		return false
	}
	return true
}
