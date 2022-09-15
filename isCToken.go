package ethprotocol

import (
	"strings"

	"github.com/0xVanfer/ethaddr"
)

// For use of compound-like protocols:
// Compound, benqi, traderjoe.
func (p *Protocol) IsCToken(address string) bool {
	switch p.ProtocolName {
	// compound
	case ethaddr.CompoundProtocol:
		for _, ctoken := range ethaddr.CompoundCTokenList[p.Network] {
			if strings.EqualFold(address, ctoken) {
				return true
			}
		}
		return false
	// benqi
	case ethaddr.AaveV3Protocol:
		for _, ctoken := range ethaddr.BenqiCTokenList[p.Network] {
			if strings.EqualFold(address, ctoken) {
				return true
			}
		}
		return false
	// traderjoe
	case ethaddr.TraderJoeProtocol:
		for _, ctoken := range ethaddr.TraderjoeCTokenList[p.Network] {
			if strings.EqualFold(address, ctoken) {
				return true
			}
		}
		return false
	default:
		return false
	}
}
