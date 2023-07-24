package swap

import (
	"math/big"
	"testing"
)

func getLiquiditiesDetailX2Y() []LiquidityPoint {
	liquidities := []LiquidityPoint{
		{LiqudityDelta: *big.NewInt(200000), Point: -7000},
		{LiqudityDelta: *big.NewInt(300000), Point: -5000},
		{LiqudityDelta: *big.NewInt(-300000), Point: -2000},
		{LiqudityDelta: *big.NewInt(-200000), Point: -240},

		{LiqudityDelta: *big.NewInt(600000), Point: -200},
		{LiqudityDelta: *big.NewInt(-600000), Point: 40},
		{LiqudityDelta: *big.NewInt(500000), Point: 80},
		{LiqudityDelta: *big.NewInt(-500000), Point: 2000},
	}
	return liquidities
}

func getLimitOrdersDetailX2Y() []LimitOrderPoint {
	limitOrders := []LimitOrderPoint{
		{SellingY: *big.NewInt(100000000000), Point: -6200},
		{SellingY: *big.NewInt(150000000000), Point: -1000},
		{SellingY: *big.NewInt(120000000000), Point: 1200},
		{SellingX: *big.NewInt(120000000000), Point: 1800},
	}
	return limitOrders
}

func getPoolInfoDetailX2Y() PoolInfo {
	return PoolInfo{
		CurrentPoint: 1729,
		PointDelta:   40,
		LeftMostPt:   -800000,
		RightMostPt:  800000,
		Fee:          2000,
		// other test case may change following
		// liquidity and liquidityX value
		Liquidity:   *big.NewInt(500000),
		LiquidityX:  *big.NewInt(500000),
		Liquidities: getLiquiditiesDetailX2Y(),
		LimitOrders: getLimitOrdersDetailX2Y(),
	}
}

func TestSwapDetailX2Y1(t *testing.T) {
	// x2y start partial y-liquidity,
	// end partial y-liquidity
	poolInfo := getPoolInfoDetailX2Y()
	poolInfo.CurrentPoint = 1729
	poolInfo.Liquidity = *big.NewInt(500000)
	poolInfo.LiquidityX = *big.NewInt(134333)
	var amount big.Int
	amount.SetString("462592000000", 10)
	lowPt := -6789
	swapResult, _ := SwapX2Y(amount, lowPt, poolInfo)
	costX, _ := new(big.Int).SetString("462592000000", 10)
	acquireY, _ := new(big.Int).SetString("372866052521", 10)
	finalPoint := -6786
	if swapResult.AmountX.Cmp(costX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapResult.AmountX.String(), costX.String())
	}
	if swapResult.AmountY.Cmp(acquireY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapResult.AmountY.String(), acquireY.String())
	}
	if swapResult.CurrentPoint != finalPoint {
		t.Fatalf("result currentPoint not equal (%d, %d)", swapResult.CurrentPoint, finalPoint)
	}
	resultLiquidity, _ := new(big.Int).SetString("200000", 10)
	resultLiquidityX, _ := new(big.Int).SetString("151638", 10)
	if swapResult.Liquidity.Cmp(resultLiquidity) != 0 {
		t.Fatalf("Liquidity not equal (%s, %s)", swapResult.Liquidity.String(), resultLiquidity.String())
	}
	if swapResult.LiquidityX.Cmp(resultLiquidityX) != 0 {
		t.Fatalf("LiquidityX not equal (%s, %s)", swapResult.LiquidityX.String(), resultLiquidityX.String())
	}
}

