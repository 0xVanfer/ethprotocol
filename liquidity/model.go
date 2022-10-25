package liquidity

import (
	"github.com/0xVanfer/erc"
	"github.com/0xVanfer/ethprotocol/model"
	"github.com/shopspring/decimal"
)

// Struct for liquidity pools.
type LiquidityPool struct {
	ProtocolBasic *model.ProtocolBasic // protocol basic
	PoolName      string               // pool name
	LpToken       *erc.ERC20Info       // basic info of lp token
	Tokens        []*TokenOfLp         // the tokens to make up the lp
	ApyInfo       *model.ApyInfo       // apy info
	Reserve       decimal.Decimal      // tvl in amount
	ReserveUSD    decimal.Decimal      // tvl in usd
	Volume24      decimal.Decimal      // trade volume in 24 hours
	OtherInfo     *LiquidityOtherInfo  // some other infos for special protocols
}

// Struct for tokens to make up lp token.
type TokenOfLp struct {
	Basic      *erc.ERC20Info  // basic info of token
	Underlying *erc.ERC20Info  // basic info of underlying token, if has no underlying, use basic
	Reserve    decimal.Decimal // reserve of the single token in amount
	ReserveUSD decimal.Decimal // reserve of the single token in usd
}

// Struct for other infos for special protocols.
type LiquidityOtherInfo struct {
	Liability     decimal.Decimal // platypus Liability
	Cash          decimal.Decimal // platypus cash
	CoverageRatio decimal.Decimal // platypus coverage ratio
}
