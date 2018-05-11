package helper

import (
	"math/big"
	"strings"
)

func Trim0x(prefixedString string) string {
	return strings.TrimPrefix(prefixedString, "0x")
}

func HexStrToBigInt(hexString string) *big.Int {
	value := new(big.Int)
	value.SetString(Trim0x(hexString), 16)
	return value
}