func TestSwapDetailX2Y2(t *testing.T) {
	// gocalc x2y start partial liquidity
	// end full liquidity
	poolInfo := getPoolInfoDetailX2Y()
	poolInfo.CurrentPoint = 1729
	poolInfo.Liquidity = *big.NewInt(500000)
	poolInfo.LiquidityX = *big.NewInt(134333)
	var amount big.Int
	amount.SetString("100000000000000000", 10)
	lowPt := -6789
	swapResult, _ := SwapX2Y(amount, lowPt, poolInfo)
	costX, _ := new(big.Int).SetString("462592912170", 10)
	acquireY, _ := new(big.Int).SetString("372866514294", 10)
	finalPoint := -6789
	if swapResult.AmountX.Cmp(costX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapResult.AmountX.String(), costX.String())
	}
	if swapResult.AmountY.Cmp(acquireY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapResult.AmountY.String(), acquireY.String())
	}
	if swapResult.CurrentPoint != finalPoint {
		t.Fatalf("result currentPoint not equal (%d, %d)", swapResult.CurrentPoint, finalPoint)
	}
	resultLiquidity, _ := new(big.Int).SetString("200000", 10)
	resultLiquidityX, _ := new(big.Int).SetString("200000", 10)
	if swapResult.Liquidity.Cmp(resultLiquidity) != 0 {
		t.Fatalf("Liquidity not equal (%s, %s)", swapResult.Liquidity.String(), resultLiquidity.String())
	}
	if swapResult.LiquidityX.Cmp(resultLiquidityX) != 0 {
		t.Fatalf("LiquidityX not equal (%s, %s)", swapResult.LiquidityX.String(), resultLiquidityX.String())
	}
}

func TestSwapDetailX2Y3(t *testing.T) {
	// x2y start with full-x-liquidity
	// end full liquidity
	poolInfo := getPoolInfoDetailX2Y()
	poolInfo.CurrentPoint = 1731
	poolInfo.Liquidity = *big.NewInt(500000)
	poolInfo.LiquidityX = *big.NewInt(500000)
	var amount big.Int
	amount.SetString("100000000000000000", 10)
	lowPt := -6789
	swapResult, _ := SwapX2Y(amount, lowPt, poolInfo)
	costX, _ := new(big.Int).SetString("462593495113", 10)
	acquireY, _ := new(big.Int).SetString("372867205930", 10)
	finalPoint := -6789
	if swapResult.AmountX.Cmp(costX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapResult.AmountX.String(), costX.String())
	}
	if swapResult.AmountY.Cmp(acquireY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapResult.AmountY.String(), acquireY.String())
	}
	if swapResult.CurrentPoint != finalPoint {
		t.Fatalf("result currentPoint not equal (%d, %d)", swapResult.CurrentPoint, finalPoint)
	}
	resultLiquidity, _ := new(big.Int).SetString("200000", 10)
	resultLiquidityX, _ := new(big.Int).SetString("200000", 10)
	if swapResult.Liquidity.Cmp(resultLiquidity) != 0 {
		t.Fatalf("Liquidity not equal (%s, %s)", swapResult.Liquidity.String(), resultLiquidity.String())
	}
	if swapResult.LiquidityX.Cmp(resultLiquidityX) != 0 {
		t.Fatalf("LiquidityX not equal (%s, %s)", swapResult.LiquidityX.String(), resultLiquidityX.String())
	}
}

func TestSwapDetailX2Y4(t *testing.T) {
	// x2y start with full-y-liquidity
	// end full liquidity
	poolInfo := getPoolInfoDetailX2Y()
	poolInfo.CurrentPoint = 1729
	poolInfo.Liquidity = *big.NewInt(500000)
	poolInfo.LiquidityX = *big.NewInt(0)
	var amount big.Int
	amount.SetString("100000000000000000", 10)
	lowPt := -6789
	swapResult, _ := SwapX2Y(amount, lowPt, poolInfo)
	costX, _ := new(big.Int).SetString("462593035624", 10)
	acquireY, _ := new(big.Int).SetString("372866660757", 10)
	finalPoint := -6789
	if swapResult.AmountX.Cmp(costX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapResult.AmountX.String(), costX.String())
	}
	if swapResult.AmountY.Cmp(acquireY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapResult.AmountY.String(), acquireY.String())
	}
	if swapResult.CurrentPoint != finalPoint {
		t.Fatalf("result currentPoint not equal (%d, %d)", swapResult.CurrentPoint, finalPoint)
	}
	resultLiquidity, _ := new(big.Int).SetString("200000", 10)
	resultLiquidityX, _ := new(big.Int).SetString("200000", 10)
	if swapResult.Liquidity.Cmp(resultLiquidity) != 0 {
		t.Fatalf("Liquidity not equal (%s, %s)", swapResult.Liquidity.String(), resultLiquidity.String())
	}
	if swapResult.LiquidityX.Cmp(resultLiquidityX) != 0 {
		t.Fatalf("LiquidityX not equal (%s, %s)", swapResult.LiquidityX.String(), resultLiquidityX.String())
	}
}

