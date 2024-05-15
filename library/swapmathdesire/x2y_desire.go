package swapmathdesire

import (
	"math/big"

	"github.com/izumiFinance/iZiSwap-SDK-go/library/amountmath"
	"github.com/izumiFinance/iZiSwap-SDK-go/library/calc"
	"github.com/izumiFinance/iZiSwap-SDK-go/library/utils"
)

type X2YRangeRetState struct {
	// Whether user has acquire enough tokenY
	Finished bool
	// Actual cost of tokenX to buy tokenY
	CostX *big.Int
	// Amount of acquired tokenY
	AcquireY *big.Int
	// Final point after this swap
	FinalPt int
	// Sqrt price on final point
	SqrtFinalPrice_96 *big.Int
	// Liquidity of tokenX at finalPt
	LiquidityX *big.Int
}

func X2YAtPrice(desireY, sqrtPrice_96, currY *big.Int) (costX, acquireY *big.Int) {
	acquireY = new(big.Int).Set(desireY)
	if acquireY.Cmp(currY) == 1 {
		acquireY = new(big.Int).Set(currY)
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
	var costX, acquireY, newLiquidityX *big.Int
	var maxTransformLiquidityX, transformLiquidityX *big.Int

	liquidityY := new(big.Int).Sub(liquidity, liquidityX)
	// desireY * 2^96 <= 2^128 * 2^96 <= 2^224 < 2^256
	maxTransformLiquidityX = calc.MulDivCeil(desireY, utils.Pow96, sqrtPrice_96)
	// transformLiquidityX <= liquidityY <= uint128.max
	transformLiquidityX = calc.MinBigInt(maxTransformLiquidityX, liquidityY)
	// transformLiquidityX * 2^96 <= 2^128 * 2^96 <= 2^224 < 2^256
	costX = calc.MulDivCeil(transformLiquidityX, utils.Pow96, sqrtPrice_96)
	// acquireY should not > uint128.max
	acquireY = calc.MulDivFloor(transformLiquidityX, sqrtPrice_96, utils.Pow96)
	newLiquidityX = new(big.Int).Add(liquidityX, transformLiquidityX)
	return X2YAtPriceLiquidityResult{CostX: costX, AcquireY: acquireY, NewLiquidityX: newLiquidityX}
}

// struct Range {
// 	uint128 liquidity;
// 	uint160 sqrtPriceL_96;
// 	int24 leftPt;
// 	uint160 sqrtPriceR_96;
// 	int24 rightPt;
// 	uint160 sqrtRate_96;
// }

type RangeX2Y struct {
	Liquidity     *big.Int
	SqrtPriceL_96 *big.Int
	LeftPt        int
	SqrtPriceR_96 *big.Int
	RightPt       int
	SqrtRate_96   *big.Int
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

func x2YRangeComplete(rg RangeX2Y, desireY *big.Int) X2YRangeCompRet {
	var ret X2YRangeCompRet
	maxY := amountmath.GetAmountY(rg.Liquidity, rg.SqrtPriceL_96, rg.SqrtPriceR_96, rg.SqrtRate_96, false)
	if maxY.Cmp(desireY) <= 0 {
		ret.AcquireY = maxY
		ret.CostX = amountmath.GetAmountX(rg.Liquidity, rg.LeftPt, rg.RightPt, rg.SqrtPriceR_96, rg.SqrtRate_96, true)
		ret.CompleteLiquidity = true
		return ret
	}
	cl := new(big.Int).Sub(
		rg.SqrtPriceR_96,
		new(big.Int).Div(
			new(big.Int).Mul(
				desireY,
				new(big.Int).Sub(
					rg.SqrtRate_96,
					utils.Pow96,
				),
			),
			rg.Liquidity,
		),
	)

	logValue, _ := calc.GetLogSqrtPriceFloor(cl)
	ret.LocPt = logValue + 1

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
		ret.CostX = calc.MulDivCeil(
			rg.Liquidity,
			new(big.Int).Sub(sqrtPricePrMloc_96, utils.Pow96),
			new(big.Int).Sub(rg.SqrtPriceR_96, sqrtPricePrM1_96),
		)

		ret.LocPt = ret.LocPt - 1
		ret.SqrtLoc_96, _ = calc.GetSqrtPrice(ret.LocPt)

		sqrtLocA1_96 := new(big.Int).Add(
			ret.SqrtLoc_96,
			new(big.Int).Div(
				new(big.Int).Mul(
					ret.SqrtLoc_96,
					new(big.Int).Sub(
						rg.SqrtRate_96,
						utils.Pow96,
					),
				),
				utils.Pow96,
			),
		)

		acquireY256 := amountmath.GetAmountY(rg.Liquidity, sqrtLocA1_96, rg.SqrtPriceR_96, rg.SqrtRate_96, false)
		// ret.acquireY <= desireY <= uint128.max
		ret.AcquireY = calc.MinBigInt(acquireY256, desireY)
	}
	return ret
}

func X2YRange(
	currentState utils.State,
	leftPt int,
	sqrtRate_96 *big.Int,
	desireY *big.Int,
) X2YRangeRetState {
	desireY = new(big.Int).Set(desireY)
	var retState X2YRangeRetState
	retState.CostX = big.NewInt(0)
	retState.AcquireY = big.NewInt(0)
	retState.LiquidityX = big.NewInt(0)
	retState.Finished = false

	currentHasY := currentState.LiquidityX.Cmp(currentState.Liquidity) < 0
	if currentHasY && (currentState.LiquidityX.Cmp(big.NewInt(0)) > 0 || leftPt == currentState.CurrentPoint) {
		ret := x2YAtPriceLiquidity(
			desireY, currentState.SqrtPrice_96, currentState.Liquidity, currentState.LiquidityX,
		)
		retState.CostX = ret.CostX
		retState.AcquireY = ret.AcquireY
		retState.LiquidityX = ret.NewLiquidityX
		if retState.LiquidityX.Cmp(currentState.Liquidity) < 0 || retState.AcquireY.Cmp(desireY) >= 0 {
			// remaining desire y is not enough to down current price to price / 1.0001
			// but desire y may remain, so we cannot simply use (retState.acquireY >= desireY)
			retState.Finished = true
			retState.FinalPt = currentState.CurrentPoint
			retState.SqrtFinalPrice_96 = new(big.Int).Set(currentState.SqrtPrice_96)
		} else {
			desireY.Sub(desireY, retState.AcquireY)
		}
	} else if currentHasY { // all y
		currentState.CurrentPoint = currentState.CurrentPoint + 1
		// sqrt(price) + sqrt(price) * (1.0001 - 1) == sqrt(price) * 1.0001
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
		retState.LiquidityX = new(big.Int).Set(currentState.LiquidityX)
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
			retState.Finished = (desireY.Cmp(big.NewInt(0)) <= 0)
			retState.FinalPt = leftPt
			retState.SqrtFinalPrice_96 = sqrtPriceL_96
			retState.LiquidityX = new(big.Int).Set(currentState.Liquidity)
		} else {
			// locPt > leftPt
			// trade at locPt
			locRet := x2YAtPriceLiquidity(
				desireY, ret.SqrtLoc_96, currentState.Liquidity, big.NewInt(0),
			)
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
		// finishd must be false
		// retState.finished == false;
		retState.FinalPt = currentState.CurrentPoint
		retState.SqrtFinalPrice_96 = new(big.Int).Set(currentState.SqrtPrice_96)
	}

	return retState
}
