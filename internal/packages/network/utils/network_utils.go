package utils

import (
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

//
// Util functions for the entire network module / submodules
//

func ConvertToAddrInfo(ip, id, port string) *peer.AddrInfo {
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
