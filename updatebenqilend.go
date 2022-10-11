package ethprotocol

import (
	"math"

	"github.com/0xVanfer/abigen/benqi/benqiCToken"
	"github.com/0xVanfer/abigen/benqi/benqiComptroller"
	"github.com/0xVanfer/chainId"
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/internal/constants"
	"github.com/0xVanfer/ethprotocol/lend"
	"github.com/0xVanfer/types"
	"github.com/0xVanfer/utils"
)

// Internal use only! No protocol regular check!
// Update benqi lend pools by underlyings.
func (prot *Protocol) updateBenqiLend(underlyings []string) error {
	network := prot.ProtocolBasic.Network
	// avax price
	chainTokenPrice, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol(chainId.ChainTokenSymbolList[network], network, "usd")
	if err != nil {
		return err
	}
	comptroller, err := benqiComptroller.NewBenqiComptroller(types.ToAddress(ethaddr.BenqiComptrollerList[network]), *prot.ProtocolBasic.Client)
	if err != nil {
		return err
	}
	// all ctokens
	allMarkets, err := comptroller.GetAllMarkets(nil)
	if err != nil {
		return err
	}
	qiPriceUSD, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol("qi", network, "usd")
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
		qitoken, err := benqiCToken.NewBenqiCToken(ctoken, *prot.ProtocolBasic.Client)
		if err != nil {
			continue
		}
		// supply apy
		supplyRatePerSecond, err := qitoken.SupplyRatePerTimestamp(nil)
		if err != nil {
			continue
		}
		lendPool.CToken.DepositApyInfo.Apr = types.ToFloat64(supplyRatePerSecond) * math.Pow10(-18) * constants.SecondsPerYear
		lendPool.CToken.DepositApyInfo.Apy = utils.Apr2Apy(lendPool.CToken.DepositApyInfo.Apr)
		// borrow apy
		borrowRatePerSecond, err := qitoken.BorrowRatePerTimestamp(nil)
		if err != nil {
			continue
		}
		lendPool.CToken.BorrowApyInfo.Apr = types.ToFloat64(borrowRatePerSecond) * math.Pow10(-18) * constants.SecondsPerYear
		lendPool.CToken.BorrowApyInfo.Apy = utils.Apr2Apy(lendPool.CToken.BorrowApyInfo.Apr)
		// apy incentives
		supplyReward0, err := comptroller.SupplyRewardSpeeds(nil, 0, ctoken)
		if err != nil {
			continue
		}
		supplyReward1, err := comptroller.SupplyRewardSpeeds(nil, 1, ctoken)
		if err != nil {
			continue
		}
		supplyReward0PerDay := types.ToFloat64(supplyReward0) * 86400 * math.Pow10(-18)
		supplyReward1PerDay := types.ToFloat64(supplyReward1) * 86400 * math.Pow10(-18)
		supplyRewardsPerYear := supplyReward0PerDay*365*qiPriceUSD + supplyReward1PerDay*365*chainTokenPrice

		borrowReward0, err := comptroller.BorrowRewardSpeeds(nil, 0, ctoken)
		if err != nil {
			continue
		}
		borrowReward1, err := comptroller.BorrowRewardSpeeds(nil, 1, ctoken)
		if err != nil {
			continue
		}
		borrowReward0PerDay := types.ToFloat64(borrowReward0) * 86400 * math.Pow10(-18)
		borrowReward1PerDay := types.ToFloat64(borrowReward1) * 86400 * math.Pow10(-18)
		borrowRewardsPerYear := borrowReward0PerDay*365*qiPriceUSD + borrowReward1PerDay*365*chainTokenPrice

		cash, err := qitoken.GetCash(nil)
		if err != nil {
			continue
		}
		totalBorrow, err := qitoken.TotalBorrows(nil)
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
