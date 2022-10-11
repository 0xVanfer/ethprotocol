package lendaavelike

import (
	"github.com/0xVanfer/erc"
	"github.com/0xVanfer/ethprotocol/model"
)

type AToken struct {
	ProtocolBasic   *model.ProtocolBasic
	Basic           *erc.ERC20Info // basic info of the token
	UnderlyingBasic *erc.ERC20Info // basic info of the underlying token
	ApyInfo         model.ApyInfo  // deposit apy info
}

type VToken struct {
	ProtocolBasic   *model.ProtocolBasic
	Basic           *erc.ERC20Info // basic info of the token
	UnderlyingBasic *erc.ERC20Info // basic info of the underlying token
	ApyInfo         model.ApyInfo  // borrow variable apy info
}

type SToken struct {
	ProtocolBasic   *model.ProtocolBasic
	Basic           *erc.ERC20Info // basic info of the token
	UnderlyingBasic *erc.ERC20Info // basic info of the underlying token
	ApyInfo         model.ApyInfo  // borrow stable apy info
}
