package ethprotocol

import "github.com/0xVanfer/ethaddr"

// Return protocol token.
func (p *Protocol) ProtocolToken() string {
	switch p.ProtocolName {
	case ethaddr.AaveV2Protocol, ethaddr.AaveV3Protocol:
		return ethaddr.AaveTokenList[p.Network]
	case ethaddr.AlphaProtocol:
		return ethaddr.AlphaTokenList[p.Network]
	case ethaddr.AxialProtocol:
		return ethaddr.AxialTokenList[p.Network]
	case ethaddr.BalancerProtocol:
		return ethaddr.BalancerTokenList[p.Network]
	case ethaddr.BeefyProtocol:
		return ethaddr.BeefyTokenList[p.Network]
	case ethaddr.BenqiProtocol:
		return ethaddr.BenqiTokenList[p.Network]
	case ethaddr.CompoundProtocol:
		return ethaddr.CompoundTokenList[p.Network]
	case ethaddr.ConvexProtocol:
		return ethaddr.ConvexTokenList[p.Network]
	case ethaddr.CurveProtocol:
		return ethaddr.CurveTokenlist[p.Network]
	case ethaddr.KyberProtocol:
		return ethaddr.KyberTokenList[p.Network]
	case ethaddr.PangolinProtocol:
		return ethaddr.PangolinTokenList[p.Network]
	case ethaddr.SushiProtocol:
		return ethaddr.SushiTokenList[p.Network]
	case ethaddr.TraderJoeProtocol:
		return ethaddr.TraderjoeTokenList[p.Network]
	case ethaddr.UniswapProtocolV2, ethaddr.UniswapProtocolV3:
		return ethaddr.UniswapTokenList[p.Network]
	case ethaddr.VectorProtocol:
		return ethaddr.VectorTokenList[p.Network]
	default:
		return ""
	}
}
