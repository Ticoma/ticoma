package verifier

import (
	"encoding/json"
	"fmt"
	"strings"
	"ticoma/packages/nodes/utils"
)

// SecurityVerifier
type SecurityVerifier struct{}

func (sv *SecurityVerifier) GetPackageSchema(timestamped bool) string {

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

	const schemaADPTimestamped = `{
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
		timestamp: int,
	},`

	if timestamped {
		return schemaADPTimestamped
	}
	return schemaADP

}

// TODO : Remove timestamped option from here. Timestamp will be assigned by the system,
// while constructing a Pkg from string, not by the user. This is only local

func (sv *SecurityVerifier) VerifyADPTypes(pkg []byte, timestamped bool) bool {

	var schema string

	if timestamped {
		schema = sv.GetPackageSchema(true)
	} else {
		schema = sv.GetPackageSchema(false)
	}

	fmt.Println("PKG STR: ", string(pkg)) // DEBUG

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
	fmt.Println("SCHEMA ", utils.StripString(schema, true))
	fmt.Println("RES ", utils.StripString(string(res), true))

	valid := strings.Compare(utils.StripString(schema, true), utils.StripString(string(res), true)) == 0

	return valid
}
