# iZiSwap-SDK-go

### description

this sdk provide offline-interfaces to calculate exchange token amount on iZiSwap.

with providing distribution of `liquidity` and `limit order` and `currentPoint` of pool,
you can call following interfaces to compute exchange amount of `tokenX` and `tokenY`

```
// swap from tokenY to tokenX with input amount of tokenY
func SwapY2X(amount big.Int, highPt int, pool PoolInfo)

// swap from tokenY to tokenX with desired amount of tokenX
func SwapY2XDesireX(amount big.Int, highPt int, pool PoolInfo)

// swap from tokenX to tokenY with input amount of tokenX
func SwapX2Y(amount big.Int, lowPt int, pool PoolInfo) 

// swap from tokenX to tokenY with desired amount of tokenY
func SwapX2YDesireY(amount big.Int, lowPt int, pool PoolInfo) 

```

here when a pair is (tokenA, tokenB),
if address(tokenA).LowerCase() < address(tokenB).LowerCase()
then, tokenA is tokenX, tokenB is tokenY,
otherwise, tokenB is tokenX, tokenA is tokenY.

for more detail of usage, you can refer to `example.go.txt`
