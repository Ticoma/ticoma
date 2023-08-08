package verifier

import (
	"fmt"
	"math"
	. "ticoma/internal/packages/player/interfaces"
)

// EngineVerifier
type EngineVerifier struct{}

// Engine ruleset
const SECONDS_IN_MS = 1000
const MAX_BLOCKS_TRAVEL_PER_SECOND = 4
const MAX_VEL = float64(MAX_BLOCKS_TRAVEL_PER_SECOND) / float64(SECONDS_IN_MS)

// Engine verify functions

// Checks if the player traversed too many blocks in a short amount of time
func (ev *EngineVerifier) VerifyPlayerMovement(startPkg *ActionDataPackageTimestamped, endPkg *ActionDataPackageTimestamped) bool {

	startPos := startPkg.Position
	endPos := endPkg.Position
	elapsedTime := endPkg.Timestamp - startPkg.Timestamp
	blocksTraveledX := math.Abs(float64(endPos.X) - float64(startPos.X))
	blocksTraveledY := math.Abs(float64(endPos.Y) - float64(startPos.Y))
	blocksTraveledTotal := blocksTraveledX + blocksTraveledY

	velocity := blocksTraveledTotal / float64(elapsedTime)

	// DEBUG
	// fmt.Println("VELOCITY: ", velocity)
	// fmt.Println("MAX VELOCITY: ", MAX_VEL)

	if velocity > MAX_VEL {
		fmt.Printf("[ERR] Engine player movement verification (velocity too high)\nVelocity: %f, max acceptable velocity: %f", velocity, MAX_VEL)
		return false
	}

	return true

}

// Checks if destPos of last package matches the pos of currently arriving package
func (ev *EngineVerifier) VerifyLastMovePos(lastDestPos *DestPosition, pos *Position) bool {
	verX := lastDestPos.X == pos.X
	verY := lastDestPos.Y == pos.Y
	if !verX || !verY {
		return false
	} else {
		return true
	}
}
