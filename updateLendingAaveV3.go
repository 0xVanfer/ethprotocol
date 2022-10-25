package ethprotocol

import (
	"strings"

	"github.com/0xVanfer/abigen/aave/aaveUiIncentiveDataProviderV3"
	"github.com/0xVanfer/abigen/aave/aaveUiPoolDataProviderV3"
	"github.com/0xVanfer/abigen/aave/aaveUiPoolDataProviderV3Polygon"
	"github.com/0xVanfer/chainId"
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/internal/constants"
	"github.com/0xVanfer/ethprotocol/lending"
	"github.com/0xVanfer/ethprotocol/model"
	"github.com/0xVanfer/types"
	"github.com/0xVanfer/utils"
	"github.com/shopspring/decimal"
)

// Internal use only! No protocol regular check!
//
// Update aave v3 lend pools by underlyings.
func (prot *Protocol) updateLendingAaveV3(underlyings []string) error {
	// check network
	network := prot.ProtocolBasic.Network
	err := prot.CheckNetwork()
	if err != nil {
		return err
	}
	// polygon has different abi
	if network == chainId.PolygonChainName {
		return prot.updateLendingAaveV3Polygon(underlyings)
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
		// locate the pool
		for _, lendingPool := range prot.LendingPools {
			if !strings.EqualFold(*lendingPool.UnderlyingBasic.Address, underlyingAddress) {
				continue
			}
			underlyingPriceUSD, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol(*lendingPool.UnderlyingBasic.Symbol, network, "usd")
			if err != nil {
				continue
			}

			for _, incentiveReward := range incentiveInfo {
				if !strings.EqualFold(types.ToString(incentiveReward.UnderlyingAsset), underlyingAddress) {
					continue
				}
				aSupply, _ := lendingPool.AToken.Basic.TotalSupply()
				aSupplyUSD := aSupply.Mul(underlyingPriceUSD)
				aRewardTokens := incentiveReward.AIncentiveData.RewardsTokenInformation
				var incentiveTotalAApy decimal.Decimal
				if aSupplyUSD.IsZero() {
					incentiveTotalAApy = decimal.Zero
				} else {
					for _, aRewardToken := range aRewardTokens {
						rewardTokenPrice, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol(aRewardToken.RewardTokenSymbol, network, "usd")
						if err != nil {
							continue
						}
						rewardPerYearUSD := types.ToDecimal(aRewardToken.EmissionPerSecond).Mul(constants.SecondsPerYear).Div(decimal.New(1, int32(aRewardToken.RewardTokenDecimals))).Mul(rewardTokenPrice) // todo int32 may have bug
						apy := rewardPerYearUSD.Div(aSupplyUSD)
						incentiveTotalAApy = incentiveTotalAApy.Add(apy)
					}
				}

				vSupply, _ := lendingPool.VToken.Basic.TotalSupply()
				vSupplyUSD := vSupply.Mul(underlyingPriceUSD)
				vRewardTokens := incentiveReward.VIncentiveData.RewardsTokenInformation
				var incentiveTotalVApy decimal.Decimal
				if vSupplyUSD.IsZero() {
					incentiveTotalVApy = decimal.Zero
				} else {
					for _, vRewardToken := range vRewardTokens {
						rewardTokenPrice, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol(vRewardToken.RewardTokenSymbol, network, "usd")
						if err != nil {
							continue
						}
						rewardPerYearUSD := types.ToDecimal(vRewardToken.EmissionPerSecond).Mul(constants.SecondsPerYear).Div(decimal.New(1, int32(vRewardToken.RewardTokenDecimals))).Mul(rewardTokenPrice) // todo int32 may have bug
						apy := rewardPerYearUSD.Div(vSupplyUSD)
						incentiveTotalVApy = incentiveTotalVApy.Add(apy)
					}
				}

				lendingPool.AToken.ApyInfo = &model.ApyInfo{
					Apy:               types.ToDecimal(assetInfo.LiquidityIndex).Mul(types.ToDecimal(assetInfo.LiquidityRate)).Div(constants.RAYUnit).Div(constants.RAYUnit),
					IncentiveTotalApy: incentiveTotalAApy,
				}
				lendingPool.AToken.ApyInfo.Generate()
				lendingPool.VToken.ApyInfo = &model.ApyInfo{
					Apy:               types.ToDecimal(assetInfo.VariableBorrowIndex).Mul(types.ToDecimal(assetInfo.VariableBorrowRate)).Div(constants.RAYUnit).Div(constants.RAYUnit),
					IncentiveTotalApy: incentiveTotalVApy,
				}
				lendingPool.VToken.ApyInfo.Generate()
			}

			// status
			supplyCap := types.ToDecimal(assetInfo.SupplyCap)
			suppliedCap := types.ToDecimal(assetInfo.AvailableLiquidity).Div(decimal.New(1, int32(*lendingPool.UnderlyingBasic.Decimals)))
			supplyCapRemain := supplyCap.Sub(suppliedCap)
			if supplyCap.IsZero() {
				supplyCapRemain = decimal.Zero
			}

			totalVBorrow := types.ToDecimal(assetInfo.TotalScaledVariableDebt).Div(decimal.New(1, int32(types.ToInt64((assetInfo.Decimals)))))
			totalSBorrow := types.ToDecimal(assetInfo.TotalPrincipalStableDebt).Div(decimal.New(1, int32(types.ToInt64((assetInfo.Decimals)))))

			lendingPool.Status = lending.LendingPoolStatus{
				CollateralFactor: types.ToDecimal(assetInfo.BaseLTVasCollateral).Div(decimal.NewFromInt(100)),
				LiquidationLimit: types.ToDecimal(assetInfo.ReserveLiquidationThreshold).Div(decimal.NewFromInt(100)),
				AllowBorrow:      assetInfo.BorrowingEnabled,
				AllowCollateral:  assetInfo.UsageAsCollateralEnabled,

				SupplyLimit:     supplyCap,
				SupplyCapacity:  supplyCapRemain,
				TotalSupply:     suppliedCap.Add(totalVBorrow).Add(totalSBorrow),
				TotalVBorrow:    totalVBorrow,
				TotalSBorrow:    totalSBorrow,
				UtilizationRate: (totalVBorrow.Add(totalSBorrow)).Div(suppliedCap.Add(totalVBorrow).Add(totalSBorrow)),

				EModeCategoryId:       int(assetInfo.EModeCategoryId),
				EModeCollateralFactor: types.ToDecimal(assetInfo.EModeLtv).Div(decimal.NewFromInt(100)),
				EModeLiquidationLimit: types.ToDecimal(assetInfo.EModeLiquidationThreshold).Div(decimal.NewFromInt(100)),

				BorrowableInIsolation: assetInfo.BorrowableInIsolation,
			}
		}
	}
	return nil
}

