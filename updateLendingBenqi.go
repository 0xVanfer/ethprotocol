package ethprotocol

import (
	"strings"

	"github.com/0xVanfer/abigen/benqi/benqiCToken"
	"github.com/0xVanfer/abigen/benqi/benqiComptroller"
	"github.com/0xVanfer/chainId"
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/internal/constants"
	"github.com/0xVanfer/ethprotocol/model"
	"github.com/0xVanfer/types"
	"github.com/0xVanfer/utils"
	"github.com/shopspring/decimal"
)

// Internal use only! No protocol regular check!
//
// Update benqi lend pools by underlyings.
func (prot *Protocol) updateLendingBenqi(underlyings []string) error {
	// check network
	network := prot.ProtocolBasic.Network
	err := prot.CheckNetwork()
	if err != nil {
		return err
	}
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
			qitoken, err := benqiCToken.NewBenqiCToken(ctoken, *prot.ProtocolBasic.Client)
			if err != nil {
				continue
			}
			// supply apy
			supplyRatePerSecond, err := qitoken.SupplyRatePerTimestamp(nil)
			if err != nil {
				continue
			}
			supplyApr := types.ToDecimal(supplyRatePerSecond).Div(constants.WEIUnit).Mul(constants.SecondsPerYear)
			// borrow apy
			borrowRatePerSecond, err := qitoken.BorrowRatePerTimestamp(nil)
			if err != nil {
				continue
			}
			borrowApr := types.ToDecimal(borrowRatePerSecond).Div(constants.WEIUnit).Mul(constants.SecondsPerYear)
			// apy incentives
			supplyReward0, err := comptroller.SupplyRewardSpeeds(nil, 0, ctoken)
			if err != nil {
				continue
			}
			supplyReward1, err := comptroller.SupplyRewardSpeeds(nil, 1, ctoken)
			if err != nil {
				continue
			}
			supplyReward0PerYear := types.ToDecimal(supplyReward0).Div(constants.WEIUnit).Mul(constants.SecondsPerYear).Mul(qiPriceUSD)
			supplyReward1PerYear := types.ToDecimal(supplyReward1).Div(constants.WEIUnit).Mul(constants.SecondsPerYear).Mul(chainTokenPrice)

			borrowReward0, err := comptroller.BorrowRewardSpeeds(nil, 0, ctoken)
			if err != nil {
				continue
			}
			borrowReward1, err := comptroller.BorrowRewardSpeeds(nil, 1, ctoken)
			if err != nil {
				continue
			}
			borrowReward0PerYear := types.ToDecimal(borrowReward0).Div(constants.WEIUnit).Mul(constants.SecondsPerYear).Mul(qiPriceUSD)
			borrowReward1PerYear := types.ToDecimal(borrowReward1).Div(constants.WEIUnit).Mul(constants.SecondsPerYear).Mul(chainTokenPrice)

			cash, err := qitoken.GetCash(nil)
			if err != nil {
				continue
			}
			// fmt.Println(qitoken.Symbol(nil))

			// fmt.Println("supply 0 per day:", types.ToDecimal(supplyReward0).Div(constants.WEIUnit).Mul(constants.SecondsPerDay))
			// fmt.Println("supply 1 per day:", types.ToDecimal(supplyReward1).Div(constants.WEIUnit).Mul(constants.SecondsPerDay))

			// fmt.Println("borrow 0 per day:", types.ToDecimal(borrowReward0).Div(constants.WEIUnit).Mul(constants.SecondsPerDay))
			// fmt.Println("borrow 1 per day:", types.ToDecimal(borrowReward1).Div(constants.WEIUnit).Mul(constants.SecondsPerDay))

			// fmt.Println("cash:", cash)
			totalBorrow, err := qitoken.TotalBorrows(nil)
			if err != nil {
				continue
			}
			totalBorrows := types.ToDecimal(totalBorrow).Div(decimal.New(1, int32(*lendingPool.UnderlyingBasic.Decimals)))
			totalSupply := totalBorrows.Add(types.ToDecimal(cash).Div(decimal.New(1, int32(*lendingPool.UnderlyingBasic.Decimals))))

			var supplyAprIncentive0 decimal.Decimal
			var supplyAprIncentive1 decimal.Decimal
			var borrowAprIncentive0 decimal.Decimal
			var borrowAprIncentive1 decimal.Decimal

			if totalSupply.IsZero() {
				supplyAprIncentive0 = decimal.Zero
				supplyAprIncentive1 = decimal.Zero
			} else {
				supplyAprIncentive0 = supplyReward0PerYear.Div(totalSupply).Div(underlyingPriceUSD)
				supplyAprIncentive1 = supplyReward1PerYear.Div(totalSupply).Div(underlyingPriceUSD)
			}
			if totalBorrows.IsZero() {
				borrowAprIncentive0 = decimal.Zero
				borrowAprIncentive1 = decimal.Zero
			} else {
				borrowAprIncentive0 = borrowReward0PerYear.Div(totalBorrows).Div(underlyingPriceUSD)
				borrowAprIncentive1 = borrowReward1PerYear.Div(totalBorrows).Div(underlyingPriceUSD)
			}

			// supply
			lendingPool.CToken.SupplyApyInfo = &model.ApyInfo{
				Apr:                supplyApr,
				IncentiveToken0Apr: supplyAprIncentive0,
				IncentiveToken1Apr: supplyAprIncentive1,
			}
			lendingPool.CToken.SupplyApyInfo.Generate()

			// borrow
			lendingPool.CToken.BorrowApyInfo = &model.ApyInfo{
				Apr:                borrowApr,
				IncentiveToken0Apr: borrowAprIncentive0,
				IncentiveToken1Apr: borrowAprIncentive1,
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
