package calc

import "math/big"

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func MinBigInt(x, y *big.Int) *big.Int {
	if x.Cmp(y) <= 0 {
		return new(big.Int).Set(x)
	}
	return new(big.Int).Set(y)
}
