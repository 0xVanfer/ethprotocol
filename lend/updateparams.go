package lend

import "errors"

type LendPoolParams struct {
	CollateralFactor   float64
	LiquidationLimit   float64
	LiquidationPenalty float64
	AllowBorrow        int
	AllowCollateral    int
}

func (p *LendPool) UpdatePoolParams(params *LendPoolParams) error {
	if !p.ProtocolBasic.Regularcheck() {
		return errors.New("protocol basic must be initialized")
	}
	p.Params = params
	return nil
}
