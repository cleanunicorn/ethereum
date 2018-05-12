package account

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
)

// Account entity
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

// FromHexKey generates an Account from a string version of a private key
func FromHexKey(privateKey string) (Account, error) {
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return Account{}, err
	}

	return Account{
		Key: *key,
	}, nil
}
