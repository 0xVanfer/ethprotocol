package lend

import "errors"

// Update pool tokens info by underlying token.
func (p *LendPool) UpdateTokensByUnderlying(underlying string) error {
	if !p.ProtocolBasic.Regularcheck() {
		return errors.New("lend pool protocol basic must be initialized")
	}
	// underlying basic
	err := p.UnderlyingBasic.Init(underlying, p.ProtocolBasic.Network, *p.ProtocolBasic.Client)
	if err != nil {
		return err
	}
	// avsc basic
	if p.PoolType.IsAaveLike {
		p.AToken.UnderlyingBasic = p.UnderlyingBasic
		p.VToken.UnderlyingBasic = p.UnderlyingBasic
		p.SToken.UnderlyingBasic = p.UnderlyingBasic

		p.AToken.ProtocolBasic = p.ProtocolBasic
		p.VToken.ProtocolBasic = p.ProtocolBasic
		p.SToken.ProtocolBasic = p.ProtocolBasic

		err = p.AToken.UpdateATokenByUnderlying(underlying)
		if err != nil {
			return err
		}
		err = p.VToken.UpdateVTokenByUnderlying(underlying)
		if err != nil {
			return err
		}
		err = p.AToken.UpdateATokenByUnderlying(underlying)
		if err != nil {
			return err
		}
	}
	if p.PoolType.IsCompoundLike {
		p.CToken.ProtocolBasic = p.ProtocolBasic
		p.CToken.UnderlyingBasic = p.UnderlyingBasic

		err = p.CToken.UpdateCTokenByUnderlying(underlying)
		if err != nil {
			return err
		}
	}
	return nil
}

// Update pool tokens info by a token.
func (p *LendPool) UpdateTokensByAToken(atoken string) error {
	underlyingAddress, err := p.AToken.GetUnderlyingAddress(atoken)
	if err != nil {
		return err
	}
	err = p.UpdateTokensByUnderlying(underlyingAddress)
	if err != nil {
		return err
	}
	return nil
}

// Update pool tokens info by v token.
func (p *LendPool) UpdateTokensByVToken(vtoken string) error {
	underlyingAddress, err := p.VToken.GetUnderlyingAddress(vtoken)
	if err != nil {
		return err
	}
	err = p.UpdateTokensByUnderlying(underlyingAddress)
	if err != nil {
		return err
	}
	return nil
}

// Update pool tokens info by s token.
func (p *LendPool) UpdateTokensBySToken(stoken string) error {
	underlyingAddress, err := p.SToken.GetUnderlyingAddress(stoken)
	if err != nil {
		return err
	}
	err = p.UpdateTokensByUnderlying(underlyingAddress)
	if err != nil {
		return err
	}
	return nil
}

// Update pool tokens info by c token.
func (p *LendPool) UpdateTokensByCToken(ctoken string) error {
	underlyingAddress, err := p.CToken.GetUnderlyingAddress(ctoken)
	if err != nil {
		return err
	}
	err = p.UpdateTokensByUnderlying(underlyingAddress)
	if err != nil {
		return err
	}
	return nil
}
