package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"ticoma/packages/network/nodes"

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
		testGameNode(ctx)
	}
}

func testGameNode(ctx context.Context) {

	// relayIp := os.Getenv("RELAY_IP")
	// relayAddr := os.Getenv("RELAY_ADDR")
	// relayPort := "1337"

	// Test setup GameNode
	// gn := nodes.NewGameNode()
	// gn.InitGameNode(&ctx, relayAddr, relayIp, relayPort, true)

	// Test setup relay
	sgn := nodes.NewStandaloneGameNode()
	sgn.GameNodeRelay.SetupRelay("0.0.0.0", "1337")

	fmt.Println("Relay peerID ", sgn.RelayHost.ID())
	fmt.Println("Relay Addr ", sgn.RelayHost.Addrs())

	fmt.Println("Done")

	// for {
	// 	gn.Greet(ctx, topic, "test msg")
	// 	time.Sleep(2 * time.Second)
	// }
}
