package odb

import (
	"context"
	"github.com/ipfs/go-libipfs/files"
	icore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/kubo/core"
	"github.com/ipfs/kubo/core/coreapi"
	"github.com/ipfs/kubo/core/node/libp2p"
	"github.com/ipfs/kubo/repo/fsrepo"
	"github.com/mitchellh/mapstructure"
	"os"
)

func createNode(ctx context.Context, repoPath string) (*core.IpfsNode, icore.CoreAPI, error) {
	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, nil, err
	}

	nodeOptions := &core.BuildCfg{
		Online:  true,
		Routing: libp2p.DHTClientOption, // DHTOption
		Repo:    repo,
		ExtraOpts: map[string]bool{
			"pubsub": true,
		},
	}

	node, err := core.NewNode(ctx, nodeOptions)
	if err != nil {
		return nil, nil, err
	}

	coreAPI, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, nil, err
	}

	return node, coreAPI, nil
}

func getUnixfsNode(path string) (files.Node, error) {
	st, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	f, err := files.NewSerialFile(path, false, st)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func structToMap(v interface{}) (map[string]interface{}, error) {
	vMap := &map[string]interface{}{}

	err := mapstructure.Decode(v, &vMap)
	if err != nil {
		return nil, err
	}

	return *vMap, nil
}
