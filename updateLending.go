package ethprotocol

import (
	"errors"
	"fmt"
	"strings"

	"github.com/0xVanfer/abigen/aave/aaveUiIncentiveDataProviderV3"
	"github.com/0xVanfer/abigen/aave/aaveUiPoolDataProviderV3"
	"github.com/0xVanfer/abigen/aave/aaveUiPoolDataProviderV3Polygon"
	"github.com/0xVanfer/abigen/benqi/benqiCToken"
	"github.com/0xVanfer/abigen/benqi/benqiComptroller"
	"github.com/0xVanfer/abigen/traderjoe/traderjoeCToken"
	"github.com/0xVanfer/abigen/traderjoe/traderjoeComptroller"
	"github.com/0xVanfer/abigen/traderjoe/traderjoeRewardDistributor"
	"github.com/0xVanfer/chainId"
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/internal/constants"
	"github.com/0xVanfer/ethprotocol/internal/requests"
	"github.com/0xVanfer/ethprotocol/lending"
	"github.com/0xVanfer/ethprotocol/model"
	"github.com/0xVanfer/types"
	"github.com/0xVanfer/utils"
	"github.com/shopspring/decimal"
)

// Update lending pool tokens.
func (prot *Protocol) UpdateLendingPoolTokens() error {
	// protocol basic must be initialized
	if !prot.ProtocolBasic.Regularcheck() {
		return errors.New("protocol basic must be initialized")
	}
	tokenLists := ethaddr.LendingTokenListsMap[prot.ProtocolBasic.ProtocolName]
	if !tokenLists.RegularCheck() {
		return errors.New("protocol not supported lending")
	}

	// use a token map or c token map to get the underlying list
	var tokenListMap map[string]map[string]string

	if tokenLists.ATokenList != nil {
		tokenListMap = *ethaddr.AaveLikeATokenListMap[prot.ProtocolBasic.ProtocolName]
	} else if tokenLists.CTokenList != nil {
		tokenListMap = *ethaddr.CompoundLikeCTokenListMap[prot.ProtocolBasic.ProtocolName]
	} else {
		return errors.New("protocol not supported")
	}
	// use underlying to update tokens
	for underlying := range tokenListMap[prot.ProtocolBasic.Network] {
		var newPool lending.LendingPool
		// will not return err
		_ = newPool.Init(*prot.ProtocolBasic)
		_ = newPool.UpdateTokensByUnderlying(underlying)
		prot.LendingPools = append(prot.LendingPools, &newPool)
	}
	return nil
}

// Update some of the protocol's lend pools apys by given underlying addresses.
//
// If "underlyings" is empty, update all pools.
func (prot *Protocol) UpdateLendingApys(underlyings ...string) error {
	// protocol basic must be initialized
	if !prot.ProtocolBasic.Regularcheck() {
		return errors.New("protocol basic must be initialized")
	}
	switch prot.ProtocolBasic.ProtocolName {
	// aave v2
	case ethaddr.AaveV2Protocol:
		return prot.updateAaveV2LendingApy(underlyings)
	// aave v3
	case ethaddr.AaveV3Protocol:
		return prot.updateAaveV3LendingApy(underlyings)
	// benqi
	case ethaddr.BenqiProtocol:
		return prot.updateBenqiLendingApy(underlyings)
	// tradejoe
	case ethaddr.TraderJoeProtocol:
		return prot.updateTraderjoeLendingApy(underlyings)
	default:
		return errors.New(prot.ProtocolBasic.ProtocolName + " lend pools not supported")
	}
}

// Internal use only! No protocol regular check!
//
// Update aave v2 lend pools by underlyings.
func (prot *Protocol) updateAaveV2LendingApy(underlyings []string) error {
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

// Internal use only! No protocol regular check!
//
// Update aave v3 lend pools by underlyings.
func (prot *Protocol) updateAaveV3LendingApy(underlyings []string) error {
	// check network
	network := prot.ProtocolBasic.Network
	err := prot.CheckNetwork()
	if err != nil {
		return err
	}
	// polygon has different abi
	if network == chainId.PolygonChainName {
		return prot.updateAaveV3LendingApyPolygon(underlyings)
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
func (prot *Protocol) updateAaveV3LendingApyPolygon(underlyings []string) error {
	// check network
	network := prot.ProtocolBasic.Network
	err := prot.CheckNetwork()
	if err != nil {
		return err
	}
	if network != chainId.PolygonChainName {
		return prot.updateAaveV3LendingApy(underlyings)
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

// Internal use only! No protocol regular check!
//
// Update benqi lend pools by underlyings.
func (prot *Protocol) updateBenqiLendingApy(underlyings []string) error {
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

// Internal use only! No protocol regular check!
//
// Update traderjoe lend pools by underlyings.
func (prot *Protocol) updateTraderjoeLendingApy(underlyings []string) error {
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
