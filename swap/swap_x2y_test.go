package swap

import (
	"math/big"
	"testing"
)

func getLiquiditiesX2Y() []LiquidityPoint {
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

func getLimitOrdersY() []LimitOrderPoint {
	limitOrders := []LimitOrderPoint{
		{SellingY: big.NewInt(100000000000), Point: -3000},
		{SellingY: big.NewInt(150000000000), Point: -1000},
		{SellingY: big.NewInt(120000000000), Point: 1200},
	}
	return limitOrders
}

func getPoolInfoX2Y() PoolInfo {
	return PoolInfo{
		CurrentPoint: 1887,
		PointDelta:   40,
		LeftMostPt:   -800000,
		RightMostPt:  800000,
		Fee:          2000,
		Liquidity:    big.NewInt(700000),
		LiquidityX:   big.NewInt(246660),
		Liquidities:  getLiquiditiesX2Y(),
		LimitOrders:  getLimitOrdersY(),
	}
}

func TestSwapX2Y1(t *testing.T) {
	poolInfo := getPoolInfoX2Y()
	var amount big.Int
	amount.SetString("100000000000000000000000", 10)
	swapAmount, _ := SwapX2Y(&amount, -6123, poolInfo)
	costX, _ := new(big.Int).SetString("410079196782", 10)
	acquireY, _ := new(big.Int).SetString("371715048235", 10)
	if swapAmount.AmountX.Cmp(costX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapAmount.AmountX.String(), costX.String())
	}
	if swapAmount.AmountY.Cmp(acquireY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapAmount.AmountY.String(), acquireY.String())
	}
}
func TestSwapX2Y2(t *testing.T) {
	poolInfo := getPoolInfoX2Y()
	var amount big.Int
	amount.SetString("410079196782", 10)
	swapAmount, _ := SwapX2Y(&amount, -6123, poolInfo)
	costX, _ := new(big.Int).SetString("410079196782", 10)
	acquireY, _ := new(big.Int).SetString("371715048235", 10)
	if swapAmount.AmountX.Cmp(costX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapAmount.AmountX.String(), costX.String())
	}
	if swapAmount.AmountY.Cmp(acquireY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapAmount.AmountY.String(), acquireY.String())
	}
}
func TestSwapX2Y3(t *testing.T) {
	poolInfo := getPoolInfoX2Y()
	var amount big.Int
	amount.SetString("399624951498", 10)
	swapAmount, _ := SwapX2Y(&amount, -6123, poolInfo)
	costX, _ := new(big.Int).SetString("399624951497", 10)
	acquireY, _ := new(big.Int).SetString("364135750158", 10)
	if swapAmount.AmountX.Cmp(costX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapAmount.AmountX.String(), costX.String())
	}
	if swapAmount.AmountY.Cmp(acquireY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapAmount.AmountY.String(), acquireY.String())
	}
}

func TestSwapX2Y4(t *testing.T) {
	poolInfo := getPoolInfoX2Y()
	var amount big.Int
	amount.SetString("368662456348", 10)
	swapAmount, _ := SwapX2Y(&amount, -6123, poolInfo)
	costX, _ := new(big.Int).SetString("368662456348", 10)
	acquireY, _ := new(big.Int).SetString("341243701400", 10)
	if swapAmount.AmountX.Cmp(costX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapAmount.AmountX.String(), costX.String())
	}
	if swapAmount.AmountY.Cmp(acquireY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapAmount.AmountY.String(), acquireY.String())
	}
}

func TestSwapX2Y5(t *testing.T) {
	poolInfo := getPoolInfoX2Y()
	var amount big.Int
	amount.SetString("245774970898", 10)
	swapAmount, _ := SwapX2Y(&amount, -6123, poolInfo)
	costX, _ := new(big.Int).SetString("245774970898", 10)
	acquireY, _ := new(big.Int).SetString("245800546380", 10)
	if swapAmount.AmountX.Cmp(costX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapAmount.AmountX.String(), costX.String())
	}
	if swapAmount.AmountY.Cmp(acquireY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapAmount.AmountY.String(), acquireY.String())
	}
}

