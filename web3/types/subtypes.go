package types

import (
	"strconv"
)

type ComplexNumber string

func (c ComplexNumber) Int64() int64 {
	n, _ := strconv.ParseInt(string(c), 0, 64)
	return n
}
