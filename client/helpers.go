package client

import (
	"math/big"
	"strings"
)

func trim0x(prefixedString string) string {
	return strings.TrimPrefix(prefixedString, "0x")
}

func hexStrToBigInt(hexString string) *big.Int {
	value := new(big.Int)
	value.SetString(trim0x(hexString), 16)
	return value
}
