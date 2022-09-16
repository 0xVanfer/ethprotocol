package aavelike

import (
	"github.com/0xVanfer/erc"
	"github.com/0xVanfer/ethprotocol/internal/apy"
	"github.com/0xVanfer/ethprotocol/model"
)

type AToken struct {
	ProtocolBase    model.ProtocolBase
	Basic           *erc.ERC20Info // basic info of the token
	UnderlyingBasic *erc.ERC20Info // basic info of the underlying token
	ApyInfo         *apy.ApyInfo   // deposit apy info
}

type VToken struct {
	ProtocolBase    model.ProtocolBase
	Basic           erc.ERC20Info // basic info of the token
	UnderlyingBasic erc.ERC20Info // basic info of the underlying token
	ApyInfo         apy.ApyInfo   // borrow variable apy info
}

type SToken struct {
	ProtocolBase    model.ProtocolBase
	Basic           erc.ERC20Info // basic info of the token
	UnderlyingBasic erc.ERC20Info // basic info of the underlying token
	ApyInfo         apy.ApyInfo   // borrow stable apy info
}
