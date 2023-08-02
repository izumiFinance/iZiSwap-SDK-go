package swap

import (
	"math/big"
	"testing"
)

func getLiquiditiesDetailY2X() []LiquidityPoint {
	liquidities := []LiquidityPoint{
		{LiqudityDelta: big.NewInt(200000), Point: -7000},
		{LiqudityDelta: big.NewInt(300000), Point: -5000},
		{LiqudityDelta: big.NewInt(-300000), Point: -2000},
		{LiqudityDelta: big.NewInt(-200000), Point: -240},

		{LiqudityDelta: big.NewInt(600000), Point: -200},
		{LiqudityDelta: big.NewInt(-600000), Point: 40},
		{LiqudityDelta: big.NewInt(500000), Point: 80},
		{LiqudityDelta: big.NewInt(-500000), Point: 2000},
	}
	return liquidities
}

func getLimitOrdersDetailY2X() []LimitOrderPoint {
	limitOrders := []LimitOrderPoint{
		{SellingY: big.NewInt(80000000000), Point: -6400},
		// some test case may change order at -6200
		{SellingX: big.NewInt(100000000000), Point: -6200},
		{SellingX: big.NewInt(150000000000), Point: -1000},
		{SellingX: big.NewInt(120000000000), Point: 1200},
	}
	return limitOrders
}

func getPoolInfoDetailY2X() PoolInfo {
	return PoolInfo{
		CurrentPoint: -6216,
		PointDelta:   40,
		LeftMostPt:   -800000,
		RightMostPt:  800000,
		Fee:          2000,
		// other test case may change following
		// liquidity and liquidityX value
		Liquidity:   big.NewInt(200000),
		LiquidityX:  big.NewInt(31891),
		Liquidities: getLiquiditiesDetailY2X(),
		LimitOrders: getLimitOrdersDetailY2X(),
	}
}

func TestSwapDetailY2X1(t *testing.T) {
	// y2x start partial x-liquidity,
	// end partial x-liquidity
	poolInfo := getPoolInfoDetailY2X()
	poolInfo.CurrentPoint = -6215
	poolInfo.Liquidity = big.NewInt(200000)
	poolInfo.LiquidityX = big.NewInt(31891)
	var amount big.Int
	amount.SetString("328168800000", 10)
	swapResult, _ := SwapY2X(&amount, 1560, poolInfo)
	costY, _ := new(big.Int).SetString("328168800000", 10)
	acquireX, _ := new(big.Int).SetString("373337423211", 10)

	if swapResult.AmountY.Cmp(costY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapResult.AmountY.String(), costY.String())
	}
	if swapResult.AmountX.Cmp(acquireX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapResult.AmountX.String(), acquireX.String())
	}
	if swapResult.CurrentPoint != 1559 {
		t.Fatalf("result currentPoint not equal (%d, %d)", swapResult.CurrentPoint, 1559)
	}
	resultLiquidity, _ := new(big.Int).SetString("500000", 10)
	resultLiquidityX, _ := new(big.Int).SetString("59052", 10)
	if swapResult.Liquidity.Cmp(resultLiquidity) != 0 {
		t.Fatalf("Liquidity not equal (%s, %s)", swapResult.Liquidity.String(), resultLiquidity.String())
	}
	if swapResult.LiquidityX.Cmp(resultLiquidityX) != 0 {
		t.Fatalf("LiquidityX not equal (%s, %s)", swapResult.LiquidityX.String(), resultLiquidityX.String())
	}
}

