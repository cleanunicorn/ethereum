package client

import (
	"math/big"
	"strings"
)

func trim0x(PrefixedString string) string {
	return strings.TrimPrefix(PrefixedString, "0x")
}

func hexStrToBigInt(HexString string) *big.Int {
	value := new(big.Int)
	value.SetString(trim0x(HexString), 16)
	return value
}
