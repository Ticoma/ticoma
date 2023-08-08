package internal

import (
	"context"
	"fmt"
	"os"

	gamenode "ticoma/internal/packages/network/gamenode"
	player "ticoma/internal/packages/player"

	"github.com/joho/godotenv"
)

// conf
var nodeConfig = gamenode.NodeConfig{}

func Main(ctx context.Context, c chan player.PlayerInterface, isRelay bool) {

	// Load env
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	relayIp := os.Getenv("RELAY_IP")
	relayAddr := os.Getenv("RELAY_ADDR")
	relayPort := "1337"
	nodeConfig.RelayAddr = relayAddr
	nodeConfig.RelayPort = relayPort
	nodeConfig.RelayIp = relayIp

	// TODO:
	// Check if all .env vars are != nil based on flags
	// to prevent 10-line errs from libp2p

	if isRelay {
		runStandaloneGameNode(ctx)
	} else {
		runPlayerNode(c, ctx)
	}
}

func runPlayerNode(c chan player.PlayerInterface, ctx context.Context) {

	nodeConfig.IsRelay = false
	p := player.New(ctx, &nodeConfig)

	c <- p
}

func runStandaloneGameNode(ctx context.Context) {

	gn := gamenode.New()
	nodeConfig.IsRelay = true
	gn.InitGameNode(ctx, &nodeConfig)

	fmt.Println("===============================================================")
	fmt.Println("Relay ID: ", gn.GetPeerInfo().ID)
	fmt.Println("===============================================================")
}
