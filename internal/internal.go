package internal

import (
	"context"
	"fmt"
	"os"

	"ticoma/internal/pkgs/gamenode"
	"ticoma/internal/pkgs/gamenode/network/libp2p/node"
	"ticoma/internal/pkgs/player"
	"ticoma/types"

	"github.com/joho/godotenv"
)

var nodeConfig node.NodeConfig

func Main(ctx context.Context, pc chan player.Player, cc chan types.ChatMessage, isRelay bool) {

	// Load env
	err := godotenv.Load()
	if err != nil {
		panic("[ERR] couldn't load .env file")
	}

	// Network conf
	port := os.Getenv("PORT")
	relayIp := os.Getenv("RELAY_IP")
	relayAddr := os.Getenv("RELAY_ADDR")
	relayPort := os.Getenv("RELAY_PORT")

	// Check if conf is OK
	if port == "" || relayIp == "" || relayAddr == "" || relayPort == "" {
		panic("[ERR] .env config is incomplete")
	}

	nodeConfig = node.NodeConfig{
		RelayAddr: relayAddr,
		RelayPort: relayPort,
		RelayIp:   relayIp,
		Port:      port,
	}

	if isRelay {
		runStandaloneGameNode(ctx)
	} else {
		runPlayerNode(pc, cc, ctx)
	}
}

func runPlayerNode(pc chan player.Player, cc chan types.ChatMessage, ctx context.Context) {
	// p := player.New(ctx, 0)
	// p.Init(ctx, cc, false, &nodeConfig)
	// fmt.Printf("Player id: %d connected!\n", 0)
	// pc <- p
}

func runStandaloneGameNode(ctx context.Context) {
	gn := gamenode.New()
	gn.Init(ctx, true, &nodeConfig)
	fmt.Println("===============================================================")
	fmt.Println("Relay ID: ", gn.NetworkNode.Host.GetPeerInfo().ID.String())
	fmt.Println("===============================================================")
}
