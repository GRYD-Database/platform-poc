package storage

import (
	orbitdb "berty.tech/go-orbit-db"
	"berty.tech/go-orbit-db/iface"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

type InputData struct {
	DatasetKey string `mapstructure:"datasetKey" json:"datasetKey"`
	ID         string `mapstructure:"id" json:"id"`
	Dataset    string `mapstructure:"dataset" json:"dataset"`
	Date       string `mapstructure:"date" json:"date"`
	DataType   string `mapstructure:"dataType" json:"dataType"`
	Data       string `mapstructure:"data" json:"data"`
}

// Ledger holds the dataset key and the wallet that inserted the data
type Ledger struct {
	Key    string `mapstructure:"key" json:"-"`
	Wallet string `mapstructure:"wallet" json:"-"`
}

type Storage struct {
	logger   *logrus.Logger
	pg       *pgxpool.Pool
	odbStore orbitdb.DocumentStore
	ledger   orbitdb.DocumentStore
	owner    common.Address
}

func New(owner common.Address, logger *logrus.Logger, pool *pgxpool.Pool, store orbitdb.DocumentStore, ledger orbitdb.DocumentStore) *Storage {
	return &Storage{
		logger:   logger,
		pg:       pool,
		owner:    owner,
		odbStore: store,
		ledger:   ledger,
	}
}

func (s *Storage) Ledger(ctx context.Context, wallet, datasetKey string) error {
	record := Ledger{
		Key:    datasetKey,
		Wallet: wallet,
	}

	ledger, err := structToMap(record)
	if err != nil {
		return fmt.Errorf("unable to add recrod to ledger: %w", err)
	}

	_, err = s.ledger.Put(ctx, ledger)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) AddRecord(ctx context.Context, storage *[]InputData) error {
	for _, row := range *storage {
		entity, err := structToMap(row)
		if err != nil {
			s.logger.Error("failed to encode struct into map: ", err)
			return err
		}

		_, err = s.odbStore.Put(ctx, entity)
		if err != nil {
			s.logger.Error("failed to add data into odb: ", err)
			return err
		}
	}

	return nil
}

func (s *Storage) GetRecordByID(ctx context.Context, id string) (*InputData, error) {
	record, err := s.odbStore.Get(ctx, id, &iface.DocumentStoreGetOptions{CaseInsensitive: false})
	if err != nil {
		return nil, err
	}

	var data InputData
	for _, row := range record {
		err = mapstructure.Decode(row, &data)
		if err != nil {
			s.logger.Error("failed to decode map into struct: ", err)
			return nil, err
		}
	}

	return &data, nil
}

func (s *Storage) GetWalletByDatasetKey(ctx context.Context, key string) (*Ledger, error) {
	record, err := s.ledger.Get(ctx, key, &iface.DocumentStoreGetOptions{CaseInsensitive: false})
	if err != nil {
		return nil, err
	}

	var data Ledger
	for _, row := range record {
		err = mapstructure.Decode(row, &data)
		if err != nil {
			s.logger.Error("failed to decode map into struct: ", err)
			return nil, err
		}
	}

	return &data, nil
}

func structToMap(v interface{}) (map[string]interface{}, error) {
	vMap := &map[string]interface{}{}

	err := mapstructure.Decode(v, &vMap)
	if err != nil {
		return nil, err
	}

	return *vMap, nil
}
