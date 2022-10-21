package liquidity

import (
	"errors"

	"github.com/0xVanfer/ethprotocol/model"
)

func (p *LiquidityPool) Init(protocolBasic model.ProtocolBasic) error {
	// regular check
	if !protocolBasic.Regularcheck() {
		return errors.New("protocol basic must not be initialized")
	}

	p.ProtocolBasic = &protocolBasic
	return nil
}
