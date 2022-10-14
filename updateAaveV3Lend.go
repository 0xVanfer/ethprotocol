package ethprotocol

import (
	"fmt"
	"math"
	"strings"

	"github.com/0xVanfer/abigen/aave/aaveUiIncentiveDataProviderV3"
	"github.com/0xVanfer/abigen/aave/aaveUiPoolDataProviderV3"
	"github.com/0xVanfer/abigen/aave/aaveUiPoolDataProviderV3Polygon"
	"github.com/0xVanfer/chainId"
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/internal/constants"
	"github.com/0xVanfer/ethprotocol/lending"
	"github.com/0xVanfer/types"
	"github.com/0xVanfer/utils"
)

// Internal use only! No protocol regular check!
//
// Update aave v3 lend pools by underlyings.
func (prot *Protocol) updateAaveV3Lend(underlyings []string) error {
	network := prot.ProtocolBasic.Network
	if !utils.ContainInArrayX(network, []string{chainId.AvalancheChainName, chainId.PolygonChainName}) {
		fmt.Println("Aave lend V3", network, "not supported.")
		return nil
	}
	// polygon has different abi
	if network == chainId.PolygonChainName {
		return prot.updateAaveV3LendPolygon(underlyings)
	}
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
			if !utils.ContainInArrayX(underlyingAddress, underlyings) {
				continue
			}
		}
		var lendingPool lending.LendingPool
		err := lendingPool.Init(*prot.ProtocolBasic)
		if err != nil {
			return err // must be fatal error
		}
		err = lendingPool.UpdateTokensByUnderlying(underlyingAddress)
		if err != nil {
			return err
		}
		underlyingPriceUSD, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol(*lendingPool.UnderlyingBasic.Symbol, network, "usd")
		if err != nil {
			continue
		}
		lendingPool.AToken.ApyInfo.Apy = types.ToFloat64(assetInfo.LiquidityIndex) * types.ToFloat64(assetInfo.LiquidityRate) * math.Pow10(-54)
		lendingPool.VToken.ApyInfo.Apy = types.ToFloat64(assetInfo.VariableBorrowIndex) * types.ToFloat64(assetInfo.VariableBorrowRate) * math.Pow10(-54)
		lendingPool.AToken.ApyInfo.Apr = utils.Apy2Apr(lendingPool.AToken.ApyInfo.Apy)
		lendingPool.VToken.ApyInfo.Apr = utils.Apy2Apr(lendingPool.VToken.ApyInfo.Apy)

		for _, incentiveReward := range incentiveInfo {
			if !strings.EqualFold(types.ToString(incentiveReward.UnderlyingAsset), underlyingAddress) {
				continue
			}
			aSupply, _ := lendingPool.AToken.Basic.TotalSupply()
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
				lendingPool.AToken.ApyInfo.ApyIncentive += apy
			}

			vSupply, _ := lendingPool.VToken.Basic.TotalSupply()
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
				lendingPool.VToken.ApyInfo.ApyIncentive += apy
			}
			lendingPool.AToken.ApyInfo.AprIncentive = utils.Apy2Apr(lendingPool.AToken.ApyInfo.ApyIncentive)
			lendingPool.VToken.ApyInfo.AprIncentive = utils.Apy2Apr(lendingPool.VToken.ApyInfo.ApyIncentive)
		}
		prot.LendingPools = append(prot.LendingPools, &lendingPool)
	}
	return nil
}

// Internal use only! No protocol regular check!
//
// Update aave v3 lend pools by underlyings.
func (prot *Protocol) updateAaveV3LendPolygon(underlyings []string) error {
	network := prot.ProtocolBasic.Network
	if network != chainId.PolygonChainName {
		return prot.updateAaveV3Lend(underlyings)
	}
	// address provider, used in contracts
	addressProviderAddress := types.ToAddress(ethaddr.AavePoolAddressProviderV3List[network])
	// get the basic info and rewards of a pool
	uiPoolDataProvider, _ := aaveUiPoolDataProviderV3Polygon.NewAaveUiPoolDataProviderV3Polygon(types.ToAddress(ethaddr.AaveUiPoolDataProveiderV3List[network]), *prot.ProtocolBasic.Client)
	allInfo, _, _ := uiPoolDataProvider.GetReservesData(nil, addressProviderAddress)
	// get the incentive rewards info of a pool
	uiIncentiveDataProvider, _ := aaveUiIncentiveDataProviderV3.NewAaveUiIncentiveDataProviderV3(types.ToAddress(ethaddr.AaveUiIncentiveDataProveiderV3List[network]), *prot.ProtocolBasic.Client)
	incentiveInfo, _ := uiIncentiveDataProvider.GetReservesIncentivesData(nil, addressProviderAddress)
	for _, assetInfo := range allInfo {
		// select from underlyings needed
		underlyingAddress := types.ToLowerString(assetInfo.UnderlyingAsset)
		if len(underlyings) != 0 {
			if !utils.ContainInArrayX(underlyingAddress, underlyings) {
				continue
			}
		}
		var lendingPool lending.LendingPool
		err := lendingPool.Init(*prot.ProtocolBasic)
		if err != nil {
			return err // must be fatal error
		}
		err = lendingPool.UpdateTokensByUnderlying(underlyingAddress)
		if err != nil {
			return err
		}
		underlyingPriceUSD, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol(*lendingPool.UnderlyingBasic.Symbol, network, "usd")
		if err != nil {
			continue
		}
		lendingPool.AToken.ApyInfo.Apy = types.ToFloat64(assetInfo.LiquidityIndex) * types.ToFloat64(assetInfo.LiquidityRate) * math.Pow10(-54)
		lendingPool.VToken.ApyInfo.Apy = types.ToFloat64(assetInfo.VariableBorrowIndex) * types.ToFloat64(assetInfo.VariableBorrowRate) * math.Pow10(-54)
		lendingPool.AToken.ApyInfo.Apr = utils.Apy2Apr(lendingPool.AToken.ApyInfo.Apy)
		lendingPool.VToken.ApyInfo.Apr = utils.Apy2Apr(lendingPool.VToken.ApyInfo.Apy)

		for _, incentiveReward := range incentiveInfo {
			if !strings.EqualFold(types.ToString(incentiveReward.UnderlyingAsset), underlyingAddress) {
				continue
			}
			aSupply, _ := lendingPool.AToken.Basic.TotalSupply()
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
				lendingPool.AToken.ApyInfo.ApyIncentive += apy
			}

			vSupply, _ := lendingPool.VToken.Basic.TotalSupply()
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
				lendingPool.VToken.ApyInfo.ApyIncentive += apy
			}
			lendingPool.AToken.ApyInfo.AprIncentive = utils.Apy2Apr(lendingPool.AToken.ApyInfo.ApyIncentive)
			lendingPool.VToken.ApyInfo.AprIncentive = utils.Apy2Apr(lendingPool.VToken.ApyInfo.ApyIncentive)
		}
		prot.LendingPools = append(prot.LendingPools, &lendingPool)
	}
	return nil
}
