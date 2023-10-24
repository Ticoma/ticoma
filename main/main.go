//go:build !noX11

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"ticoma/client"
	"ticoma/internal"
	"ticoma/internal/pkgs/player"
	"ticoma/types"
)

func main() {
	pc := make(chan player.Player)        // channel for Player interface
	crc := make(chan types.CachedRequest) // channel for verified requests
	ctx, cancel := context.WithCancel(context.Background())

	clientF := flag.Bool("client", false, "True = run with GUI, False = headless. Defaults to true")
	relayF := flag.Bool("relay", false, "True = only seed standalone gamenode, no player node. Defaults to false")
	fullscreenF := flag.Bool("fullscreen", false, "True = run in fullscreen mode. Defaults to false")

	flag.Parse()

	go internal.Main(ctx, pc, crc, relayF)
	if *clientF {
		fmt.Println("Starting client")
		client.Main(pc, crc, fullscreenF)
	}
	fmt.Scanln()
	cancel()
	os.Exit(0)
}
