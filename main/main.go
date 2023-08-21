package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"ticoma/client"
	"ticoma/internal"
	"ticoma/internal/packages/player"
	"ticoma/types"
)

func main() {

	pc := make(chan player.Player)     // channel for Player interface
	cc := make(chan types.ChatMessage) // channel for chat messages (all in one chan for now)
	ctx, cancel := context.WithCancel(context.Background())

	clientF := flag.Bool("client", false, "true if internal + client, defaults to false")
	relayF := flag.Bool("relay", false, "true if only want to run relay (seed standalone gamenode)")
	fullscreenF := flag.Bool("fullscreen", false, "true if want to run in fullscreen mode")

	flag.Parse()

	go internal.Main(ctx, pc, cc, *relayF)
	if *clientF {
		fmt.Println("Starting client")
		client.Main(pc, cc, fullscreenF)
	}
	fmt.Scanln()
	cancel()
	os.Exit(0)
}
