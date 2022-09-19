package lend

import (
	"errors"

	"github.com/0xVanfer/chainId"
	"github.com/0xVanfer/erc"
	"github.com/0xVanfer/ethprotocol/test/eth"
)

// Update pool tokens info by underlying token.
func (p *LendPool) UpdateTokensByUnderlying(underlying string) error {
	if !p.ProtocolBasic.Regularcheck() {
		return errors.New("lend pool protocol basic must be initialized")
	}
	// underlying basic
	var newUnderlying erc.ERC20Info
	err := newUnderlying.Init(underlying, p.ProtocolBasic.Network, eth.GetConnector(chainId.AvalancheChainName))
	if err != nil {
		return err
	}
	p.UnderlyingBasic = &newUnderlying
	// avsc basic
	if p.PoolType.IsAaveLike {
		p.AToken.UnderlyingBasic = &newUnderlying
		p.VToken.UnderlyingBasic = &newUnderlying
		p.SToken.UnderlyingBasic = &newUnderlying

		p.AToken.ProtocolBasic = p.ProtocolBasic
		p.VToken.ProtocolBasic = p.ProtocolBasic
		p.SToken.ProtocolBasic = p.ProtocolBasic

		_ = p.AToken.UpdateATokenByUnderlying(underlying)
		_ = p.VToken.UpdateVTokenByUnderlying(underlying)
		_ = p.SToken.UpdateSTokenByUnderlying(underlying)
	}
	if p.PoolType.IsCompoundLike {
		p.CToken.ProtocolBasic = p.ProtocolBasic
		p.CToken.UnderlyingBasic = p.UnderlyingBasic

		_ = p.CToken.UpdateCTokenByUnderlying(underlying)
	}
	return nil
}

// Update pool tokens info by a token.
func (p *LendPool) UpdateTokensByAToken(atoken string) error {
	underlyingAddress, err := p.AToken.GetUnderlyingAddress(atoken)
	if err != nil {
		return err
	}
	return p.UpdateTokensByUnderlying(underlyingAddress)
}

// Update pool tokens info by v token.
func (p *LendPool) UpdateTokensByVToken(vtoken string) error {
	underlyingAddress, err := p.VToken.GetUnderlyingAddress(vtoken)
	if err != nil {
		return err
	}
	return p.UpdateTokensByUnderlying(underlyingAddress)
}

// Update pool tokens info by s token.
func (p *LendPool) UpdateTokensBySToken(stoken string) error {
	underlyingAddress, err := p.SToken.GetUnderlyingAddress(stoken)
	if err != nil {
		return err
	}
	return p.UpdateTokensByUnderlying(underlyingAddress)
}

// Update pool tokens info by c token.
func (p *LendPool) UpdateTokensByCToken(ctoken string) error {
	underlyingAddress, err := p.CToken.GetUnderlyingAddress(ctoken)
	if err != nil {
		return err
	}
	return p.UpdateTokensByUnderlying(underlyingAddress)
}
