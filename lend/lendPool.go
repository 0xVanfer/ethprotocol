package lend

import (
	"errors"

	"github.com/0xVanfer/erc"
	"github.com/0xVanfer/ethprotocol/lend/lendatokens"
	"github.com/0xVanfer/ethprotocol/lend/lendctokens"
	"github.com/0xVanfer/ethprotocol/lend/lendstokens"
	"github.com/0xVanfer/ethprotocol/lend/lendvtokens"
	"github.com/0xVanfer/ethprotocol/model"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type LendPool struct {
	ProtocolBase    model.ProtocolBase
	UnderlyingBasic erc.ERC20Info
	AToken          lendatokens.AToken
	VToken          lendvtokens.VToken
	SToken          lendstokens.SToken
	CToken          lendctokens.CToken
	Params          LendPoolParams
	PoolType        LendPoolType
}
type LendPoolType struct {
	IsAaveLike     bool
	IsCompoundLike bool
}

func (p *LendPool) Init(network string, protocolName string, client bind.ContractBackend) error {
	p.ProtocolBase.ProtocolName = protocolName
	p.ProtocolBase.Network = network
	p.ProtocolBase.Client = client
	return nil
}

func (p *LendPool) UpdateTokensByAToken() error {
	return nil
}

// Initialize lend pool tokens.
func (p *LendPool) InitTokens(underlying string, network string, protocolName string, client bind.ContractBackend) error {

	err := p.UnderlyingBasic.Init(underlying, network, client)
	if err != nil {
		return err
	}

	// if has c token， which means it is compound-like
	foundC, err := p.CToken.InitWithUnderlying(underlying, network, protocolName, client)
	if err != nil {
		return err
	}
	// if is not compound-like, try aave-like
	if !foundC {
		foundA, err := p.AToken.InitWithUnderlying(underlying, network, protocolName, client)
		if err != nil {
			return err
		}
		foundV, err := p.VToken.InitWithUnderlying(underlying, network, protocolName, client)
		if err != nil {
			return err
		}
		foundS, err := p.SToken.InitWithUnderlying(underlying, network, protocolName, client)
		if err != nil {
			return err
		}
		// both not supported by compound-like and aave-like
		if !(foundA || foundS || foundV) {
			return errors.New("underlying token not supported")
		} else {
			p.PoolType.IsAaveLike = true
			return nil
		}
	} else {
		p.PoolType.IsCompoundLike = true
		return nil
	}
}
