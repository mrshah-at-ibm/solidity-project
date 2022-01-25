package signer

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Signer struct {
	PrivateKey *ecdsa.PrivateKey
	Signer     types.EIP155Signer
}

func NewSigner(signer types.EIP155Signer, privatekey *ecdsa.PrivateKey) *Signer {
	return &Signer{
		PrivateKey: privatekey,
		Signer:     signer,
	}
}

func (s *Signer) SignTx(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
	signedTx, err := types.SignTx(tx, s.Signer, s.PrivateKey)

	return signedTx, err
}
