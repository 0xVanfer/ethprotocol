package lend

import (
	"errors"

	"github.com/0xVanfer/ethprotocol/erc"
	"github.com/0xVanfer/ethprotocol/lend/lendatokens"
	"github.com/0xVanfer/ethprotocol/lend/lendctokens"
	"github.com/0xVanfer/ethprotocol/lend/lendstokens"
	"github.com/0xVanfer/ethprotocol/lend/lendvtokens"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type LendPool struct {
	UnderlyingBasic erc.ERC20
	IsAaveLike      bool
	IsCompoundLike  bool
	AToken          lendatokens.AToken
	VToken          lendvtokens.VToken
	SToken          lendstokens.SToken
	CToken          lendctokens.CToken
	Params          LendPoolParams
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
			p.IsAaveLike = true
			return nil
		}
	} else {
		p.IsCompoundLike = true
		return nil
	}
}
