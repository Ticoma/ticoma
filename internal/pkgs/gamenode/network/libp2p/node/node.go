package node

import (
	"context"
	"fmt"
	"ticoma/internal/debug"
	"ticoma/internal/pkgs/gamenode/network/libp2p/node/host"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

type NetworkNode struct {
	Host    host.Libp2pHost
	Topic   *pubsub.Topic
	Sub     *pubsub.Subscription
	isRelay bool
}

// Config opts for NetworkNode initialization
type NodeConfig struct {
	Port      string
	RelayAddr string
	RelayIp   string
	RelayPort string
}

func New() *NetworkNode {
	return &NetworkNode{}
}

// Init NetworkNode and connect to pubsub
//
// To set up a relay NetworkNode (Standalone), isRelay must be set to true
func (nn *NetworkNode) Init(ctx context.Context, isRelay bool, nodeConfig *NodeConfig) {

	err := nn.Host.SetupHost("0.0.0.0", nodeConfig.Port)
	if err != nil {
		panic(err)
	}

	if !isRelay {
		relayInfo := convertToAddrInfo(nodeConfig.RelayIp, nodeConfig.RelayAddr, nodeConfig.RelayPort)
		nn.Host.ConnectToRelay(ctx, *relayInfo)
		debug.DebugLog("Connected to relay!", debug.NETWORK)
	}

	topic, sub := nn.Host.ConnectToPubsub(ctx, "ticoma1", isRelay)
	debug.DebugLog("Connected to pubsub!", debug.NETWORK)

	// Storage (ipfs filesystem) tests
	// h := nn.Host.GetHost()
	// s := storage.New()
	// s.StartServer(ctx, h)
	// s.Add("Hello")
	// s.Add("World!")
	// // fmt.Println(cid, err)
	// data, _ := s.GetLocal(ctx)
	// fmt.Println(data)

	nn.Topic = topic
	nn.Sub = sub
	nn.isRelay = isRelay

	go nn.PerformSnapshot(ctx)
}

// Convert string address data -> addrInfo struct
func convertToAddrInfo(ip, id, port string) *peer.AddrInfo {
	m, err := ma.NewMultiaddr("/ip4/" + ip + "/tcp/" + port + "/p2p/" + id)
	if err != nil {
		panic(err)
	}

	addrInfo, err := peer.AddrInfoFromP2pAddr(m)
	if err != nil {
		panic(err)
	}

	return addrInfo
}

// Snapshot timer prototype
func (nn *NetworkNode) PerformSnapshot(ctx context.Context) {
	var ts int64
	var tick <-chan time.Time

	for {
		ts = time.Now().UnixMilli()
		if ts%60000 == 0 {
			tick = time.Tick(time.Minute)
			nn.SnapshotMsg(ctx, time.Now().String())
			break
		}
	}

	for ctime := range tick {
		nn.SnapshotMsg(ctx, ctime.String())
	}
}

func (nn *NetworkNode) SnapshotMsg(ctx context.Context, time string) {
	nodeId := nn.Host.GetPeerInfo().ID.String()
	nn.Topic.Publish(ctx, []byte(fmt.Sprintf("Hello from node %s. Time: %v", nodeId, time)))
}
