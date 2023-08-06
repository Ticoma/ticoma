package verifier

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"ticoma/packages/nodes/interfaces"
	"ticoma/packages/nodes/utils"
	"time"
)

// SecurityVerifier
type SecurityVerifier struct{}

func (sv *SecurityVerifier) GetPackageSchema() string {
	const schemaADP = `{
		playerId: int,
		pubKey: string,
		pos: {
			posX: int,
			posY: int,
		},
		destPos: {
			destPosX: int,
			destPosY: int,
		},
	},`

	return schemaADP
}

func (sv *SecurityVerifier) VerifyADPTypes(pkg []byte) bool {

	schema := sv.GetPackageSchema()

	// fmt.Println("PKG STR: ", string(pkg)) // DEBUG

	res := []byte{}
	keySelected := false

	// Anti spam
	if len(pkg) == 0 {
		return false
	}

	dec := json.NewDecoder(strings.NewReader(string(pkg)))

	for {
		t, err := dec.Token()
		if err != nil {
			break
		}

		// fmt.Printf("[TEST] %v: %T\n", t, t)

		switch v := t.(type) {
		case json.Delim:
			if string(v) == "}" || string(v) == "]" {
				res = append(res, []byte(string(v)+", ")...)
			} else {
				res = append(res, byte(v))
			}
			if keySelected {
				keySelected = false
			}
		case string:
			if !keySelected {
				keySelected = true
				res = append(res, []byte(v+": ")...)
			} else {
				res = append(res, []byte(fmt.Sprintf("%T, ", v))...)
				keySelected = false
			}
		case float64:
			res = append(res, []byte(fmt.Sprintf("%T, ", int(v)))...)
			keySelected = false
		default:
			res = append(res, []byte(fmt.Sprintf("%T, ", v))...)
			keySelected = false
		}
	}

	// DEBUG
	fmt.Println("[ADP TYPES] SCHEMA ", utils.StripString(schema, true))
	fmt.Println("[ADP TYPES] RES ", utils.StripString(string(res), true))

	valid := strings.Compare(utils.StripString(schema, true), utils.StripString(string(res), true)) == 0

	return valid
}

// Try to construct an ADPT based on provided string pkg
// (Types of pkg should be verified at this point)
func (sv *SecurityVerifier) ConstructADPT(pkgBytes []byte) (interfaces.ActionDataPackageTimestamped, error) {

	const EXPECTED_VAL_LENGTH_IN_ADP = 6 // [playerId, pubKey, posX, posY, destX, destY]

	// Type check
	validPkgTypes := sv.VerifyADPTypes(pkgBytes)
	if !validPkgTypes {
		return interfaces.ActionDataPackageTimestamped{}, fmt.Errorf("[SEC VER] - Couldn't verify package types.")
	}

	// If types are OK, try extract vals
	vals := utils.ExtractValsFromStrPkg(string(pkgBytes))
	if len(vals) != EXPECTED_VAL_LENGTH_IN_ADP {
		return interfaces.ActionDataPackageTimestamped{}, fmt.Errorf("[SEC VER] - Couldn't extract - pkg values length don't match schema.")
	}

	playerId, err := strconv.Atoi(vals[0])
	if err != nil {
		return interfaces.ActionDataPackageTimestamped{}, fmt.Errorf("[SEC VER] - Couldn't assign playerId from vals.")
	}

	var positions []int
	// conv vals[2:5] to ints
	for i := 2; i < len(vals); i++ {
		fmt.Println("VALS[i] ", vals[i])
		pos, err := strconv.Atoi(vals[i])
		if err != nil {
			// fmt.Println("[SEC VER] - Err while converting string val to int")
			return interfaces.ActionDataPackageTimestamped{}, fmt.Errorf("[SEC VER] - Err while converting string val to int")
		}
		positions = append(positions, pos)
	}

	timestamp := time.Now().Unix()

	ADPT := interfaces.ActionDataPackageTimestamped{
		ActionDataPackage: &interfaces.ActionDataPackage{
			PlayerId:     playerId,
			PubKey:       vals[1],
			Position:     &interfaces.Position{X: positions[0], Y: positions[1]},
			DestPosition: &interfaces.DestPosition{X: positions[2], Y: positions[3]},
		},
		Timestamp: timestamp,
	}

	return ADPT, nil
}
