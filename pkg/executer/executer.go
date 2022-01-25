package executer

import (
	"context"
	"errors"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mrshah-at-ibm/kaleido-project/pkg/config"
	"github.com/mrshah-at-ibm/kaleido-project/pkg/mrstoken"
	"github.com/mrshah-at-ibm/kaleido-project/pkg/signer"
	"go.uber.org/zap"
)

func NewExecuter(l *zap.Logger) (*Executer, error) {
	e := &Executer{
		Logger: l.Sugar().Named("Executer"),
	}

	conf, err := config.ReadConfig()
	if err != nil {
		e.Logger.Error(err)
		return nil, err
	}
	backend, err := ethclient.Dial(conf.NodeUrl)
	if err != nil {
		e.Logger.Error(err)
		return nil, err
	}

	e.Client = backend

	addrClaimed, err := config.ClaimAddress()
	if err != nil {
		e.Logger.Error(err)
		return nil, err
	}

	addr := common.HexToAddress(addrClaimed[0])
	e.Address = addr

	key, err := EnsurePrivateKey(addrClaimed[0])
	if err != nil {
		e.Logger.Error(err)
		return nil, err
	}

	var signer1 types.EIP155Signer
	chainid, err := backend.ChainID(context.TODO())
	if err != nil {
		e.Logger.Error(err)
		return nil, err
	}
	signer1 = types.NewEIP155Signer(chainid)
	s := signer.NewSigner(signer1, key)

	e.Signer = s

	contractaddress, err := config.ReadContractAddress("default")
	e.ContractAddress = contractaddress

	return e, nil
}

func (e *Executer) DeployContract() error {
	if e.ContractAddress == "" {
		transactOpts := &bind.TransactOpts{
			From: e.Address,
			// Nonce:     "",
			Signer: e.Signer.SignTx,
			// Value:     "",
			// GasPrice:  "",
			// GasFeeCap: "",
			// GasTipCap: "",
			GasLimit: 100000000,
			Context:  context.TODO(),
			// NoSend:    "",
		}
		addr, tx, crt, err := mrstoken.DeployMRSToken(transactOpts, e.Client, "mrshah token", "MRS", "www.mrshah.space")
		if err != nil {
			e.Logger.Error(err)
			return err
		}

		ctx, _ := context.WithTimeout(context.TODO(), 10*time.Second)
		e.Logger.Info("Waiting for contract to be deployed ---")
		_, err = bind.WaitDeployed(ctx, e.Client, tx)
		if err != nil {
			e.Logger.Errorw("Wait deployed", "error", err)
			return err
		}

		e.Logger.Info("Waiting for tx to be mined")
		_, err = bind.WaitMined(context.TODO(), e.Client, tx)
		if err != nil {
			e.Logger.Error(err)
			return err
		}

		e.Contract = crt
		err = config.SaveContractAddress("default", addr.String(), "")
		if err != nil {
			e.Logger.Error(err)
			return err
		}
		e.Logger.Info("Contract deployed successfully")

	} else {
		e.Logger.Info("Contract is already deployed")
		var err error
		e.Contract, err = mrstoken.NewMRSToken(common.HexToAddress(e.ContractAddress), e.Client)
		if err != nil {
			e.Logger.Error(err)
			return err
		}
	}
	return nil
}

func (e *Executer) MintToken(to string) (*types.Receipt, error) {
	if e.Contract == nil {
		e.Logger.Error("Contract is not initialized")
		return nil, errors.New("Contract is not initialized")
	}

	nonce, err := e.Client.PendingNonceAt(context.TODO(), e.Address)
	if err != nil {
		e.Logger.Error(err)
		return nil, err
	}

	e.Logger.Infow("PendingNonce", "nonce", nonce)
	transactOpts := &bind.TransactOpts{
		From:   e.Address,
		Nonce:  big.NewInt(int64(e.Nonce)),
		Signer: e.Signer.SignTx,
		// Value:     "",
		// GasPrice:  "",
		// GasFeeCap: "",
		// GasTipCap: "",
		GasLimit: 1000000,
		Context:  context.TODO(),
		// NoSend:    "",
	}

	e.Logger.Infow("Sending Transaction", "transaction", transactOpts, "nonce", transactOpts.Nonce, "account", transactOpts.From, "to", to)
	tr, err := e.Contract.Mint(transactOpts, common.HexToAddress(to))
	e.Nonce++
	if err != nil {
		if strings.Contains(err.Error(), "nonce too low") {
			e.MintToken(to)
		}

		e.Logger.Error(err)
		return nil, err
	}

	receipt, err := bind.WaitMined(context.TODO(), e.Client, tr)
	if err != nil {
		e.Logger.Error(err)
		return nil, err
	}

	e.Logger.Infow("Transaction", "transaction", tr, "receipt", receipt)
	return receipt, nil
}

