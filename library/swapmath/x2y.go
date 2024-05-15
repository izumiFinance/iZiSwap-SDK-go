package swapmath

import (
	"math/big"

	"github.com/izumiFinance/iZiSwap-SDK-go/library/amountmath"
	"github.com/izumiFinance/iZiSwap-SDK-go/library/calc"
	"github.com/izumiFinance/iZiSwap-SDK-go/library/utils"
)

type X2YRangeRetState struct {
	// whether user run out of amountX
	Finished bool
	// actual cost of tokenX to buy tokenY
	CostX *big.Int
	// amount of acquired tokenY
	AcquireY *big.Int
	// final point after this swap
	FinalPt int
	// sqrt price on final point
	SqrtFinalPrice_96 *big.Int
	// liquidity of tokenX at finalPt
	LiquidityX *big.Int
}

func X2YAtPrice(amountX, sqrtPrice_96, currY *big.Int) (costX, acquireY *big.Int) {

	l := calc.MulDivFloor(amountX, sqrtPrice_96, utils.Pow96)
	acquireY = calc.MulDivFloor(l, sqrtPrice_96, utils.Pow96)
	if acquireY.Cmp(currY) > 0 {
		acquireY.Set(currY)
	}
	l = calc.MulDivCeil(acquireY, utils.Pow96, sqrtPrice_96)
	costX = calc.MulDivCeil(l, utils.Pow96, sqrtPrice_96)

	return costX, acquireY
}

type X2YAtPriceLiquidityResult struct {
	CostX         *big.Int
	AcquireY      *big.Int
	NewLiquidityX *big.Int
}

func x2YAtPriceLiquidity(amountX, sqrtPrice_96, liquidity, liquidityX *big.Int) X2YAtPriceLiquidityResult {
	var costX, acquireY, newLiquidityX *big.Int
	var maxTransformLiquidityX, transformLiquidityX *big.Int

	liquidityY := new(big.Int).Sub(liquidity, liquidityX)
	maxTransformLiquidityX = calc.MulDivFloor(amountX, sqrtPrice_96, utils.Pow96)
	transformLiquidityX = calc.MinBigInt(maxTransformLiquidityX, liquidityY)

	costX = calc.MulDivCeil(transformLiquidityX, utils.Pow96, sqrtPrice_96)
	acquireY = calc.MulDivFloor(transformLiquidityX, sqrtPrice_96, utils.Pow96)
	newLiquidityX = new(big.Int).Add(liquidityX, transformLiquidityX)

	return X2YAtPriceLiquidityResult{CostX: costX, AcquireY: acquireY, NewLiquidityX: newLiquidityX}
}

type X2YRangeCompRet struct {
	// cost of tokenX to buy tokenY
	CostX *big.Int
	// amount of acquired tokenY
	AcquireY *big.Int
	// whether all liquidity is used
	CompleteLiquidity bool
	// location point after this swap
	LocPt int
	// sqrt location after this swap
	SqrtLoc_96 *big.Int
}

type RangeX2Y struct {
	// total liquidity in this range
	Liquidity *big.Int
	// sqrt price on left point
	SqrtPriceL_96 *big.Int
	// left point of this range
	LeftPt int
	// sqrt price on right point
	SqrtPriceR_96 *big.Int
	// right point of this range
	RightPt int
	// sqrt rate of this range
	SqrtRate_96 *big.Int
}

func x2YRangeComplete(rg RangeX2Y, amountX *big.Int) X2YRangeCompRet {
	var ret X2YRangeCompRet
	sqrtPricePrM1_96 := calc.MulDivCeil(rg.SqrtPriceR_96, utils.Pow96, rg.SqrtRate_96)
	sqrtPricePrMl_96, _ := calc.GetSqrtPrice(rg.RightPt - rg.LeftPt)
	maxX := calc.MulDivCeil(rg.Liquidity, new(big.Int).Sub(sqrtPricePrMl_96, utils.Pow96), new(big.Int).Sub(rg.SqrtPriceR_96, sqrtPricePrM1_96))

	if maxX.Cmp(amountX) <= 0 {
		ret.CostX = maxX
		ret.AcquireY = amountmath.GetAmountY(rg.Liquidity, rg.SqrtPriceL_96, rg.SqrtPriceR_96, rg.SqrtRate_96, false)
		ret.CompleteLiquidity = true
	} else {
		sqrtValue_96 := new(big.Int).Add(
			new(big.Int).Div(
				new(big.Int).Mul(
					amountX,
					new(big.Int).Sub(rg.SqrtPriceR_96, sqrtPricePrM1_96),
				),
				rg.Liquidity,
			),
			utils.Pow96,
		)

		logValue, _ := calc.GetLogSqrtPriceFloor(sqrtValue_96)

		ret.LocPt = rg.RightPt - logValue
		ret.LocPt = calc.Min(ret.LocPt, rg.RightPt)
		ret.LocPt = calc.Max(ret.LocPt, rg.LeftPt+1)
		ret.CompleteLiquidity = false

		if ret.LocPt == rg.RightPt {
			ret.CostX = big.NewInt(0)
			ret.AcquireY = big.NewInt(0)
			ret.LocPt = ret.LocPt - 1
			ret.SqrtLoc_96, _ = calc.GetSqrtPrice(ret.LocPt)
		} else {
			sqrtPricePrMloc_96, _ := calc.GetSqrtPrice(rg.RightPt - ret.LocPt)
			costX256 := calc.MulDivCeil(rg.Liquidity, new(big.Int).Sub(sqrtPricePrMloc_96, utils.Pow96), new(big.Int).Sub(rg.SqrtPriceR_96, sqrtPricePrM1_96))
			ret.CostX = calc.MinBigInt(costX256, amountX)
			ret.LocPt = ret.LocPt - 1
			ret.SqrtLoc_96, _ = calc.GetSqrtPrice(ret.LocPt)
			sqrtLocA1_96 := new(big.Int).Add(
				ret.SqrtLoc_96,
				new(big.Int).Div(
					new(big.Int).Mul(
						ret.SqrtLoc_96,
						new(big.Int).Sub(rg.SqrtRate_96, utils.Pow96),
					),
					utils.Pow96,
				),
			)
			ret.AcquireY = amountmath.GetAmountY(rg.Liquidity, sqrtLocA1_96, rg.SqrtPriceR_96, rg.SqrtRate_96, false)
		}
	}

	return ret
}

