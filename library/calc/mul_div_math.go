package calc

import "math/big"

// MulDivFloor performs multiplication first and then division, flooring the result.
func MulDivFloor(a, b, c *big.Int) *big.Int {
	mul := new(big.Int).Mul(a, b)
	return new(big.Int).Div(mul, c)
}

// MulDivCeil performs multiplication first and then division, ceiling the result.
func MulDivCeil(a, b, c *big.Int) *big.Int {
	mul := new(big.Int).Mul(a, b)
	sum := new(big.Int).Add(mul, new(big.Int).Sub(c, big.NewInt(1)))
	return new(big.Int).Div(sum, c)
}