// Internal use only! No protocol regular check!
//
// Update aave v3 lend pools by underlyings.
func (prot *Protocol) updateLendingAaveV3Polygon(underlyings []string) error {
	// check network
	network := prot.ProtocolBasic.Network
	err := prot.CheckNetwork()
	if err != nil {
		return err
	}
	if network != chainId.PolygonChainName {
		return prot.updateLendingAaveV3(underlyings)
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
		// locate the pool
		for _, lendingPool := range prot.LendingPools {
			if !strings.EqualFold(*lendingPool.UnderlyingBasic.Address, underlyingAddress) {
				continue
			}
			underlyingPriceUSD, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol(*lendingPool.UnderlyingBasic.Symbol, network, "usd")
			if err != nil {
				continue
			}

			for _, incentiveReward := range incentiveInfo {
				if !strings.EqualFold(types.ToString(incentiveReward.UnderlyingAsset), underlyingAddress) {
					continue
				}
				aSupply, _ := lendingPool.AToken.Basic.TotalSupply()
				aSupplyUSD := aSupply.Mul(underlyingPriceUSD)
				aRewardTokens := incentiveReward.AIncentiveData.RewardsTokenInformation
				var incentiveTotalAApy decimal.Decimal
				if aSupplyUSD.IsZero() {
					incentiveTotalAApy = decimal.Zero
				} else {
					for _, aRewardToken := range aRewardTokens {
						rewardTokenPrice, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol(aRewardToken.RewardTokenSymbol, network, "usd")
						if err != nil {
							continue
						}
						rewardPerYearUSD := types.ToDecimal(aRewardToken.EmissionPerSecond).Mul(constants.SecondsPerYear).Div(decimal.New(1, int32(aRewardToken.RewardTokenDecimals))).Mul(rewardTokenPrice) // todo int32 may have bug
						apy := rewardPerYearUSD.Div(aSupplyUSD)
						incentiveTotalAApy = incentiveTotalAApy.Add(apy)
					}
				}

				vSupply, _ := lendingPool.VToken.Basic.TotalSupply()
				vSupplyUSD := vSupply.Mul(underlyingPriceUSD)
				vRewardTokens := incentiveReward.VIncentiveData.RewardsTokenInformation
				var incentiveTotalVApy decimal.Decimal
				if vSupplyUSD.IsZero() {
					incentiveTotalVApy = decimal.Zero
				} else {
					for _, vRewardToken := range vRewardTokens {
						rewardTokenPrice, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol(vRewardToken.RewardTokenSymbol, network, "usd")
						if err != nil {
							continue
						}
						rewardPerYearUSD := types.ToDecimal(vRewardToken.EmissionPerSecond).Mul(constants.SecondsPerYear).Div(decimal.New(1, int32(vRewardToken.RewardTokenDecimals))).Mul(rewardTokenPrice) // todo int32 may have bug
						apy := rewardPerYearUSD.Div(vSupplyUSD)
						incentiveTotalVApy = incentiveTotalVApy.Add(apy)
					}
				}

				lendingPool.AToken.ApyInfo = &model.ApyInfo{
					Apy:               types.ToDecimal(assetInfo.LiquidityIndex).Mul(types.ToDecimal(assetInfo.LiquidityRate)).Div(constants.RAYUnit).Div(constants.RAYUnit),
					IncentiveTotalApy: incentiveTotalAApy,
				}
				lendingPool.AToken.ApyInfo.Generate()
				lendingPool.VToken.ApyInfo = &model.ApyInfo{
					Apy:               types.ToDecimal(assetInfo.VariableBorrowIndex).Mul(types.ToDecimal(assetInfo.VariableBorrowRate)).Div(constants.RAYUnit).Div(constants.RAYUnit),
					IncentiveTotalApy: incentiveTotalVApy,
				}
				lendingPool.VToken.ApyInfo.Generate()
			}

			// status
			supplyCap := types.ToDecimal(assetInfo.SupplyCap)
			suppliedCap := types.ToDecimal(assetInfo.AvailableLiquidity).Div(decimal.New(1, int32(*lendingPool.UnderlyingBasic.Decimals)))
			supplyCapRemain := supplyCap.Sub(suppliedCap)
			if supplyCap.IsZero() {
				supplyCapRemain = decimal.Zero
			}

			totalVBorrow := types.ToDecimal(assetInfo.TotalScaledVariableDebt).Div(decimal.New(1, int32(types.ToInt64((assetInfo.Decimals)))))
			totalSBorrow := types.ToDecimal(assetInfo.TotalPrincipalStableDebt).Div(decimal.New(1, int32(types.ToInt64((assetInfo.Decimals)))))

			lendingPool.Status = lending.LendingPoolStatus{
				CollateralFactor: types.ToDecimal(assetInfo.BaseLTVasCollateral).Div(decimal.NewFromInt(100)),
				LiquidationLimit: types.ToDecimal(assetInfo.ReserveLiquidationThreshold).Div(decimal.NewFromInt(100)),
				AllowBorrow:      assetInfo.BorrowingEnabled,
				AllowCollateral:  assetInfo.UsageAsCollateralEnabled,

				SupplyLimit:     supplyCap,
				SupplyCapacity:  supplyCapRemain,
				TotalSupply:     suppliedCap.Add(totalVBorrow).Add(totalSBorrow),
				TotalVBorrow:    totalVBorrow,
				TotalSBorrow:    totalSBorrow,
				UtilizationRate: (totalVBorrow.Add(totalSBorrow)).Div(suppliedCap.Add(totalVBorrow).Add(totalSBorrow)),

				EModeCategoryId:       int(assetInfo.EModeCategoryId),
				EModeCollateralFactor: types.ToDecimal(assetInfo.EModeLtv).Div(decimal.NewFromInt(100)),
				EModeLiquidationLimit: types.ToDecimal(assetInfo.EModeLiquidationThreshold).Div(decimal.NewFromInt(100)),

				BorrowableInIsolation: assetInfo.BorrowableInIsolation,
			}
		}
	}
	return nil
}
