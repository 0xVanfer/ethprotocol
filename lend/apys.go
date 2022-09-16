package lend

import (
	"errors"

	"github.com/0xVanfer/coingecko"
	"github.com/0xVanfer/ethaddr"
)

// BETTER NOT use this for a list of lend pools.
func (p *LendPool) GetApys(gecko *coingecko.Gecko) error {
	if p.IsAaveLike {
		switch p.ProtocolName {
		// aave v2
		case ethaddr.AaveV2Protocol:
			// todo
		// aave v3
		case ethaddr.AaveV3Protocol:
			// todo
		}
		return nil
	}
	if p.IsCompoundLike {
		switch p.ProtocolName {
		// compound
		case ethaddr.CompoundProtocol:
			// todo
		// benqi
		case ethaddr.BenqiProtocol:
			// todo
		// traderjoe
		case ethaddr.TraderJoeProtocol:
			// todo
		}
		return nil
	}
	return errors.New("must be either aave-like or compound-like")
}
