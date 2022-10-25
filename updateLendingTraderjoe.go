package ethprotocol

import (
	"strings"

	"github.com/0xVanfer/abigen/traderjoe/traderjoeCToken"
	"github.com/0xVanfer/abigen/traderjoe/traderjoeComptroller"
	"github.com/0xVanfer/abigen/traderjoe/traderjoeRewardDistributor"
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/internal/constants"
	"github.com/0xVanfer/ethprotocol/model"
	"github.com/0xVanfer/types"
	"github.com/0xVanfer/utils"
	"github.com/shopspring/decimal"
)

// Internal use only! No protocol regular check!
//
// Update traderjoe lend pools by underlyings.
func (prot *Protocol) updateLendingTraderjoe(underlyings []string) error {
	// check network
	network := prot.ProtocolBasic.Network
	err := prot.CheckNetwork()
	if err != nil {
		return err
	}
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
		// locate the pool
		for _, lendingPool := range prot.LendingPools {
			if !strings.EqualFold(*lendingPool.CToken.Basic.Address, types.ToString(ctoken)) {
				continue
			}
			// select from underlyings needed
			underlyingAddress := *lendingPool.UnderlyingBasic.Address
			if len(underlyings) != 0 {
				if !utils.ContainInArrayX(underlyingAddress, underlyings) {
					continue
				}
			}
			underlyingPriceUSD, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol(*lendingPool.UnderlyingBasic.Symbol, network, "usd")
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
			supplyApr := types.ToDecimal(supplyRatePerSecond).Div(constants.WEIUnit).Mul(constants.SecondsPerYear)

			// borrow apy
			borrowRatePerSecond, err := cToken.BorrowRatePerSecond(nil)
			if err != nil {
				continue
			}
			borrowApr := types.ToDecimal(borrowRatePerSecond).Div(constants.WEIUnit).Mul(constants.SecondsPerYear)

			// apy incentives
			joeSupplyReward, _ := rewarder.RewardSupplySpeeds(nil, 0, types.ToAddress(ctoken))
			joeBorrowReward, _ := rewarder.RewardBorrowSpeeds(nil, 0, types.ToAddress(ctoken))
			joeSupplyPerDay := types.ToDecimal(joeSupplyReward).Mul(constants.SecondsPerDay).Div(constants.WEIUnit)
			joeBorrowPerDay := types.ToDecimal(joeBorrowReward).Mul(constants.SecondsPerDay).Div(constants.WEIUnit)
			supplyRewardsPerYear := joeSupplyPerDay.Mul(decimal.NewFromInt(365)).Mul(joePrice)
			borrowRewardsPerYear := joeBorrowPerDay.Mul(decimal.NewFromInt(365)).Mul(joePrice)

			cash, err := cToken.GetCash(nil)
			if err != nil {
				continue
			}
			totalBorrow, err := cToken.TotalBorrows(nil)
			if err != nil {
				continue
			}

			totalBorrows := types.ToDecimal(totalBorrow).Div(decimal.New(1, int32(*lendingPool.UnderlyingBasic.Decimals)))
			totalSupply := totalBorrows.Add(types.ToDecimal(cash).Div(decimal.New(1, int32(*lendingPool.UnderlyingBasic.Decimals))))

			var supplyAprIncentive decimal.Decimal
			var borrowAprIncentive decimal.Decimal

			if totalSupply.IsZero() {
				supplyAprIncentive = decimal.Zero
			} else {
				supplyAprIncentive = supplyRewardsPerYear.Div(totalSupply).Div(underlyingPriceUSD)
			}
			if totalBorrows.IsZero() {
				borrowAprIncentive = decimal.Zero
			} else {
				borrowAprIncentive = borrowRewardsPerYear.Div(totalBorrows).Div(underlyingPriceUSD)
			}

			// supply
			lendingPool.CToken.SupplyApyInfo = &model.ApyInfo{
				Apr:               supplyApr,
				IncentiveTotalApr: supplyAprIncentive,
			}
			lendingPool.CToken.SupplyApyInfo.Generate()

			// borrow
			lendingPool.CToken.BorrowApyInfo = &model.ApyInfo{
				Apr:               borrowApr,
				IncentiveTotalApr: borrowAprIncentive,
			}
			lendingPool.CToken.BorrowApyInfo.Generate()

			// status
			lendingPool.Status.TotalSupply = totalSupply
			lendingPool.Status.TotalCBorrow = totalBorrows
			lendingPool.Status.UtilizationRate = totalBorrows.Div(totalSupply)
		}
	}
	return nil
}
