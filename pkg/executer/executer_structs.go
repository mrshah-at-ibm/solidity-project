package executer

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mrshah-at-ibm/kaleido-project/pkg/mrstoken"
	"github.com/mrshah-at-ibm/kaleido-project/pkg/signer"
	"go.uber.org/zap"
	// "golang.org/x/sync/semaphore"
)

type Executer struct {
	Client          *ethclient.Client
	Logger          *zap.SugaredLogger
	Address         common.Address
	Contract        *mrstoken.MRSToken
	ContractAddress string
	Signer          *signer.Signer
	Nonce           uint64
	// Semaphore       *semaphore.Weighted
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ExecuterInterface
type ExecuterInterface interface {
	DeployContract() error
	MintToken(to string) (*types.Receipt, error)
	BurnToken(tokenid string) (*types.Receipt, error)
	TransferToken(from string, to string, tokenid string) (*types.Receipt, error)
	BalanceOf(owner string) (int, error)
	GetBalance(owner string) (int, error)
}
