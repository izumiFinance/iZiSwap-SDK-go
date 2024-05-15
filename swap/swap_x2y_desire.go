package swap

import (
	"fmt"
	"math/big"

	"github.com/izumiFinance/iZiSwap-SDK-go/library/calc"
	"github.com/izumiFinance/iZiSwap-SDK-go/library/swapmathdesire"
	"github.com/izumiFinance/iZiSwap-SDK-go/library/utils"
)

func SwapX2YDesireY(
	desireY *big.Int,
	lowPt int,
	pool PoolInfo,
) (SwapResult, error) {
	if desireY.Cmp(big.NewInt(0)) <= 0 {
		return SwapResult{}, fmt.Errorf("AP")
	}
	desireY = new(big.Int).Set(desireY)

	lowPt = calc.Max(lowPt, pool.LeftMostPt)
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

	orderData := InitX2Y(
		pool.Liquidities,
		pool.LimitOrders,
		pool.CurrentPoint,
	)

	for lowPt <= currentPoint && !finished {
		// clear limit order first
		if orderData.IsLimitOrder(currentPoint) {
			currY := orderData.UnsafeGetLimitSellingY()

			costX, acquireY := swapmathdesire.X2YAtPrice(desireY, sqrtPrice_96, currY)

			if acquireY.Cmp(desireY) >= 0 {
				finished = true
			}

			feeAmount := calc.MulDivCeil(
				costX,
				new(big.Int).SetInt64(fee),
				new(big.Int).SetInt64(1e6-fee),
			)
			desireY.Sub(desireY, acquireY)
			amountX.Add(amountX, costX)
			amountX.Add(amountX, feeAmount)
			amountY.Add(amountY, acquireY)
			orderData.ConsumeLimitOrder(false)
		}
		if finished {
			break
		}
		searchStart := currentPoint - 1
		// second, clear the liquid if the currentPoint is an endpoint
		if orderData.IsLiquidity(currentPoint) {
			if liquidity.Cmp(big.NewInt(0)) > 0 {
				st := utils.State{
					LiquidityX:   new(big.Int).Set(liquidityX),
					Liquidity:    new(big.Int).Set(liquidity),
					CurrentPoint: currentPoint,
					SqrtPrice_96: sqrtPrice_96,
				}
				retState := swapmathdesire.X2YRange(st, currentPoint, sqrtRate_96, desireY)
				finished = retState.Finished

				feeAmount := calc.MulDivCeil(
					retState.CostX,
					new(big.Int).SetInt64(fee),
					new(big.Int).SetInt64(1e6-fee),
				)

				amountX.Add(amountX, retState.CostX)
				amountX.Add(amountX, feeAmount)
				amountY.Add(amountY, retState.AcquireY)
				desireY.Sub(desireY, retState.AcquireY)

				currentPoint = retState.FinalPt
				sqrtPrice_96 = retState.SqrtFinalPrice_96
				liquidityX = retState.LiquidityX
			}
			if !finished {
				delta := orderData.UnsafeGetDeltaLiquidity()
				liquidity.Sub(liquidity, delta)
				currentPoint -= 1
				sqrtPrice_96, _ = calc.GetSqrtPrice(currentPoint)
				liquidityX.SetInt64(0)
			}
		}
		if finished || currentPoint < lowPt {
			break
		}

		nextPt := orderData.MoveX2Y(searchStart, pointDelta)
		if nextPt < lowPt {
			nextPt = lowPt
		}

		// in [nextPt, st.currentPoint)
		if liquidity.Cmp(big.NewInt(0)) == 0 {
			// no liquidity in the range [nextPt, st.currentPoint]
			currentPoint = nextPt
			sqrtPrice_96, _ = calc.GetSqrtPrice(currentPoint)
			// liquidityX must be 0
		} else {
			st := utils.State{
				LiquidityX:   new(big.Int).Set(liquidityX),
				Liquidity:    new(big.Int).Set(liquidity),
				CurrentPoint: currentPoint,
				SqrtPrice_96: sqrtPrice_96,
			}
			retState := swapmathdesire.X2YRange(
				st, nextPt, sqrtRate_96, desireY,
			)
			finished = retState.Finished

			feeAmount := calc.MulDivCeil(
				retState.CostX,
				new(big.Int).SetInt64(fee),
				new(big.Int).SetInt64(1e6-fee),
			)

			amountY.Add(amountY, retState.AcquireY)
			amountX.Add(amountX, retState.CostX)
			amountX.Add(amountX, feeAmount)
			desireY.Sub(desireY, retState.AcquireY)

			currentPoint = retState.FinalPt
			sqrtPrice_96 = retState.SqrtFinalPrice_96
			liquidityX = retState.LiquidityX
		}

		if currentPoint <= lowPt {
			break
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
