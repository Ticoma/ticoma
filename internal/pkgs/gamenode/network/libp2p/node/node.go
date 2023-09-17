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
		debug.DebugLog("[NETWORK NODE] - Connected to relay!", debug.NETWORK)
	}

	topicName := "ticoma1"
	topic, sub := nn.Host.ConnectToPubsub(ctx, topicName, isRelay)
	debug.DebugLog("[NETWORK NODE] - Connected to pubsub!", debug.NETWORK)

	nn.Topic = topic
	nn.Sub = sub
	nn.isRelay = isRelay
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
	debug.DebugLog("[NETWORK NODE] - Snapshot ticker started", debug.NETWORK)
	snapshotTick := ticker()

	for {
		<-snapshotTick.C
		snapshotTick = ticker()
		nn.SnapshotMsg(ctx, time.Now().UnixMilli())
	}
}

// Returns a new ticker that triggers at the start of each minute
func ticker() *time.Ticker {
	return time.NewTicker(time.Second * time.Duration(60-time.Now().Second()))
}

func (nn *NetworkNode) SnapshotMsg(ctx context.Context, ms int64) {
	time := time.Unix(0, ms*int64(time.Millisecond))
	err := nn.Topic.Publish(ctx, []byte(fmt.Sprintf("Snapshot tick. Time: %v", time)))
	if err != nil {
		debug.DebugLog(fmt.Sprintf("[NETWORK NODE] - Err while performing snapshot. Err: %s", err.Error()), debug.NETWORK)
	}
}
