package ethprotocol

import (
	"math"

	"github.com/0xVanfer/abigen/traderjoe/traderjoeCToken"
	"github.com/0xVanfer/abigen/traderjoe/traderjoeComptroller"
	"github.com/0xVanfer/abigen/traderjoe/traderjoeRewardDistributor"
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/internal/constants"
	"github.com/0xVanfer/ethprotocol/lend"
	"github.com/0xVanfer/types"
	"github.com/0xVanfer/utils"
)

// Internal use only! No protocol regular check!
//
// Update compound-like lend pools by underlyings.
func (prot *Protocol) updateTraderjoeLend(underlyings []string) error {
	network := prot.ProtocolBasic.Network

	comptroller, err := traderjoeComptroller.NewTraderjoeComptroller(types.ToAddress(ethaddr.TraderjoeJoetrollerList[network]), *prot.ProtocolBasic.Client)
	if err != nil {
		return err
	}
	rewarder, err := traderjoeRewardDistributor.NewTraderjoeRewardDistributor(types.ToAddress(ethaddr.TraderjoeRewardDistributorList[network]), *prot.ProtocolBasic.Client)
	if err != nil {
		return err
	}
	// all ctokens
	allMarkets, err := comptroller.GetAllMarkets(nil)
	if err != nil {
		return err
	}
	joePrice, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol("joe", network, "usd")
	if err != nil {
		return err
	}
	for _, ctoken := range allMarkets {
		var lendPool lend.LendPool
		err = lendPool.Init(*prot.ProtocolBasic)
		if err != nil {
			return err // must be fatal error
		}
		err = lendPool.UpdateTokensByCToken(types.ToString(ctoken))
		if err != nil {
			continue
		}
		// select from underlyings needed
		underlyingAddress := *lendPool.UnderlyingBasic.Address
		if len(underlyings) != 0 {
			if !utils.ContainInArrayX(underlyingAddress, underlyings) {
				continue
			}
		}
		underlyingPriceUSD, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol(*lendPool.UnderlyingBasic.Symbol, network, "usd")
		if err != nil {
			continue
		}
		cToken, err := traderjoeCToken.NewTraderjoeCToken(ctoken, *prot.ProtocolBasic.Client)
		if err != nil {
			continue
		}
		// supply apy
		supplyRatePerSecond, err := cToken.SupplyRatePerSecond(nil)
		if err != nil {
			continue
		}
		lendPool.CToken.DepositApyInfo.Apr = types.ToFloat64(supplyRatePerSecond) * math.Pow10(-18) * constants.SecondsPerYear
		lendPool.CToken.DepositApyInfo.Apy = utils.Apr2Apy(lendPool.CToken.DepositApyInfo.Apr)
		// borrow apy
		borrowRatePerSecond, err := cToken.BorrowRatePerSecond(nil)
		if err != nil {
			continue
		}
		lendPool.CToken.BorrowApyInfo.Apr = types.ToFloat64(borrowRatePerSecond) * math.Pow10(-18) * constants.SecondsPerYear
		lendPool.CToken.BorrowApyInfo.Apy = utils.Apr2Apy(lendPool.CToken.BorrowApyInfo.Apr)

		// apy incentives
		joeSupplyReward, _ := rewarder.RewardSupplySpeeds(nil, 0, types.ToAddress(ctoken))
		joeBorrowReward, _ := rewarder.RewardBorrowSpeeds(nil, 0, types.ToAddress(ctoken))
		joeSupplyPerDay := types.ToFloat64(joeSupplyReward) * 86400 * math.Pow10(-18)
		joeBorrowPerDay := types.ToFloat64(joeBorrowReward) * 86400 * math.Pow10(-18)
		supplyRewardsPerYear := joeSupplyPerDay * 365 * joePrice
		borrowRewardsPerYear := joeBorrowPerDay * 365 * joePrice

		cash, err := cToken.GetCash(nil)
		if err != nil {
			continue
		}
		totalBorrow, err := cToken.TotalBorrows(nil)
		if err != nil {
			continue
		}
		totalBorrows := types.ToFloat64(totalBorrow) * math.Pow10(-*lendPool.UnderlyingBasic.Decimals)
		totalSupply := types.ToFloat64(totalBorrow)*math.Pow10(-*lendPool.UnderlyingBasic.Decimals) + types.ToFloat64(cash)*math.Pow10(-*lendPool.UnderlyingBasic.Decimals)
		supplyAprIncentive := supplyRewardsPerYear / totalSupply / underlyingPriceUSD
		borrowAprIncentive := borrowRewardsPerYear / totalBorrows / underlyingPriceUSD
		if totalSupply == 0 {
			supplyAprIncentive = 0
		}
		if totalBorrows == 0 {
			borrowAprIncentive = 0
		}
		lendPool.CToken.DepositApyInfo.AprIncentive = supplyAprIncentive
		lendPool.CToken.BorrowApyInfo.AprIncentive = borrowAprIncentive
		lendPool.CToken.DepositApyInfo.ApyIncentive = utils.Apr2Apy(lendPool.CToken.DepositApyInfo.AprIncentive)
		lendPool.CToken.BorrowApyInfo.ApyIncentive = utils.Apr2Apy(lendPool.CToken.BorrowApyInfo.AprIncentive)

		prot.LendPools = append(prot.LendPools, &lendPool)
	}
	return nil
}
