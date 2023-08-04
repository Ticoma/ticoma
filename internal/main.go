package main

import (
	"context"
	"fmt"
	"os"

	// nodes "ticoma/packages/nodes"
	gameNode "ticoma/packages/network/nodes"
)

func main() {
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
		gameNode.InitGameNode()
	}
}
