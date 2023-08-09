package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"ticoma/client"
	"ticoma/internal"
	"ticoma/internal/packages/player"
)

func main() {

	c := make(chan player.Player)
	ctx, cancel := context.WithCancel(context.Background())

	clientF := flag.Bool("client", false, "true if internal + client, defaults to false")
	relayF := flag.Bool("relay", false, "true if only want to run relay (seed standalone gamenode)")
	flag.Parse()

	go internal.Main(ctx, c, *relayF)
	if *clientF {
		client.Main(c)
	}
	fmt.Scanln()
	cancel()
	os.Exit(0)
}
