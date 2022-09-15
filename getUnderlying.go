package ethprotocol

import (
	"errors"
	"strings"

	"github.com/0xVanfer/ethaddr"
)

func (p *Protocol) GetUnderlying(address string) (string, error) {
	switch p.ProtocolName {
	// aave(aave v2)
	case ethaddr.AaveV2Protocol:
		for underlying, atoken := range ethaddr.AaveATokenV2List[p.Network] {
			if strings.EqualFold(address, atoken) {
				return underlying, nil
			}
		}
		return "", errors.New(p.ProtocolName + " a token not found")
	// aavev3
	case ethaddr.AaveV3Protocol:
		for underlying, atoken := range ethaddr.AaveATokenV3List[p.Network] {
			if strings.EqualFold(address, atoken) {
				return underlying, nil
			}
		}
		return "", errors.New(p.ProtocolName + " a token not found")
	// compound
	case ethaddr.CompoundProtocol:
		for underlying, ctoken := range ethaddr.CompoundCTokenList[p.Network] {
			if strings.EqualFold(address, ctoken) {
				return underlying, nil
			}
		}
		return "", errors.New(p.ProtocolName + " c token not found")
	// benqi
	case ethaddr.BenqiProtocol:
		for underlying, ctoken := range ethaddr.BenqiCTokenList[p.Network] {
			if strings.EqualFold(address, ctoken) {
				return underlying, nil
			}
		}
		return "", errors.New(p.ProtocolName + " c token not found")
	// traderjoe
	case ethaddr.TraderJoeProtocol:
		for underlying, ctoken := range ethaddr.TraderjoeCTokenList[p.Network] {
			if strings.EqualFold(address, ctoken) {
				return underlying, nil
			}
		}
		return "", errors.New(p.ProtocolName + " c token not found")
	default:
		return "", errors.New("protocol not supported")
	}
}
