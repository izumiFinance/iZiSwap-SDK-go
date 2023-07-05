package orders

import "math/big"

type Liquidity struct {
	LiqudityDelta big.Int
	Point         int
}

type LimitOrder struct {
	SellingX big.Int
	SellingY big.Int
	Point    int
}

type OrderData struct {
	Liquidities   []Liquidity
	LiquidityIdx  int
	LimitOrders   []LimitOrder
	LimitOrderIdx int
}

func (orderData *OrderData) IsLiquidity(point int) bool {
	if orderData.LiquidityIdx < 0 || orderData.LiquidityIdx >= len(orderData.Liquidities) {
		return false
	}
	return orderData.Liquidities[orderData.LiquidityIdx].Point == point
}

func (orderData *OrderData) IsLimitOrder(point int) bool {
	if orderData.LimitOrderIdx < 0 || orderData.LimitOrderIdx >= len(orderData.LimitOrders) {
		return false
	}
	return orderData.LimitOrders[orderData.LimitOrderIdx].Point == point
}

func (orderData *OrderData) UnsafeGetDeltaLiquidity() big.Int {
	return orderData.Liquidities[orderData.LiquidityIdx].LiqudityDelta
}

func (orderData *OrderData) UnsafeGetLimitSellingX() big.Int {
	return orderData.LimitOrders[orderData.LimitOrderIdx].SellingX
}

func (orderData *OrderData) UnsafeGetLimitSellingY() big.Int {
	return orderData.LimitOrders[orderData.LimitOrderIdx].SellingY
}
