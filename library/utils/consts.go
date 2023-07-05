package utils

import "math/big"

var Pow96 = new(big.Int).Exp(big.NewInt(2), big.NewInt(96), nil)
