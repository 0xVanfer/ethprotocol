package lend

import (
	"github.com/0xVanfer/erc"
	"github.com/0xVanfer/ethprotocol/lend/lendtoken"
	"github.com/0xVanfer/ethprotocol/model"
)

type LendPool struct {
	ProtocolBasic   *model.ProtocolBasic // Protocol basic info.
	UnderlyingBasic *erc.ERC20Info       // The underlying erc20 info.
	PoolType        LendPoolType         // Only support aave-like and compound-like.
	AToken          lendtoken.AToken     // AToken info, if the protocol is aave-like.
	VToken          lendtoken.VToken     // VToken info, if the protocol is aave-like.
	SToken          lendtoken.SToken     // SToken info, if the protocol is aave-like.
	CToken          lendtoken.CToken     // CToken info, if the protocol is compound-like.
	Params          LendPoolParams       // Pool params(almost constant).
	Status          LendPoolStatus       // Pool status.
}

type LendPoolType struct {
	IsAaveLike     bool // Whether the protocol is aave-like.
	IsCompoundLike bool // Whether the protocol is compound-like.
}

type LendPoolParams struct {
	CollateralFactor   float64 // The maximum of collateral factor.
	LiquidationLimit   float64 // Liquidation will occur when liquidation limit is reached.
	LiquidationPenalty float64 // Penalty when liquidation occurs.
	AllowBorrow        int     // 1 for true, 2 for false
	AllowCollateral    int     // 1 for true, 2 for false
}

type LendPoolStatus struct {
	TotalSupply    float64 // Total amount supplied into the pool.
	TotalBorrow    float64 // Total amount borrowed from the pool.
	SupplyLimit    float64 // The supply limit of the pool.
	SupplyCapacity float64 // SupplyCapacity = SupplyLimit - TotalSupply.
}
