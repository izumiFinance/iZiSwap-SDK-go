package amountmath

import (
	"math/big"

	"github.com/KyberNetwork/iZiSwap-SDK-go/library/calc"
	"github.com/KyberNetwork/iZiSwap-SDK-go/library/utils"
)

func GetAmountY(
	liquidity *big.Int,
	sqrtPriceL_96 *big.Int,
	sqrtPriceR_96 *big.Int,
	sqrtRate_96 *big.Int,
	upper bool,
) *big.Int {
	var amount *big.Int
	numerator := new(big.Int).Sub(sqrtPriceR_96, sqrtPriceL_96)
	denominator := new(big.Int).Sub(sqrtRate_96, utils.Pow96)
	if !upper {
		// You should replace MulDivMath.mulDivFloor with equivalent Go function
		amount = calc.MulDivFloor(liquidity, numerator, denominator)
	} else {
		// You should replace MulDivMath.mulDivCeil with equivalent Go function
		amount = calc.MulDivCeil(liquidity, numerator, denominator)
	}
	return amount
}

func GetAmountX(
	liquidity *big.Int,
	leftPt int,
	rightPt int,
	sqrtPriceR_96 *big.Int,
	sqrtRate_96 *big.Int,
	upper bool,
) *big.Int {
	var amount *big.Int
	// You should replace LogPowMath.getSqrtPrice with equivalent Go function
	sqrtPricePrPl_96, _ := calc.GetSqrtPrice(rightPt - leftPt)

	temp := new(big.Int).Mul(sqrtPriceR_96, utils.Pow96)
	sqrtPricePrM1_96 := new(big.Int).Div(temp, sqrtRate_96)

	numerator := new(big.Int).Sub(sqrtPricePrPl_96, utils.Pow96)
	denominator := new(big.Int).Sub(sqrtPriceR_96, sqrtPricePrM1_96)
	if !upper {
		// You should replace MulDivMath.mulDivFloor with equivalent Go function
		amount = calc.MulDivFloor(liquidity, numerator, denominator)
	} else {
		// You should replace MulDivMath.mulDivCeil with equivalent Go function
		amount = calc.MulDivCeil(liquidity, numerator, denominator)
	}
	return amount
}
