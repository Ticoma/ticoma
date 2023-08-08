package internal

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	gamenode "ticoma/internal/packages/network/gamenode"
	player "ticoma/internal/packages/player"
)

func Main(ctx context.Context, c chan player.Player, isRelay bool) {

	// Load env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if isRelay {
		runStandaloneGameNode(ctx)
	} else {
		runPlayerNode(c, ctx)
	}
}

// func nodeProcess(ctx context.Context, isRelay bool) {
// 	select {
// 	case <-ctx.Done():
// 		fmt.Println("EXIT")
// 		return
// 	default:
// 		if isRelay {
// 			runStandaloneGameNode(ctx)
// 		} else {
// 			runPlayerNode(ctx)
// 		}
// 	}
// }

func runPlayerNode(c chan player.Player, ctx context.Context) {

	relayIp := os.Getenv("RELAY_IP")
	relayAddr := os.Getenv("RELAY_ADDR")
	relayPort := "1337"

	nodeConfig := gamenode.GameNodeConfig{
		RelayAddr:          relayAddr,
		RelayIp:            relayIp,
		RelayPort:          relayPort,
		EnableDebugLogging: true,
	}

	pn := player.NewPlayerNode()
	pn.InitPlayerNode(ctx, &nodeConfig)

	p := &player.Player{
		PlayerNode: pn,
	}

	c <- *p
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
