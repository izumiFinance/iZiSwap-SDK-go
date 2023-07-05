package swap

import (
	"math/big"
)

type SwapAmount struct {
	AmountX big.Int
	AmountY big.Int
}

type PoolInfo struct {
	CurrentPoint int
	PointDelta   int
	LeftMostPt   int
	RightMostPt  int
	Fee          int
	Liquidity    big.Int
	LiquidityX   big.Int
	Liquidities  []LiquidityPoint
	LimitOrders  []LimitOrderPoint
}
