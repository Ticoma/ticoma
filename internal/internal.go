package internal

import (
	"context"
	"flag"
	"fmt"
	"os"

	"ticoma/internal/packages/gamenode"
	"ticoma/internal/packages/gamenode/network/libp2p/node"
	"ticoma/internal/packages/player"

	"github.com/joho/godotenv"
)

var nodeConfig = &node.NodeConfig{}

var portFlag = flag.String("port", "1337", "Listening port (default 1337)")
var idFlag = flag.Int("id", 0, "Player id")

func Main(ctx context.Context, c chan player.Player, isRelay bool) {

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
	nodeConfig.Port = *portFlag

	// TODO:
	// Check if all .env vars are != nil based on flags
	// to prevent 10-line errs from libp2p

	if isRelay {
		runStandaloneGameNode(ctx)
	} else {
		runPlayerNode(c, ctx)
	}
}

func runPlayerNode(c chan player.Player, ctx context.Context) {
	p := player.New(ctx, *idFlag)
	p.Init(ctx, false, nodeConfig)
	fmt.Printf("Player id: %d connected!\n", *idFlag)
	c <- p
}

func runStandaloneGameNode(ctx context.Context) {
	gn := gamenode.New()
	gn.Init(ctx, true, nodeConfig)
	fmt.Println("===============================================================")
	fmt.Println("Relay ID: ", gn.NetworkNode.Host.GetPeerInfo().ID.String())
	fmt.Println("===============================================================")
}
