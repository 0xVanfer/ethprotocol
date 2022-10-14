package lending

import (
	"github.com/0xVanfer/erc"
	"github.com/0xVanfer/ethprotocol/lending/lendingtoken"
	"github.com/0xVanfer/ethprotocol/model"
)

type LendingPool struct {
	ProtocolBasic   *model.ProtocolBasic // Protocol basic info.
	UnderlyingBasic *erc.ERC20Info       // The underlying erc20 info.
	PoolType        LendingPoolType      // Only support aave-like and compound-like.
	AToken          lendingtoken.AToken  // AToken info, if the protocol is aave-like.
	VToken          lendingtoken.VToken  // VToken info, if the protocol is aave-like.
	SToken          lendingtoken.SToken  // SToken info, if the protocol is aave-like.
	CToken          lendingtoken.CToken  // CToken info, if the protocol is compound-like.
	Params          LendingPoolParams    // Pool params(almost constant).
	Status          LendingPoolStatus    // Pool status.
}

type LendingPoolType struct {
	IsAaveLike     bool // Whether the protocol is aave-like.
	IsCompoundLike bool // Whether the protocol is compound-like.
}

type LendingPoolParams struct {
	CollateralFactor   float64 // The maximum of collateral factor.
	LiquidationLimit   float64 // Liquidation will occur when liquidation limit is reached.
	LiquidationPenalty float64 // Penalty when liquidation occurs.
	AllowBorrow        int     // 1 for true, 2 for false
	AllowCollateral    int     // 1 for true, 2 for false
}

type LendingPoolStatus struct {
	BorrowLimit     float64 // The borrow limit of the pool.
	SupplyLimit     float64 // The supply limit of the pool.
	SupplyCapacity  float64 // SupplyCapacity = SupplyLimit - TotalSupply.
	TotalSupply     float64 // Total amount supplied into the pool.
	TotalBorrow     float64 // Total amount borrowed from the pool.
	UtilizationRate float64 // = TotalBorrow / TotalSupply
}