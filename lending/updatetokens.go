package lending

import (
	"errors"

	"github.com/0xVanfer/erc"
)

// Update pool tokens info by underlying token.
func (p *LendingPool) UpdateTokensByUnderlying(underlying string) error {
	if !p.ProtocolBasic.Regularcheck() {
		return errors.New("lend pool protocol basic must be initialized")
	}
	// underlying basic
	var newUnderlying erc.ERC20Info
	err := newUnderlying.Init(underlying, p.ProtocolBasic.Network, *p.ProtocolBasic.Client)
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
		p.CToken.UnderlyingBasic = &newUnderlying
		p.CToken.ProtocolBasic = p.ProtocolBasic

		_ = p.CToken.UpdateCTokenByUnderlying(underlying)
	}
	return nil
}

// Update pool tokens info by a token.
func (p *LendingPool) UpdateTokensByAToken(atoken string) error {
	p.AToken.ProtocolBasic = p.ProtocolBasic
	underlyingAddress, err := p.AToken.GetUnderlyingAddress(atoken)
	if err != nil {
		return err
	}
	return p.UpdateTokensByUnderlying(underlyingAddress)
}

// Update pool tokens info by v token.
func (p *LendingPool) UpdateTokensByVToken(vtoken string) error {
	p.VToken.ProtocolBasic = p.ProtocolBasic
	underlyingAddress, err := p.VToken.GetUnderlyingAddress(vtoken)
	if err != nil {
		return err
	}
	return p.UpdateTokensByUnderlying(underlyingAddress)
}

// Update pool tokens info by s token.
func (p *LendingPool) UpdateTokensBySToken(stoken string) error {
	p.SToken.ProtocolBasic = p.ProtocolBasic
	underlyingAddress, err := p.SToken.GetUnderlyingAddress(stoken)
	if err != nil {
		return err
	}
	return p.UpdateTokensByUnderlying(underlyingAddress)
}

// Update pool tokens info by c token.
func (p *LendingPool) UpdateTokensByCToken(ctoken string) error {
	p.CToken.ProtocolBasic = p.ProtocolBasic
	underlyingAddress, err := p.CToken.GetUnderlyingAddress(ctoken)
	if err != nil {
		return err
	}
	return p.UpdateTokensByUnderlying(underlyingAddress)
}
