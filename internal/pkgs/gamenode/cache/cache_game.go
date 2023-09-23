package cache

import (
	"fmt"
	"ticoma/internal/pkgs/gamenode/cache/verifier/security"
	"ticoma/types"
)

//
// All game related request functions, handlers for Cache
//

// Sub-handler for all game related requests
func (nc *NodeCache) handleGameRequest(peerID string, reqPrefix string, reqS *interface{}) error {
	switch reqPrefix {
	case security.MOVE_PREFIX:
		return nc.updatePlayerPos(peerID, (*reqS).(types.PlayerPosition))
	case security.CHAT_PREFIX:
		// TODO: think about private msg implementation
		return nil
	default:
		return fmt.Errorf("[NODE CACHE] - Unknown game request (prefix %s)", reqPrefix)
	}

}

func (nc *NodeCache) updatePlayerPos(peerID string, pp types.PlayerPosition) error {

	if !nc.accountOnline(peerID) {
		return fmt.Errorf("[NODE CACHE] - Can't update pos. Player is offline!")
	}

	p := nc.Memory[peerID]

	currPos := p.Curr.PlayerPosition
	validVelocity := nc.EngineVerifier.VerifyMoveVelocity(&currPos, &pp)
	if !validVelocity {
		return fmt.Errorf("[NODE CACHE] - Coulnd't verify move direction or position. %s", "")
	}

	// Verify move position sequence
	validMove := nc.EngineVerifier.VerifyMoveDirection(&currPos.DestPosition, &pp.Position)

	if !validMove {
		return fmt.Errorf("[NODE CACHE] - Engine couldn't verify move. %s", "")
	}

	// Update playerStates with new data (move pStates stack)
	p.Prev.PlayerGameData.Position = types.Position(p.Prev.PlayerGameData.DestPosition)
	p.Prev.DestPosition = types.DestPosition(p.Curr.Position)
	p.Curr.Position = pp.Position
	p.Curr.DestPosition = pp.DestPosition
	nc.Memory[peerID] = p

	return nil
}
