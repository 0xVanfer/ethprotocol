package model

import (
	"github.com/0xVanfer/erc"
	"github.com/shopspring/decimal"
)

// Apys of a token.
//
// Use 0.01 for 1%
type ApyInfo struct {
	Base      *ApyBase     // basic rewards details
	Incentive ApyIncentive // incentive rewards(some protocols may have various incentive reward tokens)
}

type ApyBase struct {
	Apy          decimal.Decimal
	Apr          decimal.Decimal
	RewardToken  *erc.ERC20Info
	IsChainToken bool
}

type ApyIncentive struct {
	TotalApyIncentive decimal.Decimal // total incentive rewards
	TotalAprIncentive decimal.Decimal // total incentive rewards
	Details           []*ApyBase
}
