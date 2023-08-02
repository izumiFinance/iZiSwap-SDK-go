package swap

import (
	"math/big"
)

type SwapResult struct {
	AmountX      *big.Int
	AmountY      *big.Int
	CurrentPoint int
	Liquidity    *big.Int
	LiquidityX   *big.Int
}

type PoolInfo struct {
	CurrentPoint int
	PointDelta   int
	LeftMostPt   int
	RightMostPt  int
	Fee          int
	Liquidity    *big.Int
	LiquidityX   *big.Int
	Liquidities  []LiquidityPoint
	LimitOrders  []LimitOrderPoint
}
