package lend

import "github.com/0xVanfer/ethprotocol/apy"

type LendPool struct {
	AToken             AVSCToken
	VToken             AVSCToken
	SToken             AVSCToken
	CToken             AVSCToken
	DepositApy         apy.ApyInfo
	BorrowApy          apy.ApyInfo
	CollateralFactor   float64
	LiquidationLimit   float64
	LiquidationPenalty float64
	AllowBorrow        int
	AllowCollateral    int
}
