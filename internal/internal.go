package internal

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	gamenode "ticoma/internal/packages/network/gamenode"
	playernode "ticoma/internal/packages/nodes"
)

func Main(isRelay bool) {

	// Load env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancel := context.WithCancel(context.Background())
	go nodeProcess(ctx, isRelay)

	var opt string
	fmt.Println("Press anything to stop")
	fmt.Scanln(&opt)
	cancel()
	os.Exit(0)
}

func nodeProcess(ctx context.Context, isRelay bool) {
	select {
	case <-ctx.Done():
		fmt.Println("EXIT")
		return
	default:
		if isRelay {
			runStandaloneGameNode(ctx)
		} else {
			runPlayerNode(ctx)
		}
	}
}

func runPlayerNode(ctx context.Context) {

	relayIp := os.Getenv("RELAY_IP")
	relayAddr := os.Getenv("RELAY_ADDR")
	relayPort := "1337"

	nodeConfig := gamenode.GameNodeConfig{
		RelayAddr:          relayAddr,
		RelayIp:            relayIp,
		RelayPort:          relayPort,
		EnableDebugLogging: true,
	}

	pn := playernode.NewPlayerNode()
	pn.InitPlayerNode(ctx, &nodeConfig)

	// receive
	go pn.ListenForPkgs(ctx)

	// send
	for {
		pn.SendPkg(ctx, "Hello")
		time.Sleep(time.Second * 2)
	}

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
