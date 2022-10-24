package lending

import (
	"github.com/0xVanfer/erc"
	"github.com/0xVanfer/ethprotocol/lending/lendingtoken"
	"github.com/0xVanfer/ethprotocol/model"
	"github.com/shopspring/decimal"
)

type LendingPool struct {
	ProtocolBasic     *model.ProtocolBasic // Protocol basic info.
	UnderlyingBasic   *erc.ERC20Info       // The underlying erc20 info.
	PoolType          LendingPoolType      // Only support aave-like and compound-like.
	SupportedNetworks []string             // Protocol can run on which networks.
	AToken            lendingtoken.AToken  // AToken info, if the protocol is aave-like.
	VToken            lendingtoken.VToken  // VToken info, if the protocol is aave-like.
	SToken            lendingtoken.SToken  // SToken info, if the protocol is aave-like.
	CToken            lendingtoken.CToken  // CToken info, if the protocol is compound-like.
	Status            LendingPoolStatus    // Pool status.
}

type LendingPoolType struct {
	IsAaveLike     bool // Whether the protocol is aave-like.
	IsCompoundLike bool // Whether the protocol is compound-like.
}

type LendingPoolStatus struct {
	CollateralFactor   decimal.Decimal // The maximum of collateral factor.
	LiquidationLimit   decimal.Decimal // Liquidation will occur when liquidation limit is reached.
	LiquidationPenalty decimal.Decimal // Penalty when liquidation occurs.
	AllowBorrow        bool            // If the token can be borrowed in this protocol.
	AllowCollateral    bool            // If the token can be used as collateral in this protocol.

	BorrowLimit     decimal.Decimal // The borrow limit of the pool.
	SupplyLimit     decimal.Decimal // The supply limit of the pool.
	SupplyCapacity  decimal.Decimal // SupplyCapacity = SupplyLimit - TotalSupply.
	TotalSupply     decimal.Decimal // Total amount supplied into the pool.
	TotalVBorrow    decimal.Decimal // Total variable amount borrowed from the pool.
	TotalSBorrow    decimal.Decimal // Total stable amount borrowed from the pool.
	TotalCBorrow    decimal.Decimal // Total (compound like) amount borrowed from the pool.
	UtilizationRate decimal.Decimal // = TotalBorrow / TotalSupply

	EModeCategoryId       int             // Aave v3 emode category id.
	EModeCollateralFactor decimal.Decimal // Aave v3 emode collateral factor.
	EModeLiquidationLimit decimal.Decimal // Aave v3 emode liquidation limit.

	BorrowableInIsolation bool // Aave v3 isolation mode can be borrowed.
}
