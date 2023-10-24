package tests

import (
	"fmt"
	"testing"
	"ticoma/internal/pkgs/gamenode/cache"

	assert "github.com/stretchr/testify/assert"
)

var playerID string = "foo"
var playerNickname string = "k1ng_V0n"

func TestRegisterPlayer(t *testing.T) {

	c := cache.New()

	// Invalid nicks
	invalidReq1 := []byte(fmt.Sprintf(`REGISTER_{"nickname":"%s"}`, "$KappaDaniels"))
	invalidReq2 := []byte(fmt.Sprintf(`REGISTER_{"nickname":"%s"}`, "ThisNicknameIsCertainlyTooLong"))
	invalidReq3 := []byte(fmt.Sprintf(`REGISTER_{"nickname":"%s"}`, "a"))

	_, _, err := c.Put(playerID, invalidReq1)
	assert.Error(t, err)

	_, _, err = c.Put(playerID, invalidReq2)
	assert.Error(t, err)

	_, _, err = c.Put(playerID, invalidReq3)
	assert.Error(t, err)

	// Valid nick
	registerData := []byte(fmt.Sprintf(`REGISTER_{"nickname":"%s"}`, playerNickname))
	c.Put(playerID, registerData)

	p := c.GetPlayer(playerID)

	assert.Equal(t, true, p.Prev.IsOnline)
	assert.Equal(t, true, p.Curr.IsOnline)
	assert.Equal(t, playerNickname, p.Prev.Nick)
	assert.Equal(t, playerNickname, p.Curr.Nick)
}

func TestLogoutPlayer(t *testing.T) {

	c := cache.New()

	// Register first
	registerData := []byte(fmt.Sprintf(`REGISTER_{"nickname":"%s"}`, playerNickname))
	c.Put(playerID, registerData)

	p := c.GetPlayer(playerID)

	// CHeck if registered
	assert.Equal(t, true, p.Prev.IsOnline)
	assert.Equal(t, true, p.Curr.IsOnline)

	// "foo" should be able to logout now
	logoutData := []byte("LOGOUT_")
	_, _, err := c.Put(playerID, logoutData)
	if err != nil {
		t.Errorf("Error on logout: %s", err.Error())
	}

	pl := c.GetPlayer(playerID)

	assert.Equal(t, false, pl.Curr.IsOnline)

	// player "bar" should not be able to logout
	_, _, err = c.Put("bar", []byte("LOGOUT_"))

	// Put should reject this request
	assert.Error(t, err)
}

func TestLoginPlayer(t *testing.T) {

	id := "login"

	// Register, logout, and then log back in
	c := cache.New()

	c.Put(id, []byte(fmt.Sprintf(`REGISTER_{"nickname":"%s"}`, playerNickname)))

	// Check if acc created
	assert.Equal(t, true, c.GetPlayer(id).Prev.IsOnline)
	assert.Equal(t, true, c.GetPlayer(id).Curr.IsOnline)

	// Should not be able to log in when already is logged in
	_, _, err := c.Put(id, []byte("LOGIN_"))
	assert.Error(t, err)

	c.Put(id, []byte("LOGOUT_"))

	// Logged out?
	assert.Equal(t, true, c.GetPlayer(id).Prev.IsOnline)
	assert.Equal(t, false, c.GetPlayer(id).Curr.IsOnline)

	// Log back in
	c.Put(id, []byte("LOGIN_"))

	assert.Equal(t, true, c.GetPlayer(id).Curr.IsOnline)

}
func TestDeletePlayer(t *testing.T) {

	id := "test_delete"

	// Register and delete account
	c := cache.New()

	c.Put(id, []byte(fmt.Sprintf(`REGISTER_{"nickname":"%s"}`, playerNickname)))

	// Check if acc created
	assert.Equal(t, true, c.GetPlayer(id).Prev.IsOnline)
	assert.Equal(t, true, c.GetPlayer(id).Curr.IsOnline)

	// Delete account

	c.Put(id, []byte("DELETEACC_"))

	// Cache should modify prevIsOnline to false (delete entirely)
	assert.Equal(t, false, c.GetPlayer(id).Prev.IsOnline)

}
