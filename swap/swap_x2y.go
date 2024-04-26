package swap

import (
	"fmt"
	"math/big"

	"github.com/KyberNetwork/iZiSwap-SDK-go/library/calc"
	"github.com/KyberNetwork/iZiSwap-SDK-go/library/swapmath"
	"github.com/KyberNetwork/iZiSwap-SDK-go/library/swapmathdesire"
	"github.com/KyberNetwork/iZiSwap-SDK-go/library/utils"
)

func SwapX2Y(amount *big.Int, lowPt int, pool PoolInfo) (SwapResult, error) {
	if amount.Cmp(big.NewInt(0)) <= 0 {
		return SwapResult{}, fmt.Errorf("AP")
	}

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
	fee := pool.Fee

	orderData := InitX2Y(
		pool.Liquidities,
		pool.LimitOrders,
		pool.CurrentPoint,
	)

	for lowPt <= currentPoint && !finished {
		if orderData.IsLimitOrder(currentPoint) {
			// amount <= uint128.max
			amountNoFee := new(big.Int).Mul(amount, big.NewInt(int64(1e6-fee)))
			amountNoFee.Div(amountNoFee, big.NewInt(1e6))
			if amountNoFee.Cmp(big.NewInt(0)) > 0 {

				currY := orderData.UnsafeGetLimitSellingY()
				costX, acquireY := swapmath.X2YAtPrice(amountNoFee, sqrtPrice_96, currY)

				if acquireY.Cmp(currY) < 0 || costX.Cmp(amountNoFee) >= 0 {
					finished = true
				}

				feeAmount := new(big.Int)
				if costX.Cmp(amountNoFee) >= 0 {
					feeAmount.Sub(amount, costX)
				} else {
					// costX <= amountX <= uint128.max
					feeAmount.Mul(costX, big.NewInt(int64(fee)))
					feeAmount.Div(feeAmount, big.NewInt(int64(1e6-fee)))
					mod := new(big.Int).Mul(costX, big.NewInt(int64(fee)))
					mod.Mod(mod, big.NewInt(int64(1e6-fee)))
					if mod.Cmp(big.NewInt(0)) > 0 {
						feeAmount.Add(feeAmount, big.NewInt(1))
					}
				}

				amount.Sub(amount, costX)
				amount.Sub(amount, feeAmount)
				amountX.Add(amountX, costX)
				amountX.Add(amountX, feeAmount)
				amountY.Add(amountY, acquireY)

				orderData.ConsumeLimitOrder(false)
			} else {
				finished = true
			}
		}

		if finished {
			break
		}

		searchStart := currentPoint - 1

		// step2: clear the liquidity if the currentPoint is an endpoint
		if orderData.IsLiquidity(currentPoint) {
			amountNoFee := new(big.Int).Mul(amount, big.NewInt(int64(1e6-fee)))
			amountNoFee.Div(amountNoFee, big.NewInt(int64(1e6)))
			if amountNoFee.Cmp(big.NewInt(0)) > 0 {
				if liquidity.Cmp(big.NewInt(0)) > 0 {
					st := utils.State{
						LiquidityX:   new(big.Int).Set(liquidityX),
						Liquidity:    new(big.Int).Set(liquidity),
						CurrentPoint: currentPoint,
						SqrtPrice_96: sqrtPrice_96,
					}
					retState := swapmath.X2YRange(st, currentPoint, sqrtRate_96, new(big.Int).Set(amountNoFee))
					finished = retState.Finished

					feeAmount := new(big.Int)
					if retState.CostX.Cmp(amountNoFee) >= 0 {
						feeAmount.Sub(amount, retState.CostX)
					} else {
						feeAmount.Mul(retState.CostX, big.NewInt(int64(fee)))
						feeAmount.Div(feeAmount, big.NewInt(int64(1e6-fee)))
						mod := new(big.Int).Mul(retState.CostX, big.NewInt(int64(fee)))
						mod.Mod(mod, big.NewInt(int64(1e6-fee)))
						if mod.Cmp(big.NewInt(0)) > 0 {
							feeAmount.Add(feeAmount, big.NewInt(1))
						}
					}

					amountX.Add(amountX, retState.CostX)
					amountX.Add(amountX, feeAmount)
					amountY.Add(amountY, retState.AcquireY)
					amount.Sub(amount, retState.CostX)
					amount.Sub(amount, feeAmount)
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
			} else {
				finished = true
			}
		}

		if finished || currentPoint < lowPt {
			break
		}

		nextPt := orderData.MoveX2Y(searchStart, pointDelta)
		if nextPt < lowPt {
			nextPt = lowPt
		}

		if liquidity.Cmp(big.NewInt(0)) == 0 {
			// no liquidity in the range [nextPt, st.currentPoint]
			currentPoint = nextPt
			sqrtPrice_96, _ = calc.GetSqrtPrice(currentPoint)
		} else {
			amountNoFee := new(big.Int).Mul(amount, big.NewInt(int64(1e6-fee)))
			amountNoFee.Div(amountNoFee, big.NewInt(int64(1e6)))
			if amountNoFee.Cmp(big.NewInt(0)) > 0 {
				st := utils.State{
					LiquidityX:   new(big.Int).Set(liquidityX),
					Liquidity:    new(big.Int).Set(liquidity),
					CurrentPoint: currentPoint,
					SqrtPrice_96: sqrtPrice_96,
				}
				retState := swapmath.X2YRange(st, nextPt, sqrtRate_96, new(big.Int).Set(amountNoFee))
				finished = retState.Finished
				feeAmount := new(big.Int)
				if retState.CostX.Cmp(amountNoFee) >= 0 {
					feeAmount.Sub(amount, retState.CostX)
				} else {
					feeAmount.Mul(retState.CostX, big.NewInt(int64(fee)))
					feeAmount.Div(feeAmount, big.NewInt(int64(1e6-fee)))
					mod := new(big.Int).Mul(retState.CostX, big.NewInt(int64(fee)))
					mod.Mod(mod, big.NewInt(int64(1e6-fee)))
					if mod.Cmp(big.NewInt(0)) > 0 {
						feeAmount.Add(feeAmount, big.NewInt(1))
					}
				}
				amountY.Add(amountY, retState.AcquireY)
				amountX.Add(amountX, retState.CostX)
				amountX.Add(amountX, feeAmount)
				amount.Sub(amount, retState.CostX)
				amount.Sub(amount, feeAmount)

				currentPoint = retState.FinalPt
				sqrtPrice_96 = retState.SqrtFinalPrice_96
				liquidityX = retState.LiquidityX
			} else {
				finished = true
			}
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

func SwapX2YDesireY(desireY *big.Int, lowPt int, pool PoolInfo) (SwapResult, error) {
	if desireY.Cmp(big.NewInt(0)) <= 0 {
		return SwapResult{}, fmt.Errorf("AP")
	}

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
	fee := pool.Fee

	orderData := InitX2Y(
		pool.Liquidities,
		pool.LimitOrders,
		pool.CurrentPoint,
	)

	for lowPt <= currentPoint && !finished {
		if orderData.IsLimitOrder(currentPoint) {
			currY := orderData.UnsafeGetLimitSellingY()
			costX, acquireY := swapmathdesire.X2YAtPrice(desireY, sqrtPrice_96, currY)

			if acquireY.Cmp(desireY) >= 0 {
				finished = true
			}

			feeAmount := calc.MulDivCeil(costX, big.NewInt(int64(fee)), big.NewInt(1e6))

			if desireY.Cmp(acquireY) <= 0 {
				desireY = zeroBI
			} else {
				desireY.Sub(desireY, acquireY)
			}

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
				retState := swapmathdesire.X2YRange(st, currentPoint, sqrtRate_96, new(big.Int).Set(desireY))
				finished = retState.Finished

				feeAmount := calc.MulDivCeil(retState.CostX, big.NewInt(int64(fee)), big.NewInt(int64(1e6-fee)))

				amountX.Add(amountX, retState.CostX)
				amountX.Add(amountX, feeAmount)
				amountY.Add(amountY, retState.AcquireY)
				desireY.Sub(desireY, calc.MinBigInt(desireY, retState.AcquireY))
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
		} else {
			st := utils.State{
				LiquidityX:   new(big.Int).Set(liquidityX),
				Liquidity:    new(big.Int).Set(liquidity),
				CurrentPoint: currentPoint,
				SqrtPrice_96: sqrtPrice_96,
			}
			retState := swapmathdesire.X2YRange(st, nextPt, sqrtRate_96, new(big.Int).Set(desireY))
			finished = retState.Finished
			feeAmount := calc.MulDivCeil(retState.CostX, big.NewInt(int64(fee)), big.NewInt(int64(1e6-fee)))
			amountY.Add(amountY, retState.AcquireY)
			amountX.Add(amountX, retState.CostX)
			amountX.Add(amountX, feeAmount)
			desireY.Sub(desireY, calc.MinBigInt(desireY, retState.AcquireY))

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
