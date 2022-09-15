package ethprotocol

import (
	"strings"

	"github.com/0xVanfer/ethaddr"
)

// For use of aave-like protocols:
// Aave.
func (p *Protocol) IsAToken(address string) bool {
	switch p.ProtocolName {
	// aave(aave v2)
	case ethaddr.AaveV2Protocol:
		for _, atoken := range ethaddr.AaveATokenV2List[p.Network] {
			if strings.EqualFold(address, atoken) {
				return true
			}
		}
		return false
	// aavev3
	case ethaddr.AaveV3Protocol:
		for _, atoken := range ethaddr.AaveATokenV3List[p.Network] {
			if strings.EqualFold(address, atoken) {
				return true
			}
		}
		return false
	default:
		return false
	}
}
