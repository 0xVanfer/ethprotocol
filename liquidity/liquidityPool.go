package liquidity

import (
	"github.com/0xVanfer/ethprotocol/apy"
	"github.com/0xVanfer/ethprotocol/erc"
)

// Struct for liquidity pools.
type LiquidityPool struct {
	Tokens     []TokenOfLp        // the tokens to make up the lp
	LpToken    erc.ERC20          // basic info of lp token
	ApyInfo    apy.ApyInfo        // apy info
	Reserve    float64            // tvl in amount
	ReserveUSD float64            // tvl in usd
	Volume24   float64            // trade volume in 24 hours
	OtherInfo  LiquidityOtherInfo // some other infos for special protocols
}

// Struct for tokens to make up lp token.
type TokenOfLp struct {
	Basic      erc.ERC20 // basic info of token
	Underlying erc.ERC20 // basic info of underlying token, if has no underlying, use basic
	Reserve    float64   // reserve of the single token in amount
	ReserveUSD float64   // reserve of the single token in usd
}

// Struct for other infos for special protocols.
type LiquidityOtherInfo struct {
	Liability     float64 // platypus Liability
	Cash          float64 // platypus cash
	CoverageRatio float64 // platypus coverage ratio
}
