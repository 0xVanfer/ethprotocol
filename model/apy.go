package model

import (
	"github.com/0xVanfer/erc"
)

// Apys of a token.
//
// Use 0.01 for 1%
type ApyInfo struct {
	Base      *ApyBase     // basic rewards details
	Incentive ApyIncentive // incentive rewards(some protocols may have various incentive reward tokens)
}

type ApyBase struct {
	Apy          float64
	Apr          float64
	RewardToken  *erc.ERC20Info
	IsChainToken bool
}

type ApyIncentive struct {
	TotalApyIncentive float64 // total incentive rewards
	TotalAprIncentive float64 // total incentive rewards
	Details           []*ApyBase
}
