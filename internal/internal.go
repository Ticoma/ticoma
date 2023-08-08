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
var relayIp = os.Getenv("RELAY_IP")
var relayAddr = os.Getenv("RELAY_ADDR")
var relayPort = "1337"
var nodeConfig = gamenode.NodeConfig{
	RelayAddr: relayAddr,
	RelayIp:   relayIp,
	RelayPort: relayPort,
}

func Main(ctx context.Context, c chan player.PlayerInterface, isRelay bool) {

	// Load env
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	// TODO: Check if all .env vars are != nil based on flags
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

	fmt.Println("========================")
	fmt.Println("Relay ID: ", gn.GetPeerInfo().ID)
	fmt.Println("Relay port: ", relayPort)
	fmt.Println("========================")
}
