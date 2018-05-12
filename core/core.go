package core

import (
	"fmt"
	"math/big"

	"github.com/cleanunicorn/ethereum/web3/account"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
)

// CreateWallet creates a wallet and outputs the account data
func CreateWallet() (account.Account, error) {
	log.Debug("Creating wallet")

	key, err := crypto.GenerateKey()
	if err != nil {
		return account.Account{}, fmt.Errorf("Could not generate private key %s", err)
	}

	a := account.Account{
		Key: *key,
	}

	return a, nil
}

// CreateSigner creates a signer specific for the network as specified in EIP155
// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-155.md
func CreateSigner(network int64) gethtypes.EIP155Signer {
	return gethtypes.NewEIP155Signer(big.NewInt(network))
}

// SignTx uses the account and the rest of the parameters to sign a transaction and return the signed transaction
func SignTx(
	signer gethtypes.EIP155Signer,
	account account.Account,
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
