package internal

import (
	"context"
	"fmt"
	"os"

	gamenode "ticoma/internal/packages/network/gamenode"
	player "ticoma/internal/packages/player"

	"github.com/joho/godotenv"
)

func Main(ctx context.Context, c chan player.PlayerInterface, isRelay bool) {

	// Load env
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	if isRelay {
		runStandaloneGameNode(ctx)
	} else {
		runPlayerNode(c, ctx)
	}
}

func runPlayerNode(c chan player.PlayerInterface, ctx context.Context) {

	relayIp := os.Getenv("RELAY_IP")
	relayAddr := os.Getenv("RELAY_ADDR")
	relayPort := "1337"

	nodeConfig := gamenode.GameNodeConfig{
		RelayAddr:          relayAddr,
		RelayIp:            relayIp,
		RelayPort:          relayPort,
		EnableDebugLogging: true,
	}

	p := player.New(ctx, &nodeConfig)

	c <- p
}

func runStandaloneGameNode(ctx context.Context) {

	relayPort := "1337"

	rel := gamenode.NewStandaloneGameNode()
	rel.SetupRelay("0.0.0.0", relayPort)

	fmt.Println("========================")
	fmt.Println("Relay ID: ", rel.RelayHost.ID().String())
	fmt.Println("Relay port: ", relayPort)
	fmt.Println("========================")
}