func TestSwapX2Y6(t *testing.T) {
	poolInfo := getPoolInfoX2Y()
	var amount big.Int
	amount.SetString("122887485449", 10)
	swapAmount, _ := SwapX2Y(&amount, -6123, poolInfo)
	costX, _ := new(big.Int).SetString("122887485449", 10)
	acquireY, _ := new(big.Int).SetString("134829182908", 10)
	if swapAmount.AmountX.Cmp(costX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapAmount.AmountX.String(), costX.String())
	}
	if swapAmount.AmountY.Cmp(acquireY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapAmount.AmountY.String(), acquireY.String())
	}
}

func TestSwapX2YDesire1(t *testing.T) {
	poolInfo := getPoolInfoX2Y()
	var amount big.Int
	amount.SetString("100000000000000000000000", 10)
	swapAmount, _ := SwapX2YDesireY(&amount, -6123, poolInfo)
	costX, _ := new(big.Int).SetString("410079196782", 10)
	acquireY, _ := new(big.Int).SetString("371715048235", 10)
	if swapAmount.AmountX.Cmp(costX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapAmount.AmountX.String(), costX.String())
	}
	if swapAmount.AmountY.Cmp(acquireY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapAmount.AmountY.String(), acquireY.String())
	}
}
func TestSwapX2YDesire2(t *testing.T) {
	poolInfo := getPoolInfoX2Y()
	var amount big.Int
	amount.SetString("364135750158", 10)
	swapAmount, _ := SwapX2YDesireY(&amount, -6123, poolInfo)
	costX, _ := new(big.Int).SetString("399624951497", 10)
	acquireY, _ := new(big.Int).SetString("364135750158", 10)
	if swapAmount.AmountX.Cmp(costX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapAmount.AmountX.String(), costX.String())
	}
	if swapAmount.AmountY.Cmp(acquireY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapAmount.AmountY.String(), acquireY.String())
	}
}
func TestSwapX2YDesire3(t *testing.T) {
	poolInfo := getPoolInfoX2Y()
	var amount big.Int
	amount.SetString("341243701400", 10)
	swapAmount, _ := SwapX2YDesireY(&amount, -6123, poolInfo)
	costX, _ := new(big.Int).SetString("368662456348", 10)
	acquireY, _ := new(big.Int).SetString("341243701400", 10)
	if swapAmount.AmountX.Cmp(costX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapAmount.AmountX.String(), costX.String())
	}
	if swapAmount.AmountY.Cmp(acquireY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapAmount.AmountY.String(), acquireY.String())
	}
}

func TestSwapX2YDesire4(t *testing.T) {
	poolInfo := getPoolInfoX2Y()
	var amount big.Int
	amount.SetString("245800546380", 10)
	swapAmount, _ := SwapX2YDesireY(&amount, -6123, poolInfo)
	costX, _ := new(big.Int).SetString("245774970898", 10)
	acquireY, _ := new(big.Int).SetString("245800546380", 10)
	if swapAmount.AmountX.Cmp(costX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapAmount.AmountX.String(), costX.String())
	}
	if swapAmount.AmountY.Cmp(acquireY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapAmount.AmountY.String(), acquireY.String())
	}
}

func TestSwapX2YDesire5(t *testing.T) {
	poolInfo := getPoolInfoX2Y()
	var amount big.Int
	amount.SetString("134829182908", 10)
	swapAmount, _ := SwapX2YDesireY(&amount, -6123, poolInfo)
	costX, _ := new(big.Int).SetString("122887485449", 10)
	acquireY, _ := new(big.Int).SetString("134829182908", 10)
	if swapAmount.AmountX.Cmp(costX) != 0 {
		t.Fatalf("amount x not equal (%s, %s)", swapAmount.AmountX.String(), costX.String())
	}
	if swapAmount.AmountY.Cmp(acquireY) != 0 {
		t.Fatalf("amount y not equal (%s, %s)", swapAmount.AmountY.String(), acquireY.String())
	}
}
