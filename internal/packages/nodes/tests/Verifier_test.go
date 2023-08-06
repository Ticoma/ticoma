package tests

import (
	"encoding/json"
	"fmt"
	"testing"
	. "ticoma/packages/nodes/interfaces"
	. "ticoma/packages/nodes/modules/verifier"
	"ticoma/packages/nodes/utils"
)

func TestVerifyPackageTypesRandom(t *testing.T) {

	v := NewVerifier()

	// Completely random package
	pkgStr := `{"foo": bar}`
	got := v.SecurityVerifier.VerifyADPTypes([]byte(pkgStr), false)
	want := false

	if got != want {
		t.Errorf("got %t, wanted %t", got, want)
	}

	// Empty package
	pkgStrEmpty := ``
	got2 := v.SecurityVerifier.VerifyADPTypes([]byte(pkgStrEmpty), false)
	want2 := false

	if got2 != want2 {
		t.Errorf("got %t, wanted %t", got, want)
	}

}
func TestVerifyPackageTypesIncorrect(t *testing.T) {

	v := NewVerifier()

	// Invalid pkgs (keys are not equal to schema / values are of incorrect type)

	// Incorrect key name (PLAYER_ID instead of playerId)
	pkg1 := []byte(`
	{
		"PLAYER_ID":0,
		"pubKey":"PUBKEY",
		"pos":{
			"posX":1,
			"posY":1
		},
		"destPos":{
			"destPosX":1,
			"destPosY":1
		}
	}`)

	// Incorrect type in key (posX is a boolean instead of an integer)
	pkg2 := []byte(`
	{
		"playerId":0,
		"pubKey":"PUBKEY",
		"pos":{
			"posX":true,
			"posY":1
		},
		"destPos":{
			"destPosX":1,
			"destPosY":1
		}
	}`)

	got1 := v.SecurityVerifier.VerifyADPTypes(pkg1, false)
	want1 := false

	if got1 != want1 {
		t.Errorf("got %t, wanted %t", got1, want1)
	}

	got2 := v.SecurityVerifier.VerifyADPTypes(pkg2, false)
	want2 := false

	if got2 != want2 {
		t.Errorf("got %t, wanted %t", got2, want2)
	}

}

func TestVerifyPackageTypesCorrect(t *testing.T) {

	v := NewVerifier()

	// Basic local pkg
	pkg := ActionDataPackage{
		PlayerId:     1,
		PubKey:       "PUBKEY",
		Position:     &Position{X: 1, Y: 1},
		DestPosition: &DestPosition{X: 2, Y: 2},
	}

	// quick test tmp
	// fn, jt := utils.GetInterfaceFieldNames(pkg, true)
	// fmt.Println("FN ", fn, "JT", jt)

	// to json string first
	jsonBytes, err := json.Marshal(pkg)

	if err != nil {
		fmt.Println(err)
	}

	got := v.SecurityVerifier.VerifyADPTypes(jsonBytes, false)
	want := true

	if got != want {
		t.Errorf("got %t, wanted %t", got, want)
	}

	// Verify a timestamped package
	pkg2 := ActionDataPackageTimestamped{
		ActionDataPackage: &ActionDataPackage{
			PlayerId:     1,
			PubKey:       "PUBKEY",
			Position:     &Position{X: 1, Y: 1},
			DestPosition: &DestPosition{X: 2, Y: 2},
		},
		Timestamp: 1337,
	}

	// to json string first
	jsonBytes2, err2 := json.Marshal(pkg2)

	if err2 != nil {
		fmt.Println(err2)
	}

	got2 := v.SecurityVerifier.VerifyADPTypes(jsonBytes2, true)
	want2 := true

	if got2 != want2 {
		t.Errorf("got %t, wanted %t", got2, want2)
	}

}

func TestConvStringToPkg(t *testing.T) {

	v := NewVerifier()

	// Validate string pkg and, if successful, try to convert it to a struct
	testPkg := `{"playerId":0,"pubKey":"PUBKEY","pos":{"posX":1,"posY":1},"destPos":{"destPosX":1,"destPosY":1}}`
	verified := v.SecurityVerifier.VerifyADPTypes([]byte(testPkg), false)
	verifiedWant := true

	if verified != verifiedWant {
		t.Errorf("got %t, wanted %t", verified, verifiedWant)
	}

	// Extract vals
	var vals interface{}
	vals = utils.ExtractValsFromStrPkg(testPkg)
	fmt.Println(vals)

	// types are OK, try construct

	adp := &ActionDataPackage{}
	utils.ConstructPkg(adp, vals, false)

}
