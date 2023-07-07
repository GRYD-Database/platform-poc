package odb

import (
	orbitdb "berty.tech/go-orbit-db"
	"context"
	icore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/kubo/config"
	"github.com/ipfs/kubo/core"
	"github.com/sirupsen/logrus"
)

type Database struct {
	ctx              context.Context
	ConnectionString string
	URI              string
	CachePath        string

	Logger *logrus.Logger

	IPFSNode    *core.IpfsNode
	IPFSCoreAPI icore.CoreAPI

	OrbitDB orbitdb.OrbitDB
	Store   orbitdb.DocumentStore
}

func NewDatabase(
	ctx context.Context,
	dbConnectionString string,
) (*Database, error) {
	var err error

	db := new(Database)
	db.ctx = ctx
	db.ConnectionString = dbConnectionString

	db.Logger.Debug("getting config root path ...")
	defaultPath, err := config.PathRoot()
	if err != nil {
		return nil, err
	}

	db.Logger.Debug("setting up plugins ...")

	db.Logger.Debug("creating IPFS node ...")
	db.IPFSNode, db.IPFSCoreAPI, err = createNode(ctx, defaultPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}
