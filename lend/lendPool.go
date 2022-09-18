package lend

import (
	"errors"

	"github.com/0xVanfer/erc"
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/lend/lendaavelike"
	"github.com/0xVanfer/ethprotocol/lend/lendcompoundlike"
	"github.com/0xVanfer/ethprotocol/model"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type LendPool struct {
	ProtocolBasic   *model.ProtocolBasic
	UnderlyingBasic *erc.ERC20Info
	PoolType        *LendPoolType
	AToken          *lendaavelike.AToken
	VToken          *lendaavelike.VToken
	SToken          *lendaavelike.SToken
	CToken          *lendcompoundlike.CToken
	Params          *LendPoolParams
}
type LendPoolType struct {
	IsAaveLike     bool
	IsCompoundLike bool
}

func (p *LendPool) Init(network string, protocolName string, client *bind.ContractBackend) error {
	switch protocolName {
	case ethaddr.AaveV2Protocol, ethaddr.AaveV3Protocol:
		p.PoolType.IsAaveLike = true
	case ethaddr.BenqiProtocol, ethaddr.CompoundProtocol, ethaddr.TraderJoeProtocol:
		p.PoolType.IsCompoundLike = true
	default:
		return errors.New("protocol not supported")
	}
	if network == "" {
		return errors.New("network must not be empty")
	}
	p.ProtocolBasic.ProtocolName = protocolName
	p.ProtocolBasic.Network = network
	p.ProtocolBasic.Client = client
	return nil
}
