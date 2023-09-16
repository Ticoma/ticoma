package tests

import (
	"fmt"
	"os"
	"testing"
	"ticoma/internal/pkgs/gamenode/cache"

	"github.com/joho/godotenv"
	assert "github.com/stretchr/testify/assert"
)

var playerID string = "foo"

func TestMain(t *testing.M) {
	Setup()
	code := t.Run()
	os.Exit(code)
}

func Setup() {
	// Load debug from env
	err := godotenv.Load()
	if err != nil {
		panic("[ERR] couldn't load .env file in tests directory. " + err.Error())
	}
	// Init cache
	fmt.Println("Test setup OK")
}

func TestRegisterPlayer(t *testing.T) {

	c := cache.New()

	registerData := []byte("REGISTER_")
	c.Put(playerID, registerData)

	p := c.GetPlayer(playerID)

	assert.Equal(t, true, p.Prev.IsOnline)
	assert.Equal(t, true, p.Curr.IsOnline)
}

func TestLogoutPlayer(t *testing.T) {

	c := cache.New()

	// Register first
	registerData := []byte("REGISTER_")
	c.Put(playerID, registerData)

	p := c.GetPlayer(playerID)

	// CHeck if registered
	assert.Equal(t, true, p.Prev.IsOnline)
	assert.Equal(t, true, p.Curr.IsOnline)

	// "foo" should be able to logout now
	logoutData := []byte("LOGOUT_")
	c.Put(playerID, logoutData)

	pl := c.GetPlayer(playerID)

	assert.Equal(t, false, pl.Curr.IsOnline)

	// player "bar" should not be able to logout

	_, err := c.Put("bar", []byte("LOGOUT_"))

	// Put should reject this request
	assert.Error(t, err)
}

func TestLoginPlayer(t *testing.T) {

	id := "login"

	// Register, logout, and then log back in
	c := cache.New()

	c.Put(id, []byte("REGISTER_"))

	// Check if acc created
	assert.Equal(t, true, c.GetPlayer(id).Prev.IsOnline)
	assert.Equal(t, true, c.GetPlayer(id).Curr.IsOnline)

	// Should not be able to log in when already is logged in
	_, err := c.Put(id, []byte("LOGIN_"))
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

	c.Put(id, []byte("REGISTER_"))

	// Check if acc created
	assert.Equal(t, true, c.GetPlayer(id).Prev.IsOnline)
	assert.Equal(t, true, c.GetPlayer(id).Curr.IsOnline)

	// Delete account

	c.Put(id, []byte("DELETEACC_"))

	// Cache should modify prevIsOnline to false (delete entirely)
	assert.Equal(t, false, c.GetPlayer(id).Prev.IsOnline)

}
