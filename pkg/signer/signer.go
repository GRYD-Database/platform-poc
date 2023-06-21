package signer

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Signer interface {
	// EthereumAddress returns the ethereum address this signer uses.
	EthereumAddress() (common.Address, error)
}

type defaultSigner struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func (d *defaultSigner) EthereumAddress() (common.Address, error) {
	return crypto.PubkeyToAddress(*d.publicKey), nil
}

func New(hexKey string) (Signer, error) {
	privateKey, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		return nil, fmt.Errorf("eror generating private key from hex: %w ", err)
	}

	publicKey := privateKey.Public()

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("unable to generate public key")
	}

	return &defaultSigner{
		privateKey: privateKey,
		publicKey:  publicKeyECDSA,
	}, nil
}
