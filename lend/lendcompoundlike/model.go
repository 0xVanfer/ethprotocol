package lendcompoundlike

import (
	"github.com/0xVanfer/erc"
	"github.com/0xVanfer/ethprotocol/internal/apy"
	"github.com/0xVanfer/ethprotocol/model"
)

type CToken struct {
	ProtocolBasic   model.ProtocolBasic
	Basic           erc.ERC20Info // basic info of the token
	UnderlyingBasic erc.ERC20Info // basic info of the underlying token
	DepositApyInfo  apy.ApyInfo   // deposit apy info
	BorrowApyInfo   apy.ApyInfo   // deposit apy info
}
