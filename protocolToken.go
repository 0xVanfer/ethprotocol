package ethprotocol

import (
	"github.com/0xVanfer/ethaddr"
)

// Return protocol token according to network.
func (p *Protocol) GetProtocolToken() string {
	if !p.ProtocolBasic.Regularcheck() {
		return ""
	}
	// not supported protocol name
	tokenList := ethaddr.ProtocolTokenListMap[p.ProtocolBasic.ProtocolName]
	if tokenList == nil {
		return ""
	}
	return tokenList[p.ProtocolBasic.Network]
}
