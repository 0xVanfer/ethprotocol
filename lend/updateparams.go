package lend

import "errors"

// Update the pool params.
//
// Not sure how to read these params from contracts except aave v3.
func (p *LendPool) UpdatePoolParams(params LendPoolParams) error {
	if !p.ProtocolBasic.Regularcheck() {
		return errors.New("lend pool protocol basic must be initialized")
	}
	p.Params = params
	return nil
}
