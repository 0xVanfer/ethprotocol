package lend

import (
	"github.com/0xVanfer/erc"
	"github.com/0xVanfer/ethprotocol/lend/lendtoken"
	"github.com/0xVanfer/ethprotocol/model"
)

type LendPool struct {
	ProtocolBasic   *model.ProtocolBasic
	UnderlyingBasic *erc.ERC20Info
	PoolType        LendPoolType
	AToken          lendtoken.AToken
	VToken          lendtoken.VToken
	SToken          lendtoken.SToken
	CToken          lendtoken.CToken
	Params          LendPoolParams
}

type LendPoolType struct {
	IsAaveLike     bool
	IsCompoundLike bool
}

type LendPoolParams struct {
	CollateralFactor   float64
	LiquidationLimit   float64
	LiquidationPenalty float64
	AllowBorrow        int
	AllowCollateral    int
}
