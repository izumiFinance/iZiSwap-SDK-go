package swapmath

import (
	"math/big"

	"github.com/izumiFinance/swap/library/amountmath"
	"github.com/izumiFinance/swap/library/calc"
	"github.com/izumiFinance/swap/library/utils"
)

type Y2XRangeRetState struct {
	// whether user has run out of tokenY
	Finished bool
	// actual cost of tokenY to buy tokenX
	CostY *big.Int
	// actual amount of tokenX acquired
	AcquireX *big.Int
	// final point after this swap
	FinalPt int
	// sqrt price on final point
	SqrtFinalPrice_96 *big.Int
	// liquidity of tokenX at finalPt
	// if finalPt is not rightPt, liquidityX is meaningless
	LiquidityX *big.Int
}

func Y2XAtPrice(amountY *big.Int, sqrtPrice_96 *big.Int, currX *big.Int) (costY, acquireX *big.Int) {
	l := calc.MulDivFloor(amountY, utils.Pow96, sqrtPrice_96)
	// acquireX <= currX <= uint128.max
	acquireX = calc.MinBigInt(calc.MulDivFloor(l, utils.Pow96, sqrtPrice_96), currX)
	l = calc.MulDivCeil(acquireX, sqrtPrice_96, utils.Pow96)
	costY = calc.MulDivCeil(l, sqrtPrice_96, utils.Pow96)
	return costY, acquireX
}

type Y2XAtPriceLiquidityResult struct {
	CostY         *big.Int
	AcquireX      *big.Int
	NewLiquidityX *big.Int
}

func y2XAtPriceLiquidity(amountY *big.Int, sqrtPrice_96 *big.Int, liquidityX *big.Int) Y2XAtPriceLiquidityResult {
	// amountY * TwoPower.Pow96 < 2^128 * 2^96 = 2^224 < 2^256
	maxTransformLiquidityY := new(big.Int).Mul(amountY, utils.Pow96)
	maxTransformLiquidityY.Div(maxTransformLiquidityY, sqrtPrice_96)
	// transformLiquidityY <= liquidityX
	transformLiquidityY := calc.MinBigInt(maxTransformLiquidityY, liquidityX)
	// costY <= amountY
	costY := calc.MulDivCeil(transformLiquidityY, sqrtPrice_96, utils.Pow96)
	// transformLiquidityY * 2^96 < 2^224 < 2^256
	acquireX := new(big.Int).Mul(transformLiquidityY, utils.Pow96)
	acquireX.Div(acquireX, sqrtPrice_96)
	newLiquidityX := new(big.Int).Sub(liquidityX, transformLiquidityY)
	return Y2XAtPriceLiquidityResult{CostY: costY, AcquireX: acquireX, NewLiquidityX: newLiquidityX}
}

type RangeY2X struct {
	Liquidity     *big.Int
	SqrtPriceL_96 *big.Int
	LeftPt        int
	SqrtPriceR_96 *big.Int
	RightPt       int
	SqrtRate_96   *big.Int
}

type Y2XRangeCompRet struct {
	CostY             *big.Int
	AcquireX          *big.Int
	CompleteLiquidity bool
	LocPt             int
	SqrtLoc_96        *big.Int
}

func y2XRangeComplete(rg RangeY2X, amountY *big.Int) Y2XRangeCompRet {
	ret := Y2XRangeCompRet{}
	maxY := amountmath.GetAmountY(rg.Liquidity, rg.SqrtPriceL_96, rg.SqrtPriceR_96, rg.SqrtRate_96, true)
	if maxY.Cmp(amountY) <= 0 {
		// ret.costY <= maxY <= uint128.max
		ret.CostY = maxY
		ret.AcquireX = amountmath.GetAmountX(rg.Liquidity, rg.LeftPt, rg.RightPt, rg.SqrtPriceR_96, rg.SqrtRate_96, false)
		// we complete this liquidity segment
		ret.CompleteLiquidity = true
	} else {
		// we should locate highest price
		// uint160 is enough for muldiv and adding, because amountY < maxY
		sqrtLoc_96 := calc.MulDivFloor(amountY, new(big.Int).Sub(rg.SqrtRate_96, utils.Pow96), rg.Liquidity)
		sqrtLoc_96.Add(sqrtLoc_96, rg.SqrtPriceL_96)
		ret.LocPt, _ = calc.GetLogSqrtPriceFloor(sqrtLoc_96)

		ret.LocPt = calc.Max(rg.LeftPt, ret.LocPt)
		ret.LocPt = calc.Min(rg.RightPt-1, ret.LocPt)

		ret.CompleteLiquidity = false
		ret.SqrtLoc_96, _ = calc.GetSqrtPrice(ret.LocPt)
		if ret.LocPt == rg.LeftPt {
			ret.CostY = big.NewInt(0)
			ret.AcquireX = big.NewInt(0)
			return ret
		}

		costY256 := amountmath.GetAmountY(rg.Liquidity, rg.SqrtPriceL_96, ret.SqrtLoc_96, rg.SqrtRate_96, true)
		// ret.costY <= amountY <= uint128.max
		ret.CostY = calc.MinBigInt(costY256, amountY)

		// costY <= amountY even if the costY is the upperbound of the result
		// because amountY is not a real and sqrtLoc_96 <= sqrtLoc256_96
		ret.AcquireX = amountmath.GetAmountX(rg.Liquidity, rg.LeftPt, ret.LocPt, ret.SqrtLoc_96, rg.SqrtRate_96, false)

	}
	return ret
}

