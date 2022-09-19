package lend

import (
	"errors"

	"github.com/0xVanfer/erc"
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/lend/lendaavelike"
	"github.com/0xVanfer/ethprotocol/lend/lendcompoundlike"
	"github.com/0xVanfer/ethprotocol/model"
)

type LendPool struct {
	ProtocolBasic   *model.ProtocolBasic
	UnderlyingBasic *erc.ERC20Info
	PoolType        LendPoolType
	AToken          lendaavelike.AToken
	VToken          lendaavelike.VToken
	SToken          lendaavelike.SToken
	CToken          lendcompoundlike.CToken
	Params          LendPoolParams
}
type LendPoolType struct {
	IsAaveLike     bool
	IsCompoundLike bool
}

// Initialize the lend pool protocol basic and pool type.
func (p *LendPool) Init(protocolBasic model.ProtocolBasic) error {
	switch protocolBasic.ProtocolName {
	case ethaddr.AaveV2Protocol, ethaddr.AaveV3Protocol:
		p.PoolType = LendPoolType{IsAaveLike: true}
	case ethaddr.BenqiProtocol, ethaddr.CompoundProtocol, ethaddr.TraderJoeProtocol:
		p.PoolType = LendPoolType{IsCompoundLike: true}
	default:
		return errors.New("protocol not supported lend pool")
	}
	if protocolBasic.Network == "" {
		return errors.New("network must not be empty")
	}
	p.ProtocolBasic = &protocolBasic
	return nil
}
