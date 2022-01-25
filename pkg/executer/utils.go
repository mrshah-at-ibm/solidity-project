package executer

import (
	"crypto/ecdsa"
	"fmt"

	ecrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/mrshah-at-ibm/kaleido-project/pkg/config"
)

func EnsurePrivateKey(account string) (*ecdsa.PrivateKey, error) {

	var key *ecdsa.PrivateKey
	var err error

	// err := e.ReadWriteSemaphore.Acquire(context.TODO(), 1)
	// if err != nil {
	// 	return nil, err
	// }

	// 	defer e.ReadWriteSemaphore.Release(1)

	privateKey, err := config.ReadPrivateKey(account)
	if err != nil {
		return nil, err
	}
	if privateKey == nil {
		key, err := ecrypto.GenerateKey()
		if err != nil {
			return nil, fmt.Errorf("generating key: %s", err)
		}

		err = config.SavePrivateKey(account, key)
		if err != nil {
			return nil, err
		}

	} else {
		key = privateKey
	}

	return key, err
}
