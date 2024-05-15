package swap

import (
	"math/big"
	"testing"
)

func getLiquiditiesY2X() []LiquidityPoint {
	liquidities := []LiquidityPoint{
		{LiqudityDelta: big.NewInt(200000), Point: -9000},
		{LiqudityDelta: big.NewInt(300000), Point: -8000},
		{LiqudityDelta: big.NewInt(-300000), Point: -5000},
		{LiqudityDelta: big.NewInt(-200000), Point: -4000},

		{LiqudityDelta: big.NewInt(100000), Point: -2000},
		{LiqudityDelta: big.NewInt(500000), Point: -1200},
		{LiqudityDelta: big.NewInt(-500000), Point: -800},
		{LiqudityDelta: big.NewInt(-100000), Point: 800},

		{LiqudityDelta: big.NewInt(700000), Point: 1000},
		{LiqudityDelta: big.NewInt(-700000), Point: 2000},
	}
	return liquidities
}

func getLimitOrdersX() []LimitOrderPoint {
	limitOrders := []LimitOrderPoint{
		{SellingX: big.NewInt(100000000000), Point: -3000},
		{SellingX: big.NewInt(150000000000), Point: -1000},
		{SellingX: big.NewInt(120000000000), Point: 1200},
	}
	return limitOrders
}

func getPoolInfoY2X() PoolInfo {
	return PoolInfo{
		CurrentPoint: -6182,
		PointDelta:   40,
		LeftMostPt:   -800000,
		RightMostPt:  800000,
		Fee:          2000,
		Liquidity:    big.NewInt(500000),
		LiquidityX:   big.NewInt(202614),
		Liquidities:  getLiquiditiesY2X(),
		LimitOrders:  getLimitOrdersX(),
	}
}

func TestSwapY2X1(t *testing.T) {
	poolInfo := getPoolInfoY2X()
	var amount big.Int
	amount.SetString("100000000000000000000000", 10)
	swapAmount, _ := SwapY2X(&amount, 1100, poolInfo)
	costY, _ := new(big.Int).SetString("211374358247", 10)
	acquireX, _ := new(big.Int).SetString("251597283132", 10)
	if swapAmount.AmountY.Cmp(costY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapAmount.AmountY.String(), "211374358247")
	}
	if swapAmount.AmountX.Cmp(acquireX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapAmount.AmountX.String(), "251597283132")
	}
}

func TestSwapY2X2(t *testing.T) {
	poolInfo := getPoolInfoY2X()
	var amount big.Int
	amount.SetString("211374358247", 10)
	swapAmount, _ := SwapY2X(&amount, 1100, poolInfo)
	costY, _ := new(big.Int).SetString("211374358247", 10)
	acquireX, _ := new(big.Int).SetString("251597283132", 10)
	if swapAmount.AmountY.Cmp(costY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapAmount.AmountY.String(), "211374358247")
	}
	if swapAmount.AmountX.Cmp(acquireX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapAmount.AmountX.String(), "251597283132")
	}
}

func TestSwapY2X3(t *testing.T) {
	poolInfo := getPoolInfoY2X()
	var amount big.Int
	amount.SetString("190236922422", 10)
	swapAmount, _ := SwapY2X(&amount, 1100, poolInfo)
	costY, _ := new(big.Int).SetString("190236922422", 10)
	acquireX, _ := new(big.Int).SetString("228316826682", 10)
	if swapAmount.AmountY.Cmp(costY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapAmount.AmountY.String(), "211374358247")
	}
	if swapAmount.AmountX.Cmp(acquireX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapAmount.AmountX.String(), "251597283132")
	}
}

func TestSwapY2X4(t *testing.T) {
	poolInfo := getPoolInfoY2X()
	var amount big.Int
	amount.SetString("126824614948", 10)
	swapAmount, _ := SwapY2X(&amount, 1100, poolInfo)
	costY, _ := new(big.Int).SetString("126824614948", 10)
	acquireX, _ := new(big.Int).SetString("158375901172", 10)
	if swapAmount.AmountY.Cmp(costY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapAmount.AmountY.String(), "211374358247")
	}
	if swapAmount.AmountX.Cmp(acquireX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapAmount.AmountX.String(), "251597283132")
	}
}

