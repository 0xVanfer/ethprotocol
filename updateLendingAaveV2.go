package ethprotocol

import (
	"fmt"
	"strings"

	"github.com/0xVanfer/ethprotocol/internal/constants"
	"github.com/0xVanfer/ethprotocol/internal/requests"
	"github.com/0xVanfer/ethprotocol/model"
	"github.com/0xVanfer/types"
	"github.com/0xVanfer/utils"
	"github.com/shopspring/decimal"
)

// Internal use only! No protocol regular check!
//
// Update aave v2 lend pools by underlyings.
func (prot *Protocol) updateLendingAaveV2(underlyings []string) error {
	// check network
	network := prot.ProtocolBasic.Network
	err := prot.CheckNetwork()
	if err != nil {
		return err
	}

	// chain token price
	chainTokenPrice, err := prot.ProtocolBasic.Gecko.GetChainTokenPrice(network, "usd")
	if err != nil {
		return err
	}
	poolsInfo, err := requests.ReqAaveV2LendingPools(network)
	if err != nil {
		return err
	}
	// for all pools
	for _, poolInfo := range poolsInfo {
		// skip ethereum Amm pools
		if strings.EqualFold(poolInfo.Symbol[:3], "Amm") {
			continue
		}
		// select from underlyings needed
		underlyingAddress := poolInfo.UnderlyingAsset
		if len(underlyings) != 0 {
			if !utils.ContainInArrayX(underlyingAddress, underlyings) {
				continue
			}
		}
		// locate the pool
		for _, lendingPool := range prot.LendingPools {
			if !strings.EqualFold(*lendingPool.UnderlyingBasic.Address, underlyingAddress) {
				continue
			}

			// avax apr incentive
			underlyingPriceUSD, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol(*lendingPool.UnderlyingBasic.Symbol, network, "usd")
			if err != nil {
				continue
			}
			aEmissionUSD := types.ToDecimal(poolInfo.AEmissionPerSecond).Mul(constants.SecondsPerYear).Mul(chainTokenPrice).Div(constants.WEIUnit)
			if lendingPool.AToken.Basic == nil {
				fmt.Println("atoken of", underlyingAddress, "not found")
				continue
			}
			aSupply, err := lendingPool.AToken.Basic.TotalSupply()
			if err != nil {
				continue
			}

			vEmissionUSD := types.ToDecimal(poolInfo.VEmissionPerSecond).Mul(constants.SecondsPerYear).Mul(chainTokenPrice).Div(constants.WEIUnit)
			if lendingPool.VToken.Basic == nil {
				fmt.Println("vtoken of", underlyingAddress, "not found")
				continue
			}
			vSupply, err := lendingPool.VToken.Basic.TotalSupply()
			if err != nil {
				continue
			}
			// a token
			lendingPool.AToken.ApyInfo = &model.ApyInfo{
				Apr:               types.ToDecimal(poolInfo.LiquidityRate).Div(constants.RAYUnit),
				IncentiveTotalApr: aEmissionUSD.Div(aSupply).Div(underlyingPriceUSD),
			}
			lendingPool.AToken.ApyInfo.Generate()

			// v token
			lendingPool.VToken.ApyInfo = &model.ApyInfo{
				Apr:               types.ToDecimal(poolInfo.VariableBorrowRate).Div(constants.RAYUnit),
				IncentiveTotalApr: vEmissionUSD.Div(vSupply).Div(underlyingPriceUSD),
			}
			lendingPool.VToken.ApyInfo.Generate()

			// s token
			lendingPool.SToken.ApyInfo = &model.ApyInfo{
				Apr: types.ToDecimal(poolInfo.StableBorrowRate).Div(constants.RAYUnit),
			}
			lendingPool.SToken.ApyInfo.Generate()

			// status todo
			lendingPool.Status.TotalSupply = types.ToDecimal(poolInfo.TotalDeposits).Div(decimal.New(1, int32(poolInfo.Decimals)))
			lendingPool.Status.TotalVBorrow = types.ToDecimal(poolInfo.TotalCurrentVariableDebt).Div(decimal.New(1, int32(poolInfo.Decimals)))
			lendingPool.Status.TotalSBorrow = types.ToDecimal(poolInfo.TotalPrincipalStableDebt).Div(decimal.New(1, int32(poolInfo.Decimals)))
			lendingPool.Status.UtilizationRate = types.ToDecimal(poolInfo.UtilizationRate)
		}
	}
	return nil
}
