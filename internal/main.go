package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"ticoma/packages/network/nodes"
	"ticoma/packages/network/utils"
	"time"

	"github.com/joho/godotenv"
)

var flagAddress = flag.String("addr", "", "Address")

func main() {

	// Load env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	flag.Parse()

	if *flagAddress == "" {
		panic("No addr flag!")
	}

	fmt.Printf("Hello from main\n")

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
	gn := nodes.NewGameNode()
	fmt.Println("GameNode initialized")

	gn.SetupHost("127.0.0.1", "1337")
	fmt.Println("GameNode host set up")

	relayIp := os.Getenv("RELAY_IP")
	relayInfo := utils.ConvertToAddrInfo(relayIp, *flagAddress, "1337")
	gn.ConnectToRelay(ctx, *relayInfo)
	fmt.Println("GameNode connected to relay")

	gn.ReserveSlot(ctx, *relayInfo)
	fmt.Println("GameNode relay slot reserved")

	topic, _ := gn.ConnectToPubsub(ctx, "ticoma1", true)
	fmt.Println("Connected to pubsub!")

	for {
		gn.Greet(ctx, topic, "test msg")
		time.Sleep(2 * time.Second)
	}
}
