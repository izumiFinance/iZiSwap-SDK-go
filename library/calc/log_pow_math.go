package calc

import (
	"errors"
	"fmt"
	"math"
	"math/big"
)

const (
	MIN_POINT      = -887272
	MAX_POINT      = -MIN_POINT
	MIN_SQRT_PRICE = "4295128739"
	MAX_SQRT_PRICE = "1461446703485210103287273052203988822378723970342"
)

func GetSqrtPrice(point int) (*big.Int, error) {
	if point < MIN_POINT || point > MAX_POINT {
		return nil, errors.New("T")
	}
	absIdx := point
	if absIdx < 0 {
		absIdx = -absIdx
	}

	value := new(big.Int)
	if (absIdx & 1) != 0 {
		value.SetString("fffcb933bd6fad37aa2d162d1a594001", 16)
	} else {
		value.SetString("100000000000000000000000000000000", 16)
	}

	cases := []string{
		"fff97272373d413259a46990580e213a",
		"fff2e50f5f656932ef12357cf3c7fdcc",
		"ffe5caca7e10e4e61c3624eaa0941cd0",
		"ffcb9843d60f6159c9db58835c926644",
		"ff973b41fa98c081472e6896dfb254c0",
		"ff2ea16466c96a3843ec78b326b52861",
		"fe5dee046a99a2a811c461f1969c3053",
		"fcbe86c7900a88aedcffc83b479aa3a4",
		"f987a7253ac413176f2b074cf7815e54",
		"f3392b0822b70005940c7a398e4b70f3",
		"e7159475a2c29b7443b29c7fa6e889d9",
		"d097f3bdfd2022b8845ad8f792aa5825",
		"a9f746462d870fdf8a65dc1f90e061e5",
		"70d869a156d2a1b890bb3df62baf32f7",
		"31be135f97d08fd981231505542fcfa6",
		"9aa508b5b7a84e1c677de54f3e99bc9",
		"5d6af8dedb81196699c329225ee604",
		"2216e584f5fa1ea926041bedfe98",
		"48a170391f7dc42444e8fa2",
	}

	for i, c := range cases {
		if (absIdx & (1 << (i + 1))) != 0 {
			multiplier := new(big.Int)
			multiplier.SetString(c, 16)
			value.Mul(value, multiplier)
			value.Rsh(value, 128)
		}
	}

	if point > 0 {
		maxUint256 := new(big.Int)
		maxUint256.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
		value.Div(maxUint256, value)
	}

	sqrtPrice_96 := new(big.Int)
	sqrtPrice_96.Rsh(value, 32)
	if value.Mod(value, big.NewInt(int64(math.Pow(2, 32)))).Cmp(new(big.Int)) != 0 {
		sqrtPrice_96.Add(sqrtPrice_96, big.NewInt(1))
	}

	return sqrtPrice_96, nil
}

func GetLogSqrtPriceFloor(sqrtPrice_96 *big.Int) (int, error) {
	minSqrtPrice := new(big.Int)
	minSqrtPrice.SetString(MIN_SQRT_PRICE, 10)
	maxSqrtPrice := new(big.Int)
	maxSqrtPrice.SetString(MAX_SQRT_PRICE, 10)
	sqrtPrice_96_big := sqrtPrice_96

	if sqrtPrice_96_big.Cmp(minSqrtPrice) <= 0 || sqrtPrice_96_big.Cmp(maxSqrtPrice) >= 0 {
		return 0, fmt.Errorf("R")
	}

	sqrtPrice_128 := new(big.Int).Lsh(sqrtPrice_96_big, 32)
	x := new(big.Int).Set(sqrtPrice_128)
	m := new(big.Int)

	bitSize := []uint{128, 64, 32, 16, 8, 4, 2, 1}
	th := []string{
		"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
		"FFFFFFFFFFFFFFFF",
		"FFFFFFFF",
		"FFFF",
		"FF",
		"F",
		"3",
		"1",
	}
	for idx, size := range bitSize {
		currTh := new(big.Int)
		currTh.SetString(th[idx], 16)
		if x.Cmp(currTh) > 0 {
			m.Or(m, new(big.Int).SetUint64(uint64(size)))
			if size > 1 {
				x.Rsh(x, uint(size))
			}
		}
	}

	if m.Cmp(big.NewInt(128)) >= 0 {
		x.Rsh(sqrtPrice_128, uint(m.Int64()-127))
	} else {
		x.Lsh(sqrtPrice_128, uint(127-m.Int64()))
	}

	l2 := new(big.Int).Lsh(new(big.Int).Sub(m, big.NewInt(128)), 64)

	// Simulate the assembly code
	for i := 63; i >= 50; i-- {
		x.Mul(x, x)
		x.Rsh(x, 127)
		y := new(big.Int).Rsh(x, 128)
		l2 = l2.Or(l2, new(big.Int).Lsh(y, uint(i)))
		if i > 50 {
			x.Rsh(x, uint(y.Uint64()))
		}
	}

	bigIntValue, _ := new(big.Int).SetString("255738958999603826347141", 10)
	ls10001 := new(big.Int).Mul(l2, bigIntValue)

	bigIntValueF, _ := new(big.Int).SetString("3402992956809132418596140100660247210", 10)
	logFloor := new(big.Int).Rsh(new(big.Int).Sub(ls10001, bigIntValueF), 128)

	bigIntValueL, _ := new(big.Int).SetString("291339464771989622907027621153398088495", 10)
	logUpper := new(big.Int).Rsh(new(big.Int).Add(ls10001, bigIntValueL), 128)

	logValue := logFloor

	sqrtPrice, _ := GetSqrtPrice(int(logUpper.Int64()))
	if sqrtPrice.Cmp(sqrtPrice_96_big) <= 0 {
		logValue = logUpper
	}

	return int(logValue.Int64()), nil
}
