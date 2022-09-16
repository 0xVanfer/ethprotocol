package lend

type LendPoolParams struct {
	CollateralFactor   float64
	LiquidationLimit   float64
	LiquidationPenalty float64
	AllowBorrow        int
	AllowCollateral    int
}