func Y2XRange(currentState utils.State, rightPt int, sqrtRate_96 *big.Int, amountY *big.Int) Y2XRangeRetState {
	retState := Y2XRangeRetState{
		CostY:      big.NewInt(0),
		AcquireX:   big.NewInt(0),
		Finished:   false,
		LiquidityX: big.NewInt(0),
	}

	// first, if current point is not all x, we can not move right directly
	startHasY := currentState.LiquidityX.Cmp(currentState.Liquidity) < 0
	if startHasY {
		ret := y2XAtPriceLiquidity(amountY, currentState.SqrtPrice_96, currentState.LiquidityX)
		retState.LiquidityX = ret.NewLiquidityX
		retState.CostY = ret.CostY
		retState.AcquireX = ret.AcquireX
		if retState.LiquidityX.Cmp(new(big.Int)) > 0 || retState.CostY.Cmp(amountY) >= 0 {
			retState.Finished = true
			retState.FinalPt = currentState.CurrentPoint
			retState.SqrtFinalPrice_96 = currentState.SqrtPrice_96
			return retState
		} else {
			amountY.Sub(amountY, retState.CostY)
			currentState.CurrentPoint += 1
			if currentState.CurrentPoint == rightPt {
				retState.FinalPt = currentState.CurrentPoint
				retState.SqrtFinalPrice_96, _ = calc.GetSqrtPrice(rightPt)
				return retState
			}
			// sqrt(price) + sqrt(price) * (1.0001 - 1) == sqrt(price) * 1.0001
			mulDelta := new(big.Int).Mul(currentState.SqrtPrice_96, new(big.Int).Sub(sqrtRate_96, utils.Pow96))
			mulDeltaDiv := new(big.Int).Div(mulDelta, utils.Pow96)
			currentState.SqrtPrice_96 = new(big.Int).Add(currentState.SqrtPrice_96, mulDeltaDiv)
		}
	}

	sqrtPriceR_96, _ := calc.GetSqrtPrice(rightPt)

	// (uint128 liquidCostY, uint256 liquidAcquireX, bool liquidComplete, int24 locPt, uint160 sqrtLoc_96)
	ret := y2XRangeComplete(
		RangeY2X{
			Liquidity:     currentState.Liquidity,
			SqrtPriceL_96: currentState.SqrtPrice_96,
			LeftPt:        currentState.CurrentPoint,
			SqrtPriceR_96: sqrtPriceR_96,
			RightPt:       rightPt,
			SqrtRate_96:   sqrtRate_96,
		},
		amountY,
	)

	retState.CostY.Add(retState.CostY, ret.CostY)
	amountY.Sub(amountY, ret.CostY)
	retState.AcquireX.Add(retState.AcquireX, ret.AcquireX)
	if ret.CompleteLiquidity {
		retState.Finished = amountY.Cmp(big.NewInt(0)) == 0
		retState.FinalPt = rightPt
		retState.SqrtFinalPrice_96 = sqrtPriceR_96
	} else {

		//locCostY, locAcquireX, retState.LiquidityX =
		locRet := y2XAtPriceLiquidity(amountY, ret.SqrtLoc_96, currentState.Liquidity)
		locCostY := locRet.CostY
		locAcquireX := locRet.AcquireX
		retState.LiquidityX = locRet.NewLiquidityX

		retState.CostY.Add(retState.CostY, locCostY)
		retState.AcquireX.Add(retState.AcquireX, locAcquireX)
		retState.Finished = true
		retState.SqrtFinalPrice_96 = ret.SqrtLoc_96
		retState.FinalPt = ret.LocPt
	}
	return retState
}