func TestSwapDetailY2X2(t *testing.T) {
	// y2x start partial liquidity,
	// end full liquidity
	poolInfo := getPoolInfoDetailY2X()
	poolInfo.CurrentPoint = -6215
	poolInfo.Liquidity = big.NewInt(200000)
	poolInfo.LiquidityX = big.NewInt(31891)
	var amount big.Int
	amount.SetString("1000000000000000000", 10)
	swapResult, _ := SwapY2X(&amount, 1560, poolInfo)
	costY, _ := new(big.Int).SetString("328168863966", 10)
	acquireX, _ := new(big.Int).SetString("373337477835", 10)

	if swapResult.AmountY.Cmp(costY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapResult.AmountY.String(), costY.String())
	}
	if swapResult.AmountX.Cmp(acquireX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapResult.AmountX.String(), acquireX.String())
	}
	if swapResult.CurrentPoint != 1560 {
		t.Fatalf("result currentPoint not equal (%d, %d)", swapResult.CurrentPoint, 1560)
	}
	resultLiquidity, _ := new(big.Int).SetString("500000", 10)
	resultLiquidityX, _ := new(big.Int).SetString("500000", 10)
	if swapResult.Liquidity.Cmp(resultLiquidity) != 0 {
		t.Fatalf("Liquidity not equal (%s, %s)", swapResult.Liquidity.String(), resultLiquidity.String())
	}
	if swapResult.LiquidityX.Cmp(resultLiquidityX) != 0 {
		t.Fatalf("LiquidityX not equal (%s, %s)", swapResult.LiquidityX.String(), resultLiquidityX.String())
	}
}

func TestSwapDetailY2X3(t *testing.T) {
	// y2x start with full-y-liquidity
	// end full liquidity
	poolInfo := getPoolInfoDetailY2X()
	poolInfo.CurrentPoint = -6216
	poolInfo.Liquidity = big.NewInt(200000)
	poolInfo.LiquidityX = big.NewInt(0)
	var amount big.Int
	amount.SetString("1000000000000000000", 10)
	swapResult, _ := SwapY2X(&amount, 1560, poolInfo)
	costY, _ := new(big.Int).SetString("328168987421", 10)
	acquireX, _ := new(big.Int).SetString("373337707209", 10)

	if swapResult.AmountY.Cmp(costY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapResult.AmountY.String(), costY.String())
	}
	if swapResult.AmountX.Cmp(acquireX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapResult.AmountX.String(), acquireX.String())
	}
	if swapResult.CurrentPoint != 1560 {
		t.Fatalf("result currentPoint not equal (%d, %d)", swapResult.CurrentPoint, 1560)
	}
	resultLiquidity, _ := new(big.Int).SetString("500000", 10)
	resultLiquidityX, _ := new(big.Int).SetString("500000", 10)
	if swapResult.Liquidity.Cmp(resultLiquidity) != 0 {
		t.Fatalf("Liquidity not equal (%s, %s)", swapResult.Liquidity.String(), resultLiquidity.String())
	}
	if swapResult.LiquidityX.Cmp(resultLiquidityX) != 0 {
		t.Fatalf("LiquidityX not equal (%s, %s)", swapResult.LiquidityX.String(), resultLiquidityX.String())
	}
}

