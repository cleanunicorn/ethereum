package types

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
)

type Account struct {
	Key ecdsa.PrivateKey
}

// Address returns the address as a string
func (a *Account) Address() string {
	return crypto.PubkeyToAddress(a.Key.PublicKey).Hex()
}

// PrivateKey returns the private key as a string
func (a *Account) PrivateKey() string {
	return hex.EncodeToString(a.Key.D.Bytes())
}

func AccountFromHexKey(privateKey string) (Account, error) {
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return Account{}, err
	}

	return Account{
		Key: *key,
	}, nil
}
