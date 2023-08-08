package main

import (
	"context"
	"flag"
	"fmt"

	"os"
	t_client "ticoma/client"
	t_internal "ticoma/internal"
	player "ticoma/internal/packages/player"
)

func main() {

	c := make(chan player.PlayerInterface)
	ctx, cancel := context.WithCancel(context.Background())

	clientF := flag.Bool("client", false, "true if internal + client, defaults to false")
	relayF := flag.Bool("relay", false, "true if only want to run relay (seed standalone gamenode)")
	flag.Parse()

	go t_internal.Main(ctx, c, *relayF)
	if *clientF {
		t_client.Main(c)
	}
	fmt.Scanln()
	cancel()
	os.Exit(0)
}
