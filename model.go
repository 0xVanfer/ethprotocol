package ethprotocol

import (
	"github.com/0xVanfer/ethprotocol/lend"
	"github.com/0xVanfer/ethprotocol/liquidity"
	"github.com/0xVanfer/ethprotocol/model"
	"github.com/0xVanfer/ethprotocol/stake"
)

type Protocol struct {
	ProtocolBasic  *model.ProtocolBasic
	LiquidityPools []*liquidity.LiquidityPool
	StakePools     []*stake.StakePool
	LendPools      []*lend.LendPool
}
