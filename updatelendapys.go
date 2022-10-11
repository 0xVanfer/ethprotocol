package ethprotocol

import (
	"errors"

	"github.com/0xVanfer/ethaddr"
)

// Update all the protocol's lend pools apys.
func (prot *Protocol) UpdateLendApys() error {
	// protocol basic must be initialized
	if !prot.ProtocolBasic.Regularcheck() {
		return errors.New("protocol basic must be initialized")
	}
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

// Update some of the protocol's lend pools apys by given underlying addresses.
//
// If "underlyings" is empty, equal to UpdateLendApys().
func (prot *Protocol) UpdateLendApyByUnderlying(underlyings []string) error {
	// protocol basic must be initialized
	if !prot.ProtocolBasic.Regularcheck() {
		return errors.New("protocol basic must be initialized")
	}
	switch prot.ProtocolBasic.ProtocolName {
	// aave v2
	case ethaddr.AaveV2Protocol:
		return prot.updateAaveV2Lend(underlyings)
	// aave v3
	case ethaddr.AaveV3Protocol:
		return prot.updateAaveV3Lend(underlyings)
	// benqi
	case ethaddr.BenqiProtocol:
		return prot.updateBenqiLend(underlyings)
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