func TestSwapDetailY2X4(t *testing.T) {
	// y2x start with full-x-liquidity
	// end full liquidity
	poolInfo := getPoolInfoDetailY2X()
	poolInfo.CurrentPoint = -6216
	poolInfo.Liquidity = big.NewInt(200000)
	poolInfo.LiquidityX = big.NewInt(200000)
	var amount big.Int
	amount.SetString("1000000000000000000", 10)
	swapResult, _ := SwapY2X(&amount, 1560, poolInfo)
	costY, _ := new(big.Int).SetString("328169134289", 10)
	acquireX, _ := new(big.Int).SetString("373337980108", 10)

	if swapResult.AmountY.Cmp(costY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapResult.AmountY.String(), costY.String())
	}
	if swapResult.AmountX.Cmp(acquireX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapResult.AmountX.String(), acquireX.String())
	}
	if swapResult.CurrentPoint != 1560 {
		t.Fatalf("result currentPoint not equal (%d, %d)", swapResult.CurrentPoint, 1560)
	}
	resultLiquidity, _ := new(big.Int).SetString("500000", 10)
	resultLiquidityX, _ := new(big.Int).SetString("500000", 10)
	if swapResult.Liquidity.Cmp(resultLiquidity) != 0 {
		t.Fatalf("Liquidity not equal (%s, %s)", swapResult.Liquidity.String(), resultLiquidity.String())
	}
	if swapResult.LiquidityX.Cmp(resultLiquidityX) != 0 {
		t.Fatalf("LiquidityX not equal (%s, %s)", swapResult.LiquidityX.String(), resultLiquidityX.String())
	}
}
func TestSwapDetailY2X5(t *testing.T) {
	// y2x start with limitorder-x, full-x-liquidity
	// end partial liquidity
	poolInfo := getPoolInfoDetailY2X()
	poolInfo.CurrentPoint = -6200
	poolInfo.Liquidity = big.NewInt(200000)
	poolInfo.LiquidityX = big.NewInt(200000)
	var amount big.Int
	amount.SetString("328166700000", 10)
	swapResult, _ := SwapY2X(&amount, 1560, poolInfo)
	costY, _ := new(big.Int).SetString("328166700000", 10)
	acquireX, _ := new(big.Int).SetString("373333544041", 10)

	if swapResult.AmountY.Cmp(costY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapResult.AmountY.String(), costY.String())
	}
	if swapResult.AmountX.Cmp(acquireX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapResult.AmountX.String(), acquireX.String())
	}
	if swapResult.CurrentPoint != 1559 {
		t.Fatalf("result currentPoint not equal (%d, %d)", swapResult.CurrentPoint, 1559)
	}
	resultLiquidity, _ := new(big.Int).SetString("500000", 10)
	resultLiquidityX, _ := new(big.Int).SetString("77101", 10)
	if swapResult.Liquidity.Cmp(resultLiquidity) != 0 {
		t.Fatalf("Liquidity not equal (%s, %s)", swapResult.Liquidity.String(), resultLiquidity.String())
	}
	if swapResult.LiquidityX.Cmp(resultLiquidityX) != 0 {
		t.Fatalf("LiquidityX not equal (%s, %s)", swapResult.LiquidityX.String(), resultLiquidityX.String())
	}
}

func TestSwapDetailY2X6(t *testing.T) {
	// y2x start with limitorder-y, full-x-liquidity
	// end partial liquidity
	poolInfo := getPoolInfoDetailY2X()

	poolInfo.LimitOrders[1].SellingX = big.NewInt(0)
	poolInfo.LimitOrders[1].SellingY = big.NewInt(100000000000)

	poolInfo.CurrentPoint = -6200
	poolInfo.Liquidity = big.NewInt(200000)
	poolInfo.LiquidityX = big.NewInt(200000)
	var amount big.Int
	amount.SetString("274262809000", 10)
	swapResult, _ := SwapY2X(&amount, 1560, poolInfo)
	costY, _ := new(big.Int).SetString("274262809000", 10)
	acquireX, _ := new(big.Int).SetString("273333568073", 10)

	if swapResult.AmountY.Cmp(costY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapResult.AmountY.String(), costY.String())
	}
	if swapResult.AmountX.Cmp(acquireX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapResult.AmountX.String(), acquireX.String())
	}
	if swapResult.CurrentPoint != 1559 {
		t.Fatalf("result currentPoint not equal (%d, %d)", swapResult.CurrentPoint, 1559)
	}
	resultLiquidity, _ := new(big.Int).SetString("500000", 10)
	resultLiquidityX, _ := new(big.Int).SetString("51121", 10)
	if swapResult.Liquidity.Cmp(resultLiquidity) != 0 {
		t.Fatalf("Liquidity not equal (%s, %s)", swapResult.Liquidity.String(), resultLiquidity.String())
	}
	if swapResult.LiquidityX.Cmp(resultLiquidityX) != 0 {
		t.Fatalf("LiquidityX not equal (%s, %s)", swapResult.LiquidityX.String(), resultLiquidityX.String())
	}
}

