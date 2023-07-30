package verifier

import (
	"fmt"
	"math"
	. "ticoma/packages/nodes/interfaces"
)

// EngineVerifier
type EngineVerifier struct{}

// Engine ruleset
const SECONDS_IN_MS = 1000
const MAX_BLOCKS_TRAVEL_PER_SECOND = 4
const MAX_VEL = MAX_BLOCKS_TRAVEL_PER_SECOND / SECONDS_IN_MS

// Engine verify functions

// Movement
func (ev *EngineVerifier) VerifyPlayerMovement(startPkg *ActionDataPackageTimestamped, endPkg *ActionDataPackageTimestamped) bool {

	startPos := startPkg.Position
	endPos := endPkg.Position
	elapsedTime := endPkg.Timestamp - startPkg.Timestamp
	blocksTraveledX := math.Abs(float64(endPos.X) - float64(startPos.X))
	blocksTraveledY := math.Abs(float64(endPos.Y) - float64(startPos.Y))
	blocksTraveledTotal := blocksTraveledX + blocksTraveledY

	velocity := int(blocksTraveledTotal) / elapsedTime

	if velocity > MAX_VEL {
		fmt.Printf("[ERR] Engine player movement verification (velocity too high)\nVelocity: %d, max acceptable velocity: %d", velocity, MAX_VEL)
		return false
	}

	return true

}
