package main

import (
	"context"
	"fmt"
	"log"
	"os"
	gamenode "ticoma/packages/network/gamenode"
	playernode "ticoma/packages/nodes"

	"github.com/joho/godotenv"
)

func main() {

	// Load env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancel := context.WithCancel(context.Background())
	go nodeProcess(ctx)

	var opt string
	fmt.Println("Press anything to stop")
	fmt.Scanln(&opt)
	cancel()
	os.Exit(0)
}

func nodeProcess(ctx context.Context) {
	select {
	case <-ctx.Done():
		fmt.Println("EXIT")
		return
	default:
		// testGameNode(ctx)
		testPlayerNode(ctx)
	}
}

func testPlayerNode(ctx context.Context) {
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

	// send
	// for {
	// 	pn.SendPkg(ctx, "Hello!")
	// 	time.Sleep(time.Second * 2)
	// }

	// receive
	// pn.ListenForPkgs(ctx)

}

func testGameNode(ctx context.Context) {

	// conf
	relayIp := os.Getenv("RELAY_IP")
	relayAddr := os.Getenv("RELAY_ADDR")
	relayPort := "1337"

	// =====================
	// Setup relay first

	rel := gamenode.NewStandaloneGameNode()
	rel.SetupRelay(relayIp, relayPort)

	// fmt.Println("Relay pID: ", rel.RelayHost.ID().String())

	nodeConfig := gamenode.GameNodeConfig{
		RelayAddr:          relayAddr,
		RelayIp:            relayIp,
		RelayPort:          relayPort,
		EnableDebugLogging: true,
	}

	fmt.Println("Local relay set up")

	// =====================
	// Test setup GameNode

	fmt.Println("Setting up GameNode")

	gn := gamenode.NewGameNode()
	gn.InitGameNode(ctx, &nodeConfig)

	fmt.Println("Done")
}
