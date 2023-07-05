package utils

import "math/big"

type State struct {
	LiquidityX   *big.Int
	Liquidity    *big.Int
	CurrentPoint int
	SqrtPrice_96 *big.Int
}