func TestSwapY2X5(t *testing.T) {
	poolInfo := getPoolInfoY2X()
	var amount big.Int
	amount.SetString("63412307474", 10)
	swapAmount, _ := SwapY2X(&amount, 1100, poolInfo)
	costY, _ := new(big.Int).SetString("63412307474", 10)
	acquireX, _ := new(big.Int).SetString("85638433523", 10)
	if swapAmount.AmountY.Cmp(costY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapAmount.AmountY.String(), "211374358247")
	}
	if swapAmount.AmountX.Cmp(acquireX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapAmount.AmountX.String(), "251597283132")
	}
}

func TestSwapY2XDesire1(t *testing.T) {
	poolInfo := getPoolInfoY2X()
	var amount big.Int
	amount.SetString("100000000000000000000000", 10)
	swapAmount, _ := SwapY2XDesireX(&amount, 1100, poolInfo)
	costY, _ := new(big.Int).SetString("211374358247", 10)
	acquireX, _ := new(big.Int).SetString("251597283132", 10)
	if swapAmount.AmountY.Cmp(costY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapAmount.AmountY.String(), "211374358247")
	}
	if swapAmount.AmountX.Cmp(acquireX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapAmount.AmountX.String(), "251597283132")
	}
}

func TestSwapY2XDesire2(t *testing.T) {
	poolInfo := getPoolInfoY2X()
	var amount big.Int
	amount.SetString("251597283132", 10)
	swapAmount, _ := SwapY2XDesireX(&amount, 1100, poolInfo)
	costY, _ := new(big.Int).SetString("211374358247", 10)
	acquireX, _ := new(big.Int).SetString("251597283132", 10)
	if swapAmount.AmountY.Cmp(costY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapAmount.AmountY.String(), "211374358247")
	}
	if swapAmount.AmountX.Cmp(acquireX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapAmount.AmountX.String(), "251597283132")
	}
}

func TestSwapY2XDesire3(t *testing.T) {
	poolInfo := getPoolInfoY2X()
	var amount big.Int
	amount.SetString("228316826682", 10)
	swapAmount, _ := SwapY2XDesireX(&amount, 1100, poolInfo)
	costY, _ := new(big.Int).SetString("190236922422", 10)
	acquireX, _ := new(big.Int).SetString("228316826682", 10)
	if swapAmount.AmountY.Cmp(costY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapAmount.AmountY.String(), "211374358247")
	}
	if swapAmount.AmountX.Cmp(acquireX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapAmount.AmountX.String(), "251597283132")
	}
}

func TestSwapY2XDesire4(t *testing.T) {
	poolInfo := getPoolInfoY2X()
	var amount big.Int
	amount.SetString("158375901172", 10)
	swapAmount, _ := SwapY2XDesireX(&amount, 1100, poolInfo)
	costY, _ := new(big.Int).SetString("126824614948", 10)
	acquireX, _ := new(big.Int).SetString("158375901172", 10)
	if swapAmount.AmountY.Cmp(costY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapAmount.AmountY.String(), "211374358247")
	}
	if swapAmount.AmountX.Cmp(acquireX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapAmount.AmountX.String(), "251597283132")
	}
}

func TestSwapY2XDesire5(t *testing.T) {
	poolInfo := getPoolInfoY2X()
	var amount big.Int
	amount.SetString("85638433523", 10)
	swapAmount, _ := SwapY2XDesireX(&amount, 1100, poolInfo)
	costY, _ := new(big.Int).SetString("63412307474", 10)
	acquireX, _ := new(big.Int).SetString("85638433523", 10)
	if swapAmount.AmountY.Cmp(costY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapAmount.AmountY.String(), "211374358247")
	}
	if swapAmount.AmountX.Cmp(acquireX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapAmount.AmountX.String(), "251597283132")
	}
}
