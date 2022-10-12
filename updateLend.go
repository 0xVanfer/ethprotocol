package ethprotocol

import (
	"errors"
	"strings"

	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/lend"
)

// Update some of the protocol's lend pools apys by given underlying addresses.
//
// If "underlyings" is empty, update all pools.
func (prot *Protocol) UpdateLendApys(underlyings ...string) error {
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
	// compound
	case ethaddr.CompoundProtocol:
		// todo
		return nil
	// tradejoe
	case ethaddr.TraderJoeProtocol:
		return prot.updateTraderjoeLend(underlyings)
	default:
		return errors.New("protocol not supported lend pools")
	}
}

// Update lend pools' params.
//
// params: map[underlying address] = pool params.
func (prot *Protocol) UpdateLendParams(params map[string]lend.LendPoolParams) error {
	// protocol basic must be initialized
	if !prot.ProtocolBasic.Regularcheck() {
		return errors.New("protocol basic must be initialized")
	}
	for _, pool := range prot.LendPools {
		for underlying, param := range params {
			if strings.EqualFold(*pool.UnderlyingBasic.Address, underlying) {
				pool.UpdatePoolParams(param)
			}
		}
	}
	return nil
}
