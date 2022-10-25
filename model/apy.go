package model

import (
	"github.com/0xVanfer/erc"
	"github.com/0xVanfer/types"
	"github.com/0xVanfer/utils"
	"github.com/shopspring/decimal"
)

// Apys of a token.
//
// Use 0.01 for 1%
type ApyInfo struct {
	Apy         decimal.Decimal
	Apr         decimal.Decimal
	RewardToken *erc.ERC20Info

	IncentiveTotalApy  decimal.Decimal
	IncentiveTotalApr  decimal.Decimal
	IncentiveToken0Apy decimal.Decimal
	IncentiveToken0Apr decimal.Decimal
	IncentiveToken1Apy decimal.Decimal
	IncentiveToken1Apr decimal.Decimal
	IncentiveToken0    *erc.ERC20Info
	IncentiveToken1    *erc.ERC20Info
}

func (i *ApyInfo) Generate() {
	// base apy, apr
	if i.Apr.IsZero() {
		i.Apr = toApr(i.Apy)
	}
	if i.Apy.IsZero() {
		i.Apy = toApy(i.Apr)
	}

	// token 0
	if i.IncentiveToken0Apr.IsZero() {
		i.IncentiveToken0Apr = toApr(i.IncentiveToken0Apy)
	}
	if i.IncentiveToken0Apy.IsZero() {
		i.IncentiveToken0Apy = toApy(i.IncentiveToken0Apr)
	}

	// token 1
	if i.IncentiveToken1Apr.IsZero() {
		i.IncentiveToken1Apr = toApr(i.IncentiveToken1Apy)
	}
	if i.IncentiveToken1Apy.IsZero() {
		i.IncentiveToken1Apy = toApy(i.IncentiveToken1Apr)
	}

	// incentive total
	if i.IncentiveTotalApr.IsZero() {
		if !i.IncentiveTotalApy.IsZero() {
			i.IncentiveTotalApr = toApr(i.IncentiveTotalApy)
		} else {
			i.IncentiveTotalApr = i.IncentiveToken0Apr.Add(i.IncentiveToken1Apr)
		}
	}
	if i.IncentiveTotalApy.IsZero() {
		if !i.IncentiveTotalApr.IsZero() {
			i.IncentiveTotalApy = toApy(i.IncentiveTotalApr)
		} else {
			i.IncentiveTotalApy = i.IncentiveToken0Apy.Add(i.IncentiveToken1Apy)
		}
	}
}

// Apy into apr.
func toApr(apy decimal.Decimal) decimal.Decimal {
	if apy.IsZero() {
		return decimal.Zero
	}
	return types.ToDecimal(utils.Apy2Apr(apy))
}

// Apr into apy.
func toApy(apr decimal.Decimal) decimal.Decimal {
	if apr.IsZero() {
		return decimal.Zero
	}
	return types.ToDecimal(utils.Apr2Apy(apr))
}
