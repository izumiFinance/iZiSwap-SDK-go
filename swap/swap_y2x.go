package swap

import (
	"fmt"
	"math/big"

	"github.com/izumiFinance/iZiSwap-SDK-go/library/calc"
	"github.com/izumiFinance/iZiSwap-SDK-go/library/swapmath"
	"github.com/izumiFinance/iZiSwap-SDK-go/library/utils"
)

func SwapY2X(amount big.Int, highPt int, pool PoolInfo) (SwapResult, error) {
	if amount.Cmp(big.NewInt(0)) <= 0 {
		return SwapResult{}, fmt.Errorf("AP")
	}

	highPt = calc.Min(highPt, pool.RightMostPt)
	var amountX, amountY big.Int

	sqrtPrice_96, _ := calc.GetSqrtPrice(pool.CurrentPoint)

	liquidityX := pool.LiquidityX
	liquidity := pool.Liquidity

	finished := false
	sqrtRate_96, _ := calc.GetSqrtPrice(1)
	pointDelta := pool.PointDelta
	currentPoint := pool.CurrentPoint
	fee := pool.Fee

	orderData := InitY2X(
		pool.Liquidities,
		pool.LimitOrders,
		pool.CurrentPoint,
	)

	for currentPoint < highPt && !finished {
		if orderData.IsLimitOrder(currentPoint) {
			// amount <= uint128.max
			amountNoFee := new(big.Int).Mul(&amount, big.NewInt(int64(1e6-fee)))
			amountNoFee.Div(amountNoFee, big.NewInt(1e6))
			if amountNoFee.Cmp(big.NewInt(0)) > 0 {
				// clear limit order first
				currX := orderData.UnsafeGetLimitSellingX()
				costY, acquireX := swapmath.Y2XAtPrice(amountNoFee, sqrtPrice_96, &currX)
				if acquireX.Cmp(&currX) < 0 || costY.Cmp(amountNoFee) >= 0 {
					finished = true
				}
				var feeAmount *big.Int
				if costY.Cmp(amountNoFee) >= 0 {
					feeAmount = new(big.Int).Sub(&amount, costY)
				} else {
					// amount <= uint128.max
					feeAmount = new(big.Int).Mul(costY, big.NewInt(int64(fee)))
					feeAmount.Div(feeAmount, big.NewInt(int64(1e6-fee)))
					mod := new(big.Int).Mod(new(big.Int).Mul(costY, big.NewInt(int64(fee))), big.NewInt(int64(1e6-fee)))
					if mod.Cmp(big.NewInt(0)) > 0 {
						feeAmount.Add(feeAmount, big.NewInt(1))
					}
				}
				amount.Sub(&amount, new(big.Int).Add(costY, feeAmount))
				amountY.Add(&amountY, new(big.Int).Add(costY, feeAmount))
				amountX.Add(&amountX, acquireX)
				orderData.ConsumeLimitOrder(true)
			} else {
				finished = true
			}
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
				liquidity.Add(&liquidity, &delta)
				liquidityX = liquidity
			}
		} else {
			// amount <= uint128.max
			amountNoFee := new(big.Int).Mul(&amount, big.NewInt(int64(1e6-fee)))
			amountNoFee.Div(amountNoFee, big.NewInt(1e6))
			if amountNoFee.Cmp(big.NewInt(0)) > 0 {
				st := utils.State{
					LiquidityX:   &liquidityX,
					Liquidity:    &liquidity,
					CurrentPoint: currentPoint,
					SqrtPrice_96: sqrtPrice_96,
				}
				retState := swapmath.Y2XRange(st, nextPoint, sqrtRate_96, new(big.Int).Set(amountNoFee))

				finished = retState.Finished
				var feeAmount *big.Int
				if retState.CostY.Cmp(amountNoFee) >= 0 {
					feeAmount = new(big.Int).Sub(&amount, retState.CostY)
				} else {
					// retState.costY <= uint128.max
					feeAmount = new(big.Int).Mul(retState.CostY, big.NewInt(int64(fee)))
					feeAmount.Div(feeAmount, big.NewInt(int64(1e6-fee)))
					mod := new(big.Int).Mod(new(big.Int).Mul(retState.CostY, big.NewInt(int64(fee))), big.NewInt(int64(1e6-fee)))
					if mod.Cmp(big.NewInt(0)) > 0 {
						feeAmount.Add(feeAmount, big.NewInt(1))
					}
				}

				amountX.Add(&amountX, retState.AcquireX)
				amountY.Add(&amountY, new(big.Int).Add(retState.CostY, feeAmount))
				amount.Sub(&amount, new(big.Int).Add(retState.CostY, feeAmount))

				currentPoint = retState.FinalPt
				sqrtPrice_96 = retState.SqrtFinalPrice_96
				liquidityX = *(retState.LiquidityX)
			} else {
				finished = true
			}

			if currentPoint == nextPoint {
				if orderData.IsLiquidity(nextPoint) {
					delta := orderData.UnsafeGetDeltaLiquidity()
					liquidity.Add(&liquidity, &delta)
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