func (e *Executer) BurnToken(tokenid string) (*types.Receipt, error) {
	if e.Contract == nil {
		e.Logger.Error("Contract is not initialized")
		return nil, errors.New("Contract is not initialized")
	}

	transactOpts := &bind.TransactOpts{
		From:   e.Address,
		Nonce:  big.NewInt(int64(e.Nonce)),
		Signer: e.Signer.SignTx,
		// Value:     "",
		// GasPrice:  "",
		// GasFeeCap: "",
		// GasTipCap: "",
		GasLimit: 1000000,
		Context:  context.TODO(),
		// NoSend:    "",
	}

	tokenint, err := strconv.Atoi(tokenid)
	if err != nil {
		if err.Error() == "nonce too low" {
			e.BurnToken(tokenid)
		}
		e.Logger.Error(err)
		return nil, err
	}

	tr, err := e.Contract.Burn(transactOpts, big.NewInt(int64(tokenint)))
	if err != nil {
		e.Logger.Error(err)
		return nil, err
	}

	if err != nil {
		// TODO: Check if following hits
		if strings.Contains(err.Error(), "nonce too low") {
			return e.BurnToken(tokenid)
		}

		e.Logger.Error(err)
		return nil, err
	}

	e.Logger.Infow("Transaction", "transaction", tr)

	receipt, err := bind.WaitMined(context.TODO(), e.Client, tr)
	if err != nil {
		e.Logger.Error(err)
		return nil, err
	}

	e.Logger.Infow("Transaction", "transaction", tr, "receipt", receipt)

	return receipt, nil

}

func (e *Executer) TransferToken(from string, to string, tokenid string) (*types.Receipt, error) {
	if e.Contract == nil {
		e.Logger.Error("Contract is not initialized")
		return nil, errors.New("Contract is not initialized")
	}

	transactOpts := &bind.TransactOpts{
		From:   e.Address,
		Nonce:  big.NewInt(int64(e.Nonce)),
		Signer: e.Signer.SignTx,
		// Value:     "",
		// GasPrice:  "",
		// GasFeeCap: "",
		// GasTipCap: "",
		GasLimit: 1000000,
		Context:  context.TODO(),
		// NoSend:    "",
	}

	tokenint, err := strconv.Atoi(tokenid)
	if err != nil {
		if err.Error() == "nonce too low" {
			e.BurnToken(tokenid)
		}
		e.Logger.Error(err)
		return nil, err
	}

	tr, err := e.Contract.TransferFrom(transactOpts, common.HexToAddress(from), common.HexToAddress(to), big.NewInt(int64(tokenint)))
	if err != nil {
		e.Logger.Error(err)
		return nil, err
	}
	e.Logger.Infow("Transaction", "transaction", tr)

	if err != nil {
		// TODO: Check if following hits
		if strings.Contains(err.Error(), "nonce too low") {
			return e.TransferToken(from, to, tokenid)
		}

		e.Logger.Error(err)
		return nil, err
	}

	e.Logger.Infow("Transaction", "transaction", tr)

	receipt, err := bind.WaitMined(context.TODO(), e.Client, tr)
	if err != nil {
		e.Logger.Error(err)
		return nil, err
	}

	e.Logger.Infow("Transaction", "transaction", tr, "receipt", receipt)

	return receipt, nil
}

func (e *Executer) BalanceOf(owner string) (int, error) {
	if e.Contract == nil {
		e.Logger.Error("Contract is not initialized")
		return 0, errors.New("contract is not initialized")
	}

	transactOpts := &bind.CallOpts{
		Pending:     true,
		BlockNumber: nil,
		From:        e.Address,
		Context:     context.TODO(),
	}

	balance, err := e.Contract.CheckBalanceOf(transactOpts, common.HexToAddress(owner))
	if err != nil {
		e.Logger.Error(err)
		return 0, err
	}
	e.Logger.Infow("Balance Value", "balance", balance)
	return int(balance.Int64()), nil
}

func (e *Executer) GetBalance(owner string) (int, error) {
	// if e.Contract == nil {
	// 	e.Logger.Error("Contract is not initialized")
	// 	return 0, errors.New("Contract is not initialized")
	// }

	// transactOpts := &bind.CallOpts{
	// 	Pending:     true,
	// 	BlockNumber: nil,
	// 	From:        e.Address,
	// 	Context:     context.TODO(),
	// }

	balance, err := e.Client.BalanceAt(context.TODO(), e.Address, nil)
	// balance, err := e.Contract.BalanceOf(transactOpts, common.HexToAddress(owner))
	if err != nil {
		e.Logger.Error(err)
		return 0, err
	}
	e.Logger.Infow("Balance Value", "balance", balance)
	return int(balance.Int64()), nil
}
