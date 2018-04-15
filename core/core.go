package core

import (
	"fmt"
	"math/big"

	"github.com/cleanunicorn/ethereum/types"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
)

// CreateWallet creates a wallet and outputs the account data
func CreateWallet() (types.Account, error) {
	log.Debug("Creating wallet")

	key, err := crypto.GenerateKey()
	if err != nil {
		return types.Account{}, fmt.Errorf("Could not generate private key %s", err)
	}

	a := types.Account{
		Key: *key,
	}

	return a, nil
}

// CreateSigner creates a signer specific to the network
func CreateSigner(network int64) gethtypes.EIP155Signer {
	return gethtypes.NewEIP155Signer(big.NewInt(network))
}

// SignTx uses the account and the rest of the parameters to sign a transaction and return the signed transaction
func SignTx(
	signer gethtypes.EIP155Signer,
	account types.Account,
	nonce uint64,
	to common.Address,
	amount *big.Int,
	gasLimit uint64,
	gasPrice *big.Int,
	data []byte,
) (*gethtypes.Transaction, error) {

	tx, err := gethtypes.SignTx(
		gethtypes.NewTransaction(
			nonce,
			to,
			amount,
			gasLimit,
			gasPrice,
			data,
		),
		signer,
		&account.Key,
	)
	if err != nil {
		return &gethtypes.Transaction{}, err
	}

	return tx, nil
}
