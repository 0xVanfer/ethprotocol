package ethprotocol

import (
	"strings"

	"github.com/0xVanfer/ethaddr"
)

// For use of aave-like protocols:
// Aave.
func (p *Protocol) IsVToken(address string) bool {
	switch p.ProtocolName {
	// aave(aave v2)
	case ethaddr.AaveV2Protocol:
		for _, vtoken := range ethaddr.AaveVTokenV2List[p.Network] {
			if strings.EqualFold(address, vtoken) {
				return true
			}
		}
		return false
	// aavev3
	case ethaddr.AaveV3Protocol:
		for _, vtoken := range ethaddr.AaveVTokenV3List[p.Network] {
			if strings.EqualFold(address, vtoken) {
				return true
			}
		}
		return false
	default:
		return false
	}
}
