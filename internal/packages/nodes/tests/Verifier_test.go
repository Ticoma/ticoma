package tests

import (
	"encoding/json"
	"fmt"
	"testing"
	. "ticoma/packages/nodes/interfaces"
	. "ticoma/packages/nodes/modules/verifier"
)

func TestVerifyPackageTypesRandom(t *testing.T) {

	v := NewVerifier()

	// Completely random package
	pkgStr := `{"foo": bar}`
	got := v.SecurityVerifier.VerifyADPTypes([]byte(pkgStr))
	want := false

	if got != want {
		t.Errorf("got %t, wanted %t", got, want)
	}

	// Empty package
	pkgStrEmpty := ``
	got2 := v.SecurityVerifier.VerifyADPTypes([]byte(pkgStrEmpty))
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

	got1 := v.SecurityVerifier.VerifyADPTypes(pkg1)
	want1 := false

	if got1 != want1 {
		t.Errorf("got %t, wanted %t", got1, want1)
	}

	got2 := v.SecurityVerifier.VerifyADPTypes(pkg2)
	want2 := false

	if got2 != want2 {
		t.Errorf("got %t, wanted %t", got2, want2)
	}

}

func TestVerifyPackageTypesCorrect(t *testing.T) {

	v := NewVerifier()

	pkg := ActionDataPackage{
		PlayerId:     1,
		PubKey:       "PUBKEY",
		Position:     &Position{X: 1, Y: 1},
		DestPosition: &DestPosition{X: 2, Y: 2},
	}

	// to json string first
	jsonBytes, err := json.Marshal(pkg)

	if err != nil {
		fmt.Println(err)
	}

	got := v.SecurityVerifier.VerifyADPTypes(jsonBytes)
	want := true

	if got != want {
		t.Errorf("got %t, wanted %t", got, want)
	}

}
