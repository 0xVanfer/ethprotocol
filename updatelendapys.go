package ethprotocol

import (
	"errors"

	"github.com/0xVanfer/ethaddr"
)

func (prot *Protocol) UpdateLendApys() error {
	switch prot.ProtocolBasic.ProtocolName {
	// aave v2
	case ethaddr.AaveV2Protocol:
		return prot.updateAaveV2Lend([]string{})
	// aave v3
	case ethaddr.AaveV3Protocol:
		return prot.updateAaveV3Lend([]string{})
	// benqi
	case ethaddr.BenqiProtocol:
		return prot.updateBenqiLend([]string{})
	case ethaddr.CompoundProtocol:
		// todo
		return nil
	case ethaddr.TraderJoeProtocol:
		// todo
		return nil
	default:
		return errors.New("protocol not supported lend pools")
	}
}
