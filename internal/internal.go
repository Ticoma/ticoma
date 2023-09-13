package internal

import (
	"context"
	"fmt"
	"os"

	"ticoma/internal/pkgs/gamenode"
	"ticoma/internal/pkgs/gamenode/cache"
	"ticoma/internal/pkgs/gamenode/network/libp2p/node"
	"ticoma/internal/pkgs/player"

	"github.com/joho/godotenv"
)

var nodeConfig node.NodeConfig

func Main(ctx context.Context, pc chan player.Player, rc chan interface{}, isRelay bool) {

	// Load env
	err := godotenv.Load()
	if err != nil {
		panic("[ERR] couldn't load .env file")
	}

	// Network conf
	port := os.Getenv("PORT")
	relayIp := os.Getenv("RELAY_IP")
	relayAddr := os.Getenv("RELAY_ADDR")
	relayPort := os.Getenv("RELAY_PORT")

	// Check if conf is OK
	if port == "" || relayIp == "" || relayAddr == "" || relayPort == "" {
		panic("[ERR] .env config is incomplete")
	}

	nodeConfig = node.NodeConfig{
		RelayAddr: relayAddr,
		RelayPort: relayPort,
		RelayIp:   relayIp,
		Port:      port,
	}

	if isRelay {
		runStandaloneGameNode(ctx)
	} else {
		runPlayerNode(pc, rc, ctx)
	}
}

func runPlayerNode(pc chan player.Player, rc chan interface{}, ctx context.Context) {
	c := cache.New()
	p := player.New(ctx)
	p.Init(ctx, rc, false, &nodeConfig)
	fmt.Printf("Player pID: %s connected!\n", p.GetPeerID())

	pid := "foo"
	moveReqPrefix := "MOVE_"
	moveReqData := fmt.Sprintf(`{"pos":{"posX":%d,"posY":%d},"destPos":{"destPosX":%d,"destPosY":%d}}`, 1, 1, 2, 2)
	c.Put(pid, []byte(moveReqPrefix+moveReqData))

	// fmt.Println(moveReqPrefix + moveReqData)
	paaa := c.GetPlayer(pid)
	fmt.Println(paaa)

	// fmt.Printf("Pos: %v", p.GetPos()) // test
	pc <- p
}

func runStandaloneGameNode(ctx context.Context) {
	gn := gamenode.New()
	gn.Init(ctx, true, &nodeConfig)
	fmt.Println("===============================================================")
	fmt.Println("Relay ID: ", gn.NetworkNode.Host.GetPeerInfo().ID.String())
	fmt.Println("===============================================================")
}