func TestSwapDetailX2Y5(t *testing.T) {
	// x2y start with limitorder-y
	// full-x-liquidity, end partial liquidity
	poolInfo := getPoolInfoDetailX2Y()
	poolInfo.CurrentPoint = 1200
	poolInfo.Liquidity = *big.NewInt(500000)
	poolInfo.LiquidityX = *big.NewInt(500000)
	var amount big.Int
	amount.SetString("462346200000", 10)
	lowPt := -6789
	swapResult, _ := SwapX2Y(amount, lowPt, poolInfo)
	costX, _ := new(big.Int).SetString("462346200000", 10)
	acquireY, _ := new(big.Int).SetString("372581497872", 10)
	finalPoint := -6789
	if swapResult.AmountX.Cmp(costX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapResult.AmountX.String(), costX.String())
	}
	if swapResult.AmountY.Cmp(acquireY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapResult.AmountY.String(), acquireY.String())
	}
	if swapResult.CurrentPoint != finalPoint {
		t.Fatalf("result currentPoint not equal (%d, %d)", swapResult.CurrentPoint, finalPoint)
	}
	resultLiquidity, _ := new(big.Int).SetString("200000", 10)
	resultLiquidityX, _ := new(big.Int).SetString("167919", 10)
	if swapResult.Liquidity.Cmp(resultLiquidity) != 0 {
		t.Fatalf("Liquidity not equal (%s, %s)", swapResult.Liquidity.String(), resultLiquidity.String())
	}
	if swapResult.LiquidityX.Cmp(resultLiquidityX) != 0 {
		t.Fatalf("LiquidityX not equal (%s, %s)", swapResult.LiquidityX.String(), resultLiquidityX.String())
	}
}

func TestSwapDetailX2Y6(t *testing.T) {
	// x2y start with limitorder-x
	// full-y-liquidity, end partial liquidity
	poolInfo := getPoolInfoDetailX2Y()

	poolInfo.LimitOrders[2].SellingX = *big.NewInt(120000000000)
	poolInfo.LimitOrders[2].SellingY = *big.NewInt(0)

	poolInfo.CurrentPoint = 1200
	poolInfo.Liquidity = *big.NewInt(500000)
	poolInfo.LiquidityX = *big.NewInt(0)

	var amount big.Int
	amount.SetString("355702300000", 10)
	lowPt := -6789
	swapResult, _ := SwapX2Y(amount, lowPt, poolInfo)
	costX, _ := new(big.Int).SetString("355702300000", 10)
	acquireY, _ := new(big.Int).SetString("252582032779", 10)
	finalPoint := -6789
	if swapResult.AmountX.Cmp(costX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapResult.AmountX.String(), costX.String())
	}
	if swapResult.AmountY.Cmp(acquireY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapResult.AmountY.String(), acquireY.String())
	}
	if swapResult.CurrentPoint != finalPoint {
		t.Fatalf("result currentPoint not equal (%d, %d)", swapResult.CurrentPoint, finalPoint)
	}
	resultLiquidity, _ := new(big.Int).SetString("200000", 10)
	resultLiquidityX, _ := new(big.Int).SetString("173521", 10)
	if swapResult.Liquidity.Cmp(resultLiquidity) != 0 {
		t.Fatalf("Liquidity not equal (%s, %s)", swapResult.Liquidity.String(), resultLiquidity.String())
	}
	if swapResult.LiquidityX.Cmp(resultLiquidityX) != 0 {
		t.Fatalf("LiquidityX not equal (%s, %s)", swapResult.LiquidityX.String(), resultLiquidityX.String())
	}
}

