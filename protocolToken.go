package ethprotocol

import (
	"errors"

	"github.com/0xVanfer/ethaddr"
)

// Return protocol token.
func (p *Protocol) ProtocolToken() (string, error) {
	if !p.ProtocolBasic.Regularcheck() {
		return "", errors.New("protocol basic must be initialized")
	}
	switch p.ProtocolBasic.ProtocolName {
	case ethaddr.AaveV2Protocol, ethaddr.AaveV3Protocol:
		return ethaddr.AaveTokenList[p.ProtocolBasic.Network], nil
	case ethaddr.AlphaProtocol:
		return ethaddr.AlphaTokenList[p.ProtocolBasic.Network], nil
	case ethaddr.AxialProtocol:
		return ethaddr.AxialTokenList[p.ProtocolBasic.Network], nil
	case ethaddr.BalancerProtocol:
		return ethaddr.BalancerTokenList[p.ProtocolBasic.Network], nil
	case ethaddr.BeefyProtocol:
		return ethaddr.BeefyTokenList[p.ProtocolBasic.Network], nil
	case ethaddr.BenqiProtocol:
		return ethaddr.BenqiTokenList[p.ProtocolBasic.Network], nil
	case ethaddr.CompoundProtocol:
		return ethaddr.CompoundTokenList[p.ProtocolBasic.Network], nil
	case ethaddr.ConvexProtocol:
		return ethaddr.ConvexTokenList[p.ProtocolBasic.Network], nil
	case ethaddr.CurveProtocol:
		return ethaddr.CurveTokenlist[p.ProtocolBasic.Network], nil
	case ethaddr.KyberProtocol:
		return ethaddr.KyberTokenList[p.ProtocolBasic.Network], nil
	case ethaddr.PangolinProtocol:
		return ethaddr.PangolinTokenList[p.ProtocolBasic.Network], nil
	case ethaddr.SushiProtocol:
		return ethaddr.SushiTokenList[p.ProtocolBasic.Network], nil
	case ethaddr.TraderJoeProtocol:
		return ethaddr.TraderjoeTokenList[p.ProtocolBasic.Network], nil
	case ethaddr.UniswapProtocolV2, ethaddr.UniswapProtocolV3:
		return ethaddr.UniswapTokenList[p.ProtocolBasic.Network], nil
	case ethaddr.VectorProtocol:
		return ethaddr.VectorTokenList[p.ProtocolBasic.Network], nil
	default:
		return "", errors.New("protocol not supported")
	}
}
