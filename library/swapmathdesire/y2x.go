package swapmathdesire

import (
	"math/big"

	"github.com/izumiFinance/iZiSwap-SDK-go/library/amountmath"
	"github.com/izumiFinance/iZiSwap-SDK-go/library/calc"
	"github.com/izumiFinance/iZiSwap-SDK-go/library/utils"
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

func Y2XAtPrice(desireX *big.Int, sqrtPrice_96 *big.Int, currX *big.Int) (costY, acquireX *big.Int) {
	acquireX = calc.MinBigInt(desireX, currX)
	l := calc.MulDivCeil(acquireX, sqrtPrice_96, utils.Pow96)
	costY = calc.MulDivCeil(l, sqrtPrice_96, utils.Pow96)
	return costY, acquireX
}

type Y2XAtPriceLiquidityResult struct {
	CostY         *big.Int
	AcquireX      *big.Int
	NewLiquidityX *big.Int
}

func y2XAtPriceLiquidity(amountY *big.Int, sqrtPrice_96 *big.Int, liquidityX *big.Int) Y2XAtPriceLiquidityResult {
	maxTransformLiquidityY := calc.MulDivCeil(amountY, sqrtPrice_96, utils.Pow96)
	transformLiquidityY := calc.MinBigInt(maxTransformLiquidityY, liquidityX)
	costY := calc.MulDivCeil(transformLiquidityY, sqrtPrice_96, utils.Pow96)
	acquireX := new(big.Int).Div(
		new(big.Int).Mul(
			transformLiquidityY,
			utils.Pow96,
		),
		sqrtPrice_96,
	)
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

func y2XRangeComplete(rg RangeY2X, desireX *big.Int) Y2XRangeCompRet {
	ret := Y2XRangeCompRet{}

	maxX := amountmath.GetAmountX(rg.Liquidity, rg.LeftPt, rg.RightPt, rg.SqrtPriceR_96, rg.SqrtRate_96, false)
	if maxX.Cmp(desireX) <= 0 {
		ret.AcquireX = new(big.Int).Set(maxX)
		ret.CostY = amountmath.GetAmountY(rg.Liquidity, rg.SqrtPriceL_96, rg.SqrtPriceR_96, rg.SqrtRate_96, true)
		ret.CompleteLiquidity = true

		return ret
	}

	sqrtPricePrPl_96, _ := calc.GetSqrtPrice(rg.RightPt - rg.LeftPt)
	sqrtPricePrM1_96 := new(big.Int).Div(
		new(big.Int).Mul(rg.SqrtPriceR_96, utils.Pow96),
		rg.SqrtRate_96,
	)
	div := new(big.Int).Sub(
		sqrtPricePrPl_96,
		calc.MulDivFloor(desireX, new(big.Int).Sub(rg.SqrtPriceR_96, sqrtPricePrM1_96), rg.Liquidity),
	)
	sqrtPriceLoc_96 := new(big.Int).Div(
		new(big.Int).Mul(rg.SqrtPriceR_96, utils.Pow96),
		div,
	)

	ret.CompleteLiquidity = false
	ret.LocPt, _ = calc.GetLogSqrtPriceFloor(sqrtPriceLoc_96)
	ret.LocPt = calc.Max(rg.LeftPt, ret.LocPt)
	ret.LocPt = calc.Min(rg.RightPt-1, ret.LocPt)
	ret.SqrtLoc_96, _ = calc.GetSqrtPrice(ret.LocPt)

	if ret.LocPt == rg.LeftPt {
		ret.AcquireX = big.NewInt(0)
		ret.CostY = big.NewInt(0)
		return ret
	}

	ret.CompleteLiquidity = false
	ret.AcquireX = calc.MinBigInt(
		amountmath.GetAmountX(
			rg.Liquidity,
			rg.LeftPt,
			ret.LocPt,
			ret.SqrtLoc_96,
			rg.SqrtRate_96,
			false,
		),
		desireX,
	)

	ret.CostY = amountmath.GetAmountY(
		rg.Liquidity,
		rg.SqrtPriceL_96,
		ret.SqrtLoc_96,
		rg.SqrtRate_96,
		true,
	)

	return ret
}

func Y2XRange(currentState utils.State, rightPt int, sqrtRate_96 *big.Int, desireX *big.Int) Y2XRangeRetState {
	retState := Y2XRangeRetState{
		CostY:      big.NewInt(0),
		AcquireX:   big.NewInt(0),
		Finished:   false,
		LiquidityX: big.NewInt(0),
	}

	startHasY := currentState.LiquidityX.Cmp(currentState.Liquidity) < 0
	if startHasY {
		ret := y2XAtPriceLiquidity(desireX, currentState.SqrtPrice_96, currentState.LiquidityX)
		retState.LiquidityX = ret.NewLiquidityX
		retState.CostY = ret.CostY
		retState.AcquireX = ret.AcquireX
		if retState.LiquidityX.Cmp(new(big.Int)) > 0 || retState.CostY.Cmp(desireX) >= 0 {
			retState.Finished = true
			retState.FinalPt = currentState.CurrentPoint
			retState.SqrtFinalPrice_96 = currentState.SqrtPrice_96
			return retState
		} else {
			desireX.Sub(desireX, retState.AcquireX)
			currentState.CurrentPoint += 1
			if currentState.CurrentPoint == rightPt {
				retState.FinalPt = currentState.CurrentPoint
				retState.SqrtFinalPrice_96, _ = calc.GetSqrtPrice(rightPt)
				return retState
			}
			mulDelta := new(big.Int).Mul(currentState.SqrtPrice_96, new(big.Int).Sub(sqrtRate_96, utils.Pow96))
			mulDeltaDiv := new(big.Int).Div(mulDelta, utils.Pow96)
			currentState.SqrtPrice_96 = new(big.Int).Add(currentState.SqrtPrice_96, mulDeltaDiv)
		}
	}

	sqrtPriceR_96, _ := calc.GetSqrtPrice(rightPt)
	ret := y2XRangeComplete(
		RangeY2X{
			Liquidity:     currentState.Liquidity,
			SqrtPriceL_96: currentState.SqrtPrice_96,
			LeftPt:        currentState.CurrentPoint,
			SqrtPriceR_96: sqrtPriceR_96,
			RightPt:       rightPt,
			SqrtRate_96:   sqrtRate_96,
		},
		desireX,
	)
	retState.CostY.Add(retState.CostY, ret.CostY)
	retState.AcquireX.Add(retState.AcquireX, ret.AcquireX)
	desireX.Sub(desireX, ret.AcquireX)

	if ret.CompleteLiquidity {
		retState.Finished = desireX.Cmp(big.NewInt(0)) == 0
		retState.FinalPt = rightPt
		retState.SqrtFinalPrice_96 = sqrtPriceR_96
	} else {
		locRet := y2XAtPriceLiquidity(desireX, ret.SqrtLoc_96, currentState.Liquidity)
		locCostY := locRet.CostY
		locAcquireX := locRet.AcquireX
		retState.LiquidityX = locRet.NewLiquidityX
		retState.CostY.Add(retState.CostY, locCostY)
		retState.AcquireX.Add(retState.AcquireX, locAcquireX)
		retState.Finished = true
		retState.FinalPt = ret.LocPt
		retState.SqrtFinalPrice_96 = ret.SqrtLoc_96
	}

	return retState
}
