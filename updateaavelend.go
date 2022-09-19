package ethprotocol

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/0xVanfer/abigen/aave/aaveUiIncentiveDataProviderV3"
	"github.com/0xVanfer/abigen/aave/aaveUiPoolDataProviderV3"
	"github.com/0xVanfer/chainId"
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/internal/apy"
	"github.com/0xVanfer/ethprotocol/internal/common"
	"github.com/0xVanfer/ethprotocol/internal/constants"
	"github.com/0xVanfer/ethprotocol/internal/requests"
	"github.com/0xVanfer/ethprotocol/lend"
	"github.com/0xVanfer/types"
)

// If underlyings is empty, read all pools.
func (prot *Protocol) updateAaveV2Lend(underlyings []string) error {
	if !prot.ProtocolBasic.Regularcheck() {
		return errors.New("protocol basic must be initialized")
	}
	network := prot.ProtocolBasic.Network
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
			if !common.ContainInArrayWithoutCapital(underlyingAddress, underlyings) {
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
		lendPool.AToken.ApyInfo.Apy = apy.Apr2Apy(lendPool.AToken.ApyInfo.Apr)
		lendPool.VToken.ApyInfo.Apy = apy.Apr2Apy(lendPool.VToken.ApyInfo.Apr)
		lendPool.SToken.ApyInfo.Apy = apy.Apr2Apy(lendPool.VToken.ApyInfo.Apr)
		lendPool.AToken.ApyInfo.ApyIncentive = apy.Apr2Apy(lendPool.AToken.ApyInfo.AprIncentive)
		lendPool.VToken.ApyInfo.ApyIncentive = apy.Apr2Apy(lendPool.AToken.ApyInfo.AprIncentive)

		prot.LendPools = append(prot.LendPools, &lendPool)
	}
	return nil
}
func (prot *Protocol) updateAaveV3Lend(underlyings []string) error {
	if !prot.ProtocolBasic.Regularcheck() {
		return errors.New("protocol basic must be initialized")
	}
	network := prot.ProtocolBasic.Network

	// address provider, used in contracts
	addressProviderAddress := types.ToAddress(ethaddr.AavePoolAddressProviderV3List[network])
	// get the basic info and rewards of a pool
	uiPoolDataProvider, _ := aaveUiPoolDataProviderV3.NewAaveUiPoolDataProviderV3(types.ToAddress(ethaddr.AaveUiPoolDataProveiderV3List[network]), *prot.ProtocolBasic.Client)
	allInfo, _, _ := uiPoolDataProvider.GetReservesData(nil, addressProviderAddress)
	// get the incentive rewards info of a pool
	uiIncentiveDataProvider, _ := aaveUiIncentiveDataProviderV3.NewAaveUiIncentiveDataProviderV3(types.ToAddress(ethaddr.AaveUiIncentiveDataProveiderV3List[network]), *prot.ProtocolBasic.Client)
	incentiveInfo, _ := uiIncentiveDataProvider.GetReservesIncentivesData(nil, addressProviderAddress)
	for _, assetInfo := range allInfo {
		// select from underlyings needed
		underlyingAddress := types.ToLowerString(assetInfo.UnderlyingAsset)
		if len(underlyings) != 0 {
			if !common.ContainInArrayWithoutCapital(underlyingAddress, underlyings) {
				continue
			}
		}
		var lendPool lend.LendPool
		err := lendPool.Init(*prot.ProtocolBasic)
		if err != nil {
			return err // must be fatal error
		}
		err = lendPool.UpdateTokensByUnderlying(underlyingAddress)
		if err != nil {
			return err
		}
		underlyingPriceUSD, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol(*lendPool.UnderlyingBasic.Symbol, network, "usd")
		if err != nil {
			continue
		}
		lendPool.AToken.ApyInfo.Apy = types.ToFloat64(assetInfo.LiquidityIndex) * types.ToFloat64(assetInfo.LiquidityRate) * math.Pow10(-54)
		lendPool.VToken.ApyInfo.Apy = types.ToFloat64(assetInfo.VariableBorrowIndex) * types.ToFloat64(assetInfo.VariableBorrowRate) * math.Pow10(-54)
		lendPool.AToken.ApyInfo.Apr = apy.Apy2Apr(lendPool.AToken.ApyInfo.Apy)
		lendPool.VToken.ApyInfo.Apr = apy.Apy2Apr(lendPool.VToken.ApyInfo.Apy)

		for _, incentiveReward := range incentiveInfo {
			if !strings.EqualFold(types.ToString(incentiveReward.UnderlyingAsset), underlyingAddress) {
				continue
			}
			aSupply, _ := lendPool.AToken.Basic.TotalSupply()
			aSupplyUSD := aSupply * underlyingPriceUSD
			aRewardTokens := incentiveReward.AIncentiveData.RewardsTokenInformation
			for _, aRewardToken := range aRewardTokens {
				rewardTokenPrice, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol(aRewardToken.RewardTokenSymbol, network, "usd")
				if err != nil {
					continue
				}
				rewardPerYearUSD := types.ToFloat64(aRewardToken.EmissionPerSecond) * constants.SecondsPerYear * math.Pow10(-types.ToInt(aRewardToken.RewardTokenDecimals)) * rewardTokenPrice
				apy := rewardPerYearUSD / aSupplyUSD
				if aSupplyUSD == 0 {
					apy = 0
				}
				lendPool.AToken.ApyInfo.ApyIncentive += apy
			}

			vSupply, _ := lendPool.VToken.Basic.TotalSupply()
			vSupplyUSD := vSupply * underlyingPriceUSD
			vRewardTokens := incentiveReward.VIncentiveData.RewardsTokenInformation
			for _, vRewardToken := range vRewardTokens {
				rewardTokenPrice, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol(vRewardToken.RewardTokenSymbol, network, "usd")
				if err != nil {
					continue
				}
				rewardPerYearUSD := types.ToFloat64(vRewardToken.EmissionPerSecond) * constants.SecondsPerYear * math.Pow10(-types.ToInt(vRewardToken.RewardTokenDecimals)) * rewardTokenPrice
				apy := rewardPerYearUSD / vSupplyUSD
				if vSupplyUSD == 0 {
					apy = 0
				}
				lendPool.VToken.ApyInfo.ApyIncentive += apy
			}
			lendPool.AToken.ApyInfo.AprIncentive = apy.Apy2Apr(lendPool.AToken.ApyInfo.ApyIncentive)
			lendPool.VToken.ApyInfo.AprIncentive = apy.Apy2Apr(lendPool.VToken.ApyInfo.ApyIncentive)
		}
		prot.LendPools = append(prot.LendPools, &lendPool)
	}
	return nil
}
