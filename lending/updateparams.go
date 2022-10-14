package lending

import "errors"

// Update the pool params.
//
// Not sure how to read these params from contracts except aave v3.
func (p *LendingPool) UpdatePoolParams(params LendingPoolParams) error {
	if !p.ProtocolBasic.Regularcheck() {
		return errors.New("lend pool protocol basic must be initialized")
	}
	p.Params = params
	return nil
}
