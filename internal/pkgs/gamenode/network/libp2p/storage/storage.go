package storage

import (
	"bytes"
	"context"

	bsnet "github.com/ipfs/boxo/bitswap/network"
	bsserver "github.com/ipfs/boxo/bitswap/server"
	"github.com/ipfs/boxo/blockservice"
	blockstore "github.com/ipfs/boxo/blockstore"
	chunker "github.com/ipfs/boxo/chunker"
	offline "github.com/ipfs/boxo/exchange/offline"
	"github.com/ipfs/boxo/ipld/merkledag"
	"github.com/ipfs/boxo/ipld/unixfs/importer/balanced"
	uih "github.com/ipfs/boxo/ipld/unixfs/importer/helpers"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	routinghelpers "github.com/libp2p/go-libp2p-routing-helpers"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/multiformats/go-multicodec"
)

type Storage interface {
	Add(data string) (cid.Cid, error)
	StartServer(ctx context.Context, h host.Host)
	GetAllLocal(ctx context.Context) ([]string, error)
}

type storage struct {
	blockstore       blockstore.Blockstore
	dagBuilderParams *uih.DagBuilderParams
}

func New() Storage {
	bs := createBlockstore()
	dbp := createDagBuilderParams(bs)
	return &storage{
		blockstore:       bs,
		dagBuilderParams: dbp,
	}
}

func createBlockstore() blockstore.Blockstore {
	ds := dsync.MutexWrap(datastore.NewMapDatastore())
	bs := blockstore.NewBlockstore(ds)
	bs = blockstore.NewIdStore(bs) // handle identity multihashes, these don't require doing any actual lookups
	return bs
}

func createDagBuilderParams(bs blockstore.Blockstore) *uih.DagBuilderParams {

	bsrv := blockservice.New(bs, offline.Exchange(bs))
	dsrv := merkledag.NewDAGService(bsrv)

	dagBuilderParams := &uih.DagBuilderParams{
		Maxlinks:  uih.DefaultLinksPerBlock,
		RawLeaves: true,
		CidBuilder: cid.V1Builder{
			Codec:    uint64(multicodec.DagPb),
			MhType:   uint64(multicodec.Sha2_256),
			MhLength: -1,
		},
		Dagserv: dsrv,
		NoCopy:  false,
	}

	return dagBuilderParams
}

func (s *storage) Add(data string) (cid.Cid, error) {
	// b, err := os.ReadFile(filePath)
	// if err != nil {
	// 	return cid.Undef, err
	// }
	fileReader := bytes.NewReader([]byte(data))

	dagBuilder, err := s.dagBuilderParams.New(chunker.NewSizeSplitter(fileReader, chunker.DefaultBlockSize))
	if err != nil {
		return cid.Undef, nil
	}

	nd, err := balanced.Layout(dagBuilder) // Arrange the graph with a balanced layout
	if err != nil {
		return cid.Undef, nil
	}
	return nd.Cid(), nil
}

func (s *storage) StartServer(ctx context.Context, h host.Host) {
	net := bsnet.NewFromIpfsHost(h, routinghelpers.Null{})
	bswap := bsserver.New(ctx, net, s.blockstore)
	net.Start(bswap)
}

func (s *storage) GetAllLocal(ctx context.Context) ([]string, error) {
	c, _ := s.blockstore.AllKeysChan(ctx)
	cidArr := []string{}

	for k := range c {
		cidArr = append(cidArr, k.String())
	}

	return cidArr, nil
}
