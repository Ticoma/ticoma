package player

import (
	"fmt"
	"ticoma/client/pkgs/camera"
	"ticoma/client/pkgs/interfaces"
	internal_player "ticoma/internal/pkgs/player"
	"time"
)

//
// Player actions package
// High-level internal_player interface wrapper for Client
//
// E.g. Move Action calls underlying Player interface (Move request) and Updates
// the client-side player position on success
//

type ClientPlayer struct {
	IsOnline       bool
	InternalPlayer internal_player.Player // Interface passed from internal
	IsActiveTile   bool                   // Helpers - hoverTile, activeTile can't be nil, and 0,0 is a valid tile
	IsHoveringTile bool
	HoverTile      *interfaces.Tile // Currently hovered game tile (empty if none)
	ActiveTile     *interfaces.Tile // Clicked tile (destination of most recent move request)
	IsMoving       bool             // Input blocker
}

func New(p *internal_player.Player) *ClientPlayer {
	return &ClientPlayer{
		IsOnline:       false,
		InternalPlayer: *p,
		IsActiveTile:   false,
		IsHoveringTile: false,
		HoverTile:      &interfaces.Tile{},
		ActiveTile:     &interfaces.Tile{},
		IsMoving:       false,
	}
}

// Send a REGISTER_ request from Client
func (cp *ClientPlayer) Register() {
	err := cp.InternalPlayer.Register()
	if err != nil {
		fmt.Println("[CLIENT] - Failed to register. Err: " + err.Error())
		return
	}
	// Change state if OK
	cp.IsOnline = true
}

// Send a LOGIN_ request from Client
func (cp *ClientPlayer) Login() {
	err := cp.InternalPlayer.Login()
	if err != nil {
		fmt.Println("[CLIENT] - Failed to login. Err: " + err.Error())
		return
	}
	cp.IsOnline = true
}

// Send two MOVE_ requests with engine-safe cooldown
func (cp *ClientPlayer) MovePlayer(cam *camera.GameCamera, destX int, destY int) {
	cp.IsMoving = true
	currPos := cp.InternalPlayer.GetPos().Position
	// Move request (1)
	err := cp.InternalPlayer.Move(&currPos.X, &currPos.Y, &destX, &destY)
	if err != nil {
		fmt.Println("[CLIENT] - Failed to move (1), err: ", err)
	}
	time.Sleep(time.Millisecond * 300)
	// Move arrive request (2)
	err = cp.InternalPlayer.Move(&destX, &destY, &destX, &destY)
	if err != nil {
		fmt.Println("[CLIENT] - Failed to move (2), err: ", err)
	}
	time.Sleep(time.Millisecond * 300)
	cp.IsMoving = false
}

// Send a message to chat
func (cp *ClientPlayer) Chat(msg *[]byte) error {
	err := cp.InternalPlayer.Chat(msg)
	if err != nil {
		return fmt.Errorf("[CLIENT] - Chat Err: %v", err)
	}
	return nil
}
