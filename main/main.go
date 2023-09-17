package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"ticoma/client"
	"ticoma/internal"
	"ticoma/internal/pkgs/player"
)

func main() {
	pc := make(chan player.Player) // channel for Player interface
	rc := make(chan interface{})   // channel for network requests
	ctx, cancel := context.WithCancel(context.Background())

	clientF := flag.Bool("client", false, "true if internal + client, defaults to false")
	relayF := flag.Bool("relay", false, "true if only want to run relay (seed standalone gamenode)")
	fullscreenF := flag.Bool("fullscreen", false, "true if want to run in fullscreen mode")

	flag.Parse()

	go internal.Main(ctx, pc, rc, *relayF)
	if *clientF {
		fmt.Println("Starting client")
		client.Main(pc, rc, fullscreenF)
	}
	fmt.Scanln()
	cancel()
	os.Exit(0)
}
