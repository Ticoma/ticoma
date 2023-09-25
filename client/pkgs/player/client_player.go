package player

import (
	"fmt"
	c "ticoma/client/pkgs/constants"
	internal_player "ticoma/internal/pkgs/player"
	"ticoma/types"
	"time"

	"github.com/solarlune/paths"
)

//
// Player actions package
// High-level internal_player interface wrapper for Client
//
// E.g. Move Action calls underlying Player interface (Move request) and Updates
// the client-side player position on success
//

type ClientPlayer struct {
	Nickname        *string
	IsOnline        bool
	InternalPlayer  internal_player.Player // Interface passed from internal
	HoverTile       *types.Position        // Currently hovered game tile
	DestTile        *types.Position        // Destination of most recent move request
	MoveQueue       []*paths.Cell          // Queued move requests
	MoveQueueCancel chan bool              // Channel for stopping an ongoing MouseMove
}

type MOVE_REQ_TYPE int

const ( // Think of better names for this
	INTENTION MOVE_REQ_TYPE = iota
	ARRIVAL
)

func New(p *internal_player.Player) *ClientPlayer {
	return &ClientPlayer{
		Nickname:        new(string),
		IsOnline:        false,
		InternalPlayer:  *p,
		HoverTile:       nil,
		DestTile:        nil,
		MoveQueue:       []*paths.Cell{},
		MoveQueueCancel: make(chan bool),
	}
}

// Send a REGISTER_ request from Client
func (cp *ClientPlayer) Register() {
	nick := "nick"
	err := cp.InternalPlayer.Register(&nick)
	if err != nil {
		fmt.Println("[CLIENT PLAYER] - Failed to register. Err: " + err.Error())
		return
	}
}

// Send a LOGIN_ request from Client
func (cp *ClientPlayer) Login() {
	err := cp.InternalPlayer.Login()
	if err != nil {
		fmt.Println("[CLIENT PLAYER] - Failed to login. Err: " + err.Error())
		return
	}
}

// Computes a path and queues n MOVE_ requests to destination tile
func (cp *ClientPlayer) MouseMove(cMap *paths.Grid) error {

	if cp.DestTile == nil {
		return fmt.Errorf("[CLIENT PLAYER] - MouseMove err: No Destination is set")
	}

	if len(cp.MoveQueue) != 0 {
		cp.MoveQueueCancel <- true
	}

	path := cMap.GetPath(float64(cp.InternalPlayer.GetPos().Position.X-1), float64(cp.InternalPlayer.GetPos().Position.Y-1), float64(cp.DestTile.X-1), float64(cp.DestTile.Y-1), false, false)

	if path != nil {

		ticker := time.NewTicker(time.Millisecond * time.Duration(c.MOVE_COOLDOWN_IN_MS))

		for i, pathTile := range path.Cells {
			if i != 0 {
				cp.MoveQueue = append(cp.MoveQueue, pathTile, pathTile)
			}
		}

		// Perform move in
		for range cp.MoveQueue {
			go func() {
				select {
				case <-cp.MoveQueueCancel:
					cp.MoveQueue = nil
					return
				case <-ticker.C:

					if len(cp.MoveQueue) <= 0 {
						return
					}

					pPos := cp.InternalPlayer.GetPos().Position
					destPos := cp.InternalPlayer.GetPos().DestPosition
					queuedX, queuedY := cp.MoveQueue[0].X+1, cp.MoveQueue[0].Y+1

					// fmt.Println("pPos", pPos.X, pPos.Y)
					// fmt.Println("destPos", destPos.X, destPos.Y)
					// fmt.Println("queued", queuedX, queuedY)

					if destPos.X == queuedX && destPos.Y == queuedY {
						cp.singleMove(&queuedX, &queuedY, &queuedX, &queuedY, ARRIVAL)
					} else {
						cp.singleMove(&pPos.X, &pPos.Y, &queuedX, &queuedY, INTENTION)
					}

					cp.MoveQueue = cp.MoveQueue[1:]
				}
			}()
		}
	}
	return nil
}

func (cp *ClientPlayer) singleMove(posX, posY, destX, destY *int, reqType MOVE_REQ_TYPE) {
	pPos := cp.InternalPlayer.GetPos().Position
	switch reqType {
	case INTENTION:
		cp.InternalPlayer.Move(&pPos.X, &pPos.Y, destX, destY)
	case ARRIVAL:
		cp.InternalPlayer.Move(destX, destY, destX, destY)
	}
}