func X2YRange(currentState utils.State, leftPt int, sqrtRate_96 *big.Int, amountX *big.Int) X2YRangeRetState {
	var retState X2YRangeRetState
	retState.CostX = big.NewInt(0)
	retState.AcquireY = big.NewInt(0)
	retState.LiquidityX = big.NewInt(0)
	retState.Finished = false

	currentHasY := currentState.LiquidityX.Cmp(currentState.Liquidity) < 0
	if currentHasY && (currentState.LiquidityX.Cmp(new(big.Int).SetInt64(0)) > 0 || leftPt == currentState.CurrentPoint) {
		ret := x2YAtPriceLiquidity(amountX, currentState.SqrtPrice_96, currentState.Liquidity, currentState.LiquidityX)
		retState.CostX = ret.CostX
		retState.AcquireY = ret.AcquireY
		retState.LiquidityX = ret.NewLiquidityX
		if retState.LiquidityX.Cmp(currentState.Liquidity) < 0 || retState.CostX.Cmp(amountX) >= 0 {
			retState.Finished = true
			retState.FinalPt = currentState.CurrentPoint
			retState.SqrtFinalPrice_96 = currentState.SqrtPrice_96
		} else {
			amountX.Sub(amountX, retState.CostX)
		}
	} else if currentHasY { // all y
		currentState.CurrentPoint = currentState.CurrentPoint + 1
		currentState.SqrtPrice_96 = new(big.Int).Add(
			currentState.SqrtPrice_96,
			new(big.Int).Div(
				new(big.Int).Mul(
					currentState.SqrtPrice_96,
					new(big.Int).Sub(sqrtRate_96, utils.Pow96),
				),
				utils.Pow96,
			),
		)
	} else {
		retState.LiquidityX = currentState.LiquidityX
	}

	if retState.Finished {
		return retState
	}

	if leftPt < currentState.CurrentPoint {
		sqrtPriceL_96, _ := calc.GetSqrtPrice(leftPt)
		ret := x2YRangeComplete(
			RangeX2Y{
				Liquidity:     currentState.Liquidity,
				SqrtPriceL_96: sqrtPriceL_96,
				LeftPt:        leftPt,
				SqrtPriceR_96: currentState.SqrtPrice_96,
				RightPt:       currentState.CurrentPoint,
				SqrtRate_96:   sqrtRate_96,
			},
			amountX,
		)
		retState.CostX.Add(retState.CostX, ret.CostX)
		amountX.Sub(amountX, ret.CostX)
		retState.AcquireY.Add(retState.AcquireY, ret.AcquireY)
		if ret.CompleteLiquidity {
			retState.Finished = (amountX.Cmp(new(big.Int).SetInt64(0)) <= 0)
			retState.FinalPt = leftPt
			retState.SqrtFinalPrice_96 = sqrtPriceL_96
			retState.LiquidityX = currentState.Liquidity
		} else {
			locRet := x2YAtPriceLiquidity(amountX, ret.SqrtLoc_96, currentState.Liquidity, big.NewInt(0))
			locCostX := locRet.CostX
			locAcquireY := locRet.AcquireY
			retState.LiquidityX = locRet.NewLiquidityX
			retState.CostX.Add(retState.CostX, locCostX)
			retState.AcquireY.Add(retState.AcquireY, locAcquireY)
			retState.Finished = true
			retState.SqrtFinalPrice_96 = ret.SqrtLoc_96
			retState.FinalPt = ret.LocPt
		}
	} else {
		retState.FinalPt = currentState.CurrentPoint
		retState.SqrtFinalPrice_96 = currentState.SqrtPrice_96
	}

	return retState
}
