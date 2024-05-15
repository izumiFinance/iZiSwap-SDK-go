package swap

import (
	"fmt"
	"math/big"

	"github.com/izumiFinance/iZiSwap-SDK-go/library/calc"
	"github.com/izumiFinance/iZiSwap-SDK-go/library/swapmathdesire"
	"github.com/izumiFinance/iZiSwap-SDK-go/library/utils"
)

func SwapY2XDesireX(
	desireX *big.Int,
	highPt int,
	pool PoolInfo,
) (SwapResult, error) {
	if desireX.Cmp(big.NewInt(0)) <= 0 {
		return SwapResult{}, fmt.Errorf("AP")
	}
	desireX = new(big.Int).Set(desireX)

	highPt = calc.Min(highPt, pool.RightMostPt)

	amountX := big.NewInt(0)
	amountY := big.NewInt(0)

	sqrtPrice_96, _ := calc.GetSqrtPrice(pool.CurrentPoint)

	liquidityX := new(big.Int).Set(pool.LiquidityX)
	liquidity := new(big.Int).Set(pool.Liquidity)

	finished := false
	sqrtRate_96, _ := calc.GetSqrtPrice(1)
	pointDelta := pool.PointDelta
	currentPoint := pool.CurrentPoint
	fee := int64(pool.Fee)

	orderData := InitY2X(
		pool.Liquidities,
		pool.LimitOrders,
		pool.CurrentPoint,
	)

	for currentPoint < highPt && !finished {
		if orderData.IsLimitOrder(currentPoint) {
			// clear limit order first
			currX := orderData.UnsafeGetLimitSellingX()
			costY, acquireX := swapmathdesire.Y2XAtPrice(desireX, sqrtPrice_96, currX)

			if acquireX.Cmp(desireX) >= 0 {
				finished = true
			}
			feeAmount := calc.MulDivCeil(
				costY,
				new(big.Int).SetInt64(fee),
				new(big.Int).SetInt64(1e6-fee),
			)

			desireX.Sub(desireX, acquireX)
			amountY.Add(amountY, new(big.Int).Add(costY, feeAmount))
			amountX.Add(amountX, acquireX)
			orderData.ConsumeLimitOrder(true)
		}

		if finished {
			break
		}

		nextPoint := orderData.MoveY2X(currentPoint, pointDelta)
		if nextPoint > highPt {
			nextPoint = highPt
		}
		// in [st.currentPoint, nextPoint)
		if liquidity.Cmp(big.NewInt(0)) == 0 {
			// no liquidity in the range [st.currentPoint, nextPoint)
			currentPoint = nextPoint
			sqrtPrice_96, _ = calc.GetSqrtPrice(currentPoint)
			if orderData.IsLiquidity(currentPoint) {
				delta := orderData.UnsafeGetDeltaLiquidity()
				liquidity.Add(liquidity, delta)
				liquidityX = liquidity
			}
		} else {
			// desireX > 0
			if desireX.Cmp(big.NewInt(0)) > 0 {
				st := utils.State{
					LiquidityX:   new(big.Int).Set(liquidityX),
					Liquidity:    new(big.Int).Set(liquidity),
					CurrentPoint: currentPoint,
					SqrtPrice_96: sqrtPrice_96,
				}
				retState := swapmathdesire.Y2XRange(st, nextPoint, sqrtRate_96, desireX)

				finished = retState.Finished
				feeAmount := calc.MulDivCeil(
					retState.CostY,
					new(big.Int).SetInt64(fee),
					new(big.Int).SetInt64(1e6-fee),
				)

				amountX.Add(amountX, retState.AcquireX)
				amountY.Add(amountY, new(big.Int).Add(retState.CostY, feeAmount))
				desireX.Sub(desireX, retState.AcquireX)

				currentPoint = retState.FinalPt
				sqrtPrice_96 = retState.SqrtFinalPrice_96
				liquidityX = retState.LiquidityX
			} else {
				finished = true
			}

			if currentPoint == nextPoint {
				if orderData.IsLiquidity(nextPoint) {
					delta := orderData.UnsafeGetDeltaLiquidity()
					liquidity.Add(liquidity, delta)
				}
				liquidityX = liquidity
			}
		}
	}

	swapResult := SwapResult{
		CurrentPoint: currentPoint,
		Liquidity:    liquidity,
		LiquidityX:   liquidityX,
		AmountX:      amountX,
		AmountY:      amountY,
	}
	return swapResult, nil
}
