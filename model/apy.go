package model

import (
	"github.com/0xVanfer/erc"
)

// Apys of a token.
//
// Use 0.01 for 1%
type ApyInfo struct {
	TotalApyIncentive float64    // total incentive rewards
	TotalAprIncentive float64    // total incentive rewards
	Base              *ApyBase   // basic rewards details
	Incentive         []*ApyBase // incentive rewards(some protocols may have various incentive reward tokens)
}

type ApyBase struct {
	Apy          float64
	Apr          float64
	RewardToken  *erc.ERC20Info
	IsChainToken bool
}