func TestSwapDetailY2X7(t *testing.T) {
	// y2x start with limitorder-y, full-x-liquidity
	// end partial liquidity
	poolInfo := getPoolInfoDetailY2X()

	poolInfo.CurrentPoint = -6200
	poolInfo.Liquidity = big.NewInt(200000)
	poolInfo.LiquidityX = big.NewInt(198640)
	var amount big.Int
	amount.SetString("328166700000", 10)
	swapResult, _ := SwapY2X(&amount, 1560, poolInfo)
	costY, _ := new(big.Int).SetString("328166700000", 10)
	acquireX, _ := new(big.Int).SetString("373333543040", 10)

	if swapResult.AmountY.Cmp(costY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapResult.AmountY.String(), costY.String())
	}
	if swapResult.AmountX.Cmp(acquireX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapResult.AmountX.String(), acquireX.String())
	}
	if swapResult.CurrentPoint != 1559 {
		t.Fatalf("result currentPoint not equal (%d, %d)", swapResult.CurrentPoint, 1559)
	}
	resultLiquidity, _ := new(big.Int).SetString("500000", 10)
	resultLiquidityX, _ := new(big.Int).SetString("76178", 10)
	if swapResult.Liquidity.Cmp(resultLiquidity) != 0 {
		t.Fatalf("Liquidity not equal (%s, %s)", swapResult.Liquidity.String(), resultLiquidity.String())
	}
	if swapResult.LiquidityX.Cmp(resultLiquidityX) != 0 {
		t.Fatalf("LiquidityX not equal (%s, %s)", swapResult.LiquidityX.String(), resultLiquidityX.String())
	}
}

func TestSwapDetailY2X8(t *testing.T) {
	// y2x start with limitorder-y, full-x-liquidity
	// end partial liquidity
	poolInfo := getPoolInfoDetailY2X()

	poolInfo.CurrentPoint = -6200
	poolInfo.Liquidity = big.NewInt(200000)
	poolInfo.LiquidityX = big.NewInt(198640)
	var amount big.Int
	amount.SetString("250000000000", 10)
	swapResult, _ := SwapY2X(&amount, 1201, poolInfo)
	costY, _ := new(big.Int).SetString("249999999999", 10)
	acquireX, _ := new(big.Int).SetString("304147179726", 10)

	if swapResult.AmountY.Cmp(costY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapResult.AmountY.String(), costY.String())
	}
	if swapResult.AmountX.Cmp(acquireX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapResult.AmountX.String(), acquireX.String())
	}
	if swapResult.CurrentPoint != 1200 {
		t.Fatalf("result currentPoint not equal (%d, %d)", swapResult.CurrentPoint, 1200)
	}
	resultLiquidity, _ := new(big.Int).SetString("500000", 10)
	resultLiquidityX, _ := new(big.Int).SetString("500000", 10)
	if swapResult.Liquidity.Cmp(resultLiquidity) != 0 {
		t.Fatalf("Liquidity not equal (%s, %s)", swapResult.Liquidity.String(), resultLiquidity.String())
	}
	if swapResult.LiquidityX.Cmp(resultLiquidityX) != 0 {
		t.Fatalf("LiquidityX not equal (%s, %s)", swapResult.LiquidityX.String(), resultLiquidityX.String())
	}
}

func TestSwapDetailY2X9(t *testing.T) {
	// y2x start with limitorder-y, full-x-liquidity
	// end partial liquidity
	poolInfo := getPoolInfoDetailY2X()

	poolInfo.CurrentPoint = -6200
	poolInfo.Liquidity = big.NewInt(200000)
	poolInfo.LiquidityX = big.NewInt(198640)
	var amount big.Int
	amount.SetString("327974071999", 10)
	swapResult, _ := SwapY2X(&amount, 1201, poolInfo)
	costY, _ := new(big.Int).SetString("327974071999", 10)
	acquireX, _ := new(big.Int).SetString("373166078203", 10)

	if swapResult.AmountY.Cmp(costY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapResult.AmountY.String(), costY.String())
	}
	if swapResult.AmountX.Cmp(acquireX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapResult.AmountX.String(), acquireX.String())
	}
	if swapResult.CurrentPoint != 1200 {
		t.Fatalf("result currentPoint not equal (%d, %d)", swapResult.CurrentPoint, 1200)
	}
	resultLiquidity, _ := new(big.Int).SetString("500000", 10)
	resultLiquidityX, _ := new(big.Int).SetString("358", 10)
	if swapResult.Liquidity.Cmp(resultLiquidity) != 0 {
		t.Fatalf("Liquidity not equal (%s, %s)", swapResult.Liquidity.String(), resultLiquidity.String())
	}
	if swapResult.LiquidityX.Cmp(resultLiquidityX) != 0 {
		t.Fatalf("LiquidityX not equal (%s, %s)", swapResult.LiquidityX.String(), resultLiquidityX.String())
	}
}
