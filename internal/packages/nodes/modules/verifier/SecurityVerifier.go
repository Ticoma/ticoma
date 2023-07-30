package verifier

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SecurityVerifier
type SecurityVerifier struct{}

func (sv *SecurityVerifier) GetPackageSchema(timestamped bool) string {

	// Schema for the classic Action Data Package (local)
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

	// Schema for ADP with a Timestamp key
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

func (sv *SecurityVerifier) VerifyADPTypes(pkg []byte, timestamped bool) bool {

	var schema string

	if timestamped {
		schema = sv.GetPackageSchema(true)
	} else {
		schema = sv.GetPackageSchema(false)
	}

	pkgStr := string(pkg)
	res := []byte{}
	keySelected := false

	// Anti spam
	if len(pkg) == 0 {
		return false
	}

	dec := json.NewDecoder(strings.NewReader(pkgStr))

	for {
		t, err := dec.Token()
		if err != nil {
			break
		}

		// fmt.Printf("%v: %T\n", t, t)

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
	// fmt.Println(StripString(schema, true))
	// fmt.Println(StripString(string(res), true))

	return strings.Compare(StripString(schema, true), StripString(string(res), true)) == 0
}

// temp here
func StripString(str string, removeLastCharToo bool) string {
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\t", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	if removeLastCharToo {
		str = strings.TrimSuffix(str, ",")
	}
	return str
}
