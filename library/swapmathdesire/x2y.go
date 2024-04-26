package swapmathdesire

import (
	"math/big"

	"github.com/izumiFinance/iZiSwap-SDK-go/library/amountmath"
	"github.com/izumiFinance/iZiSwap-SDK-go/library/calc"
	"github.com/izumiFinance/iZiSwap-SDK-go/library/utils"
)

var zeroBI = big.NewInt(0)

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

func X2YAtPrice(desireY, sqrtPrice_96, currY *big.Int) (costX, acquireY *big.Int) {
	acquireY = new(big.Int).Set(desireY)
	if acquireY.Cmp(currY) > 0 {
		acquireY.Set(currY)
	}

	l := calc.MulDivCeil(acquireY, utils.Pow96, sqrtPrice_96)
	costX = calc.MulDivCeil(l, utils.Pow96, sqrtPrice_96)

	return costX, acquireY
}

type X2YAtPriceLiquidityResult struct {
	CostX         *big.Int
	AcquireY      *big.Int
	NewLiquidityX *big.Int
}

func x2YAtPriceLiquidity(desireY, sqrtPrice_96, liquidity, liquidityX *big.Int) X2YAtPriceLiquidityResult {
	liquidityY := new(big.Int).Sub(liquidity, liquidityX)
	// desireY * 2^96 <= 2^128 * 2^96 <= 2^224 < 2^256
	maxTransformLiquidityX := calc.MulDivCeil(desireY, utils.Pow96, sqrtPrice_96)
	// transformLiquidityX <= liquidityY <= uint128.max
	transformLiquidityX := calc.MinBigInt(maxTransformLiquidityX, liquidityY)
	// transformLiquidityX * 2^96 <= 2^128 * 2^96 <= 2^224 < 2^256
	costX := calc.MulDivCeil(transformLiquidityX, utils.Pow96, sqrtPrice_96)
	// acquireY should not > uint128.max
	acquireY := calc.MulDivFloor(transformLiquidityX, sqrtPrice_96, utils.Pow96)
	newLiquidityX := new(big.Int).Add(liquidityX, transformLiquidityX)

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

func x2YRangeComplete(rg RangeX2Y, desireY *big.Int) X2YRangeCompRet {
	var ret X2YRangeCompRet

	maxY := amountmath.GetAmountY(rg.Liquidity, rg.SqrtPriceL_96, rg.SqrtPriceR_96, rg.SqrtRate_96, false)
	if maxY.Cmp(desireY) <= 0 {
		ret.AcquireY = maxY
		ret.CostX = amountmath.GetAmountX(rg.Liquidity, rg.LeftPt, rg.RightPt, rg.SqrtPriceR_96, rg.SqrtRate_96, true)
		ret.CompleteLiquidity = true
		return ret
	}

	// 1. desireY * (rg.sqrtRate_96 - 2^96)
	//    < 2^128 * 2^96
	//    = 2 ^ 224 < 2 ^ 256
	// 2. desireY < maxY = rg.liquidity * (rg.sqrtPriceR_96 - rg.sqrtPriceL_96) / (rg.sqrtRate_96 - 2^96)
	// here, '/' means div of int
	// desireY < rg.liquidity * (rg.sqrtPriceR_96 - rg.sqrtPriceL_96) / (rg.sqrtRate_96 - 2^96)
	// => desireY * (rg.sqrtRate_96 - TwoPower.Pow96) / rg.liquidity < rg.sqrtPriceR_96 - rg.sqrtPriceL_96
	// => rg.sqrtPriceR_96 - desireY * (rg.sqrtRate_96 - TwoPower.Pow96) / rg.liquidity > rg.sqrtPriceL_96
	cl := new(big.Int).Sub(
		rg.SqrtPriceR_96,
		new(big.Int).Div(
			new(big.Int).Mul(
				desireY,
				new(big.Int).Sub(rg.SqrtRate_96, utils.Pow96)),
			rg.Liquidity))
	ret.LocPt, _ = calc.GetLogSqrtPriceFloor(cl)
	ret.LocPt = ret.LocPt + 1
	ret.LocPt = calc.Min(ret.LocPt, rg.RightPt)
	ret.LocPt = calc.Max(ret.LocPt, rg.LeftPt+1)
	ret.CompleteLiquidity = false

	if ret.LocPt == rg.RightPt {
		ret.CostX = big.NewInt(0)
		ret.AcquireY = big.NewInt(0)
		ret.LocPt = ret.LocPt - 1
		ret.SqrtLoc_96, _ = calc.GetSqrtPrice(ret.LocPt)
	} else {
		// rg.rightPt - ret.locPt <= 256 * 100
		// sqrtPricePrMloc_96 <= 1.0001 ** 25600 * 2 ^ 96 = 13 * 2^96 < 2^100
		sqrtPricePrMloc_96, _ := calc.GetSqrtPrice(rg.RightPt - ret.LocPt)
		// rg.sqrtPriceR_96 * TwoPower.Pow96 < 2^160 * 2^96 = 2^256
		sqrtPricePrM1_96 := calc.MulDivCeil(rg.SqrtPriceR_96, utils.Pow96, rg.SqrtRate_96)
		// rg.liquidity * (sqrtPricePrMloc_96 - TwoPower.Pow96) < 2^128 * 2^100 = 2^228 < 2^256
		ret.CostX = calc.MulDivCeil(rg.Liquidity, new(big.Int).Sub(sqrtPricePrMloc_96, utils.Pow96), new(big.Int).Sub(rg.SqrtPriceR_96, sqrtPricePrM1_96))

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
		acquireT256 := amountmath.GetAmountY(rg.Liquidity, sqrtLocA1_96, rg.SqrtPriceR_96, rg.SqrtRate_96, false)
		// ret.acquireY <= desireY <= uint128.max
		ret.AcquireY = calc.MinBigInt(acquireT256, desireY)
	}

	return ret
}

func X2YRange(currentState utils.State, leftPt int, sqrtRate_96 *big.Int, desireY *big.Int) X2YRangeRetState {
	var retState X2YRangeRetState

	retState.CostX = big.NewInt(0)
	retState.AcquireY = big.NewInt(0)
	retState.Finished = false

	currentHasY := currentState.LiquidityX.Cmp(currentState.Liquidity) < 0
	if currentHasY && (currentState.LiquidityX.Cmp(zeroBI) > 0 || leftPt == currentState.CurrentPoint) {
		ret := x2YAtPriceLiquidity(desireY, currentState.SqrtPrice_96, currentState.Liquidity, currentState.LiquidityX)
		retState.CostX = ret.CostX
		retState.AcquireY = ret.AcquireY
		retState.LiquidityX = ret.NewLiquidityX
		if retState.LiquidityX.Cmp(currentState.Liquidity) < 0 || retState.AcquireY.Cmp(desireY) >= 0 {
			retState.Finished = true
			retState.FinalPt = currentState.CurrentPoint
			retState.SqrtFinalPrice_96 = currentState.SqrtPrice_96
		} else {
			desireY.Sub(desireY, retState.AcquireY)
		}
	} else if currentHasY {
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
			desireY,
		)
		retState.CostX.Add(retState.CostX, ret.CostX)
		desireY.Sub(desireY, ret.AcquireY)
		retState.AcquireY.Add(retState.AcquireY, ret.AcquireY)
		if ret.CompleteLiquidity {
			retState.Finished = (desireY.Cmp(zeroBI) == 0)
			retState.FinalPt = leftPt
			retState.SqrtFinalPrice_96 = sqrtPriceL_96
			retState.LiquidityX = currentState.Liquidity
		} else {
			locRet := x2YAtPriceLiquidity(desireY, ret.SqrtLoc_96, currentState.Liquidity, zeroBI)
			locCostX := locRet.CostX
			locAcquireY := locRet.AcquireY
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
