package ethprotocol

import (
	"github.com/0xVanfer/coingecko"
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/lending"
	"github.com/0xVanfer/ethprotocol/liquidity"
	"github.com/0xVanfer/ethprotocol/model"
	"github.com/0xVanfer/ethprotocol/stake"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type Protocol struct {
	ProtocolBasic  *model.ProtocolBasic
	LiquidityPools []*liquidity.LiquidityPool
	StakePools     []*stake.StakePool
	LendingPools   []*lending.LendingPool
}

type ProtocolInput struct {
	Network   string               // Network of the protocol.
	Name      string               // Name of the protocol, given by github.com/0xVanfer/ethaddr.
	Client    bind.ContractBackend // Contract backend to call the contracts. Can be nil, but most functions will not work properly.
	Coingecko coingecko.Gecko      // Coingecko. Can also input the key only.
}

// Return protocol token according to network.
func (p *Protocol) Token() string {
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
