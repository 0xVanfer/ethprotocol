package ethprotocol

import (
	"fmt"
	"math"
	"strings"

	"github.com/0xVanfer/chainId"
	"github.com/0xVanfer/ethprotocol/internal/constants"
	"github.com/0xVanfer/ethprotocol/internal/requests"
	"github.com/0xVanfer/ethprotocol/lend"
	"github.com/0xVanfer/types"
	"github.com/0xVanfer/utils"
)

// Internal use only! No protocol regular check!
//
// Update aave v2 lend pools by underlyings.
func (prot *Protocol) updateAaveV2Lend(underlyings []string) error {
	network := prot.ProtocolBasic.Network
	if !utils.ContainInArrayX(network, []string{chainId.AvalancheChainName, chainId.EthereumChainName}) {
		fmt.Println("Aave lend V2", network, "not supported.")
		return nil
	}
	// chain token price
	chainTokenPrice, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol(chainId.ChainTokenSymbolList[network], network, "usd")
	if err != nil {
		return err
	}
	poolsInfo, err := requests.ReqAaveV2LendPools(network)
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
		// new a pool
		var lendPool lend.LendPool
		err = lendPool.Init(*prot.ProtocolBasic)
		if err != nil {
			return err // must be fatal error
		}
		err = lendPool.UpdateTokensByUnderlying(underlyingAddress)
		if err != nil {
			return err
		}
		// avs apr
		lendPool.AToken.ApyInfo.Apr = types.ToFloat64(poolInfo.LiquidityRate) * math.Pow10(-27)
		lendPool.VToken.ApyInfo.Apr = types.ToFloat64(poolInfo.VariableBorrowRate) * math.Pow10(-27)
		lendPool.SToken.ApyInfo.Apr = types.ToFloat64(poolInfo.StableBorrowRate) * math.Pow10(-27)

		// avax apr incentive
		underlyingPriceUSD, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol(*lendPool.UnderlyingBasic.Symbol, network, "usd")
		if err != nil {
			continue
		}
		aEmissionUSD := types.ToFloat64(poolInfo.AEmissionPerSecond) * constants.SecondsPerYear * chainTokenPrice * math.Pow10(-18)
		if lendPool.AToken.Basic == nil {
			fmt.Println("atoken of", underlyingAddress, "not found")
			continue
		}
		aSupply, err := lendPool.AToken.Basic.TotalSupply()
		if err != nil {
			continue
		}
		aSupplyUSD := aSupply * underlyingPriceUSD
		lendPool.AToken.ApyInfo.AprIncentive = aEmissionUSD / aSupplyUSD

		vEmissionUSD := types.ToFloat64(poolInfo.VEmissionPerSecond) * constants.SecondsPerYear * chainTokenPrice * math.Pow10(-18)
		if lendPool.VToken.Basic == nil {
			fmt.Println("vtoken of", underlyingAddress, "not found")
			continue
		}
		vSupply, err := lendPool.VToken.Basic.TotalSupply()
		if err != nil {
			continue
		}
		vSupplyUSD := vSupply * underlyingPriceUSD
		lendPool.VToken.ApyInfo.AprIncentive = vEmissionUSD / vSupplyUSD

		// apr 2 apy
		lendPool.AToken.ApyInfo.Apy = utils.Apr2Apy(lendPool.AToken.ApyInfo.Apr)
		lendPool.VToken.ApyInfo.Apy = utils.Apr2Apy(lendPool.VToken.ApyInfo.Apr)
		lendPool.SToken.ApyInfo.Apy = utils.Apr2Apy(lendPool.VToken.ApyInfo.Apr)
		lendPool.AToken.ApyInfo.ApyIncentive = utils.Apr2Apy(lendPool.AToken.ApyInfo.AprIncentive)
		lendPool.VToken.ApyInfo.ApyIncentive = utils.Apr2Apy(lendPool.AToken.ApyInfo.AprIncentive)

		prot.LendPools = append(prot.LendPools, &lendPool)
	}
	return nil
}
