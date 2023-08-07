package main

import (
	"flag"
	"fmt"
	t_client "ticoma/client"
	t_internal "ticoma/internal"
)

func main() {

	clientF := flag.Bool("client", false, "true if internal + client, defaults to false")
	relayF := flag.Bool("relay", false, "true if only want to run relay (seed standalone gamenode)")
	flag.Parse()

	fmt.Println("Opts: ", *clientF, *relayF)

	t_internal.Main(*relayF)
	if *clientF {
		t_client.Main()
	}
}
