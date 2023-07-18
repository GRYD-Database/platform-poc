package odb

import (
	orbitdb "berty.tech/go-orbit-db"
	"berty.tech/go-orbit-db/accesscontroller"
	"berty.tech/go-orbit-db/stores"
	"berty.tech/go-orbit-db/stores/documentstore"
	"context"
	"fmt"
	icore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/kubo/config"
	"github.com/ipfs/kubo/core"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/event"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

const replicate = true
const create = true

type Database struct {
	ctx              context.Context
	ConnectionString string
	URI              string
	CachePath        string
	isLocal          bool
	isReplicated     bool
	createRepo       bool
	repoPath         string

	Logger *logrus.Logger

	IPFSNode    *core.IpfsNode
	IPFSCoreAPI icore.CoreAPI

	OrbitDB orbitdb.OrbitDB
	Store   orbitdb.DocumentStore
	Ledger  orbitdb.DocumentStore

	Events event.Subscription
}

func NewDatabase(
	ctx context.Context,
	dbConnectionString,
	repoPath string,
	isLocal,
	createRepo,
	isReplicated bool,
	logger *logrus.Logger,
) (*Database, error) {
	var err error

	db := new(Database)
	db.ctx = ctx
	db.ConnectionString = dbConnectionString
	db.Logger = logger
	db.isLocal = isLocal
	db.repoPath = repoPath
	db.isReplicated = isReplicated
	db.createRepo = createRepo

	if len(strings.TrimSpace(db.repoPath)) == 0 {
		db.Logger.Debug("getting config root path ...")
		db.repoPath, err = config.PathRoot()
		if err != nil {
			logger.Error("getting config root path ...", err)
			return nil, err
		}
	}

	db.Logger.Debug("setting up plugins ...")
	if err := setupPlugins(db.repoPath); err != nil {
		logger.Error("setting up plugins ...", err)
		return nil, err
	}

	db.Logger.Debug("creating IPFS node ...")
	db.IPFSNode, db.IPFSCoreAPI, err = createNode(ctx, db.repoPath)
	if err != nil {
		logger.Error("creating IPFS node ...", err)
		return nil, err
	}

	err = db.OrbitBootstrapper()
	if err != nil {
		logger.Error("err bootstrapping ODB ...", err)
		return nil, err
	}

	logger.Info("successfully bootstrapped ODB and IPFS")

	return db, nil
}

func (d *Database) OrbitBootstrapper() error {
	ac := accesscontroller.NewEmptyManifestParams()

	ac.SetAccess("write", []string{"*"})
	ac.SetAccess("read", []string{"*"})

	odb, err := orbitdb.NewOrbitDB(d.ctx, d.IPFSCoreAPI, nil)
	if err != nil {
		return fmt.Errorf("error bootstrapping ODB: %w", err)
	}
	d.OrbitDB = odb

	storetype := "docstore"

	d.Logger.Debug("initializing OrbitDB.Docs ...")

	d.Store, err = d.OrbitDB.Docs(d.ctx, d.ConnectionString, &orbitdb.CreateDBOptions{
		AccessController:  ac,
		LocalOnly:         &d.isLocal,
		StoreType:         &storetype,
		StoreSpecificOpts: documentstore.DefaultStoreOptsForMap("id"),
		Timeout:           time.Second * 600,
		Replicate:         &d.isReplicated,
		Create:            &d.createRepo,
	})
	if err != nil {
		return err
	}

	d.Ledger, err = d.OrbitDB.Docs(d.ctx, d.ConnectionString, &orbitdb.CreateDBOptions{
		AccessController:  ac,
		StoreType:         &storetype,
		StoreSpecificOpts: documentstore.DefaultStoreOptsForMap("key"),
		Timeout:           time.Second * 600,
	})
	if err != nil {
		return err
	}

	d.Logger.Debug("subscribing to EventBus ...")
	d.Events, err = d.Store.EventBus().Subscribe(new(stores.EventReady))
	if err != nil {
		return fmt.Errorf("error subscribing to odb events: %w", err)
	}

	err = d.Store.Load(d.ctx, -1)
	if err != nil {
		d.Logger.Error(err)
		return err
	}

	return nil
}

func (d *Database) GetOwnID() string {
	return d.OrbitDB.Identity().ID
}

func (d *Database) GetOwnPubKey() crypto.PubKey {
	pubKey, err := d.OrbitDB.Identity().GetPublicKey()
	if err != nil {
		return nil
	}

	return pubKey
}