func TestSwapDetailX2Y7(t *testing.T) {
	// x2y start with limitorder-y, partial-liquidity,
	// end partial liquidity
	poolInfo := getPoolInfoDetailX2Y()
	poolInfo.CurrentPoint = 1200
	poolInfo.Liquidity = *big.NewInt(500000)
	poolInfo.LiquidityX = *big.NewInt(383966)
	var amount big.Int
	amount.SetString("462346300000", 10)
	lowPt := -6789
	swapResult, _ := SwapX2Y(amount, lowPt, poolInfo)
	costX, _ := new(big.Int).SetString("462346300000", 10)
	acquireY, _ := new(big.Int).SetString("372581616273", 10)
	finalPoint := -6789
	if swapResult.AmountX.Cmp(costX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapResult.AmountX.String(), costX.String())
	}
	if swapResult.AmountY.Cmp(acquireY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapResult.AmountY.String(), acquireY.String())
	}
	if swapResult.CurrentPoint != finalPoint {
		t.Fatalf("result currentPoint not equal (%d, %d)", swapResult.CurrentPoint, finalPoint)
	}
	resultLiquidity, _ := new(big.Int).SetString("200000", 10)
	resultLiquidityX, _ := new(big.Int).SetString("161169", 10)
	if swapResult.Liquidity.Cmp(resultLiquidity) != 0 {
		t.Fatalf("Liquidity not equal (%s, %s)", swapResult.Liquidity.String(), resultLiquidity.String())
	}
	if swapResult.LiquidityX.Cmp(resultLiquidityX) != 0 {
		t.Fatalf("LiquidityX not equal (%s, %s)", swapResult.LiquidityX.String(), resultLiquidityX.String())
	}
}

func TestSwapDetailX2Y8(t *testing.T) {
	// x2y start with limitorder-y, partial-liquidity
	// end full-x-liquidityand partial or full limit order
	poolInfo := getPoolInfoDetailX2Y()
	poolInfo.CurrentPoint = 1200
	poolInfo.Liquidity = *big.NewInt(500000)
	poolInfo.LiquidityX = *big.NewInt(383966)
	var amount big.Int
	amount.SetString("325923465573", 10)
	lowPt := -6201
	swapResult, _ := SwapX2Y(amount, lowPt, poolInfo)
	costX, _ := new(big.Int).SetString("325923465573", 10)
	acquireY, _ := new(big.Int).SetString("299340764005", 10)
	finalPoint := -6200
	if swapResult.AmountX.Cmp(costX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapResult.AmountX.String(), costX.String())
	}
	if swapResult.AmountY.Cmp(acquireY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapResult.AmountY.String(), acquireY.String())
	}
	if swapResult.CurrentPoint != finalPoint {
		t.Fatalf("result currentPoint not equal (%d, %d)", swapResult.CurrentPoint, finalPoint)
	}
	resultLiquidity, _ := new(big.Int).SetString("200000", 10)
	resultLiquidityX, _ := new(big.Int).SetString("200000", 10)
	if swapResult.Liquidity.Cmp(resultLiquidity) != 0 {
		t.Fatalf("Liquidity not equal (%s, %s)", swapResult.Liquidity.String(), resultLiquidity.String())
	}
	if swapResult.LiquidityX.Cmp(resultLiquidityX) != 0 {
		t.Fatalf("LiquidityX not equal (%s, %s)", swapResult.LiquidityX.String(), resultLiquidityX.String())
	}
}

func TestSwapDetailX2Y9(t *testing.T) {
	// gocalc x2y start with limitorder-y, partial-liquidity
	// end partial-x-liquidity
	poolInfo := getPoolInfoDetailX2Y()
	poolInfo.CurrentPoint = 1200
	poolInfo.Liquidity = *big.NewInt(500000)
	poolInfo.LiquidityX = *big.NewInt(383966)
	var amount big.Int
	amount.SetString("275923400000", 10)
	lowPt := -6200
	swapResult, _ := SwapX2Y(amount, lowPt, poolInfo)
	costX, _ := new(big.Int).SetString("275923400000", 10)
	acquireY, _ := new(big.Int).SetString("272496469260", 10)
	finalPoint := -6200
	if swapResult.AmountX.Cmp(costX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapResult.AmountX.String(), costX.String())
	}
	if swapResult.AmountY.Cmp(acquireY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapResult.AmountY.String(), acquireY.String())
	}
	if swapResult.CurrentPoint != finalPoint {
		t.Fatalf("result currentPoint not equal (%d, %d)", swapResult.CurrentPoint, finalPoint)
	}
	resultLiquidity, _ := new(big.Int).SetString("200000", 10)
	resultLiquidityX, _ := new(big.Int).SetString("152001", 10)
	if swapResult.Liquidity.Cmp(resultLiquidity) != 0 {
		t.Fatalf("Liquidity not equal (%s, %s)", swapResult.Liquidity.String(), resultLiquidity.String())
	}
	if swapResult.LiquidityX.Cmp(resultLiquidityX) != 0 {
		t.Fatalf("LiquidityX not equal (%s, %s)", swapResult.LiquidityX.String(), resultLiquidityX.String())
	}
}
