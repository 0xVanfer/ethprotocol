package ethprotocol

import (
	"errors"
	"fmt"
	"math"
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
	"github.com/0xVanfer/types"
	"github.com/0xVanfer/utils"
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
	// every update apy func should check network
	network := prot.ProtocolBasic.Network
	supportedNetworks := []string{
		chainId.AvalancheChainName,
		chainId.EthereumChainName,
	}
	if !utils.ContainInArrayX(network, supportedNetworks) {
		fmt.Println("Aave lend V2", network, "not supported.")
		return nil
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
			// avs apr
			lendingPool.AToken.ApyInfo.Apr = types.ToFloat64(poolInfo.LiquidityRate) * math.Pow10(-27)
			lendingPool.VToken.ApyInfo.Apr = types.ToFloat64(poolInfo.VariableBorrowRate) * math.Pow10(-27)
			lendingPool.SToken.ApyInfo.Apr = types.ToFloat64(poolInfo.StableBorrowRate) * math.Pow10(-27)

			// avax apr incentive
			underlyingPriceUSD, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol(*lendingPool.UnderlyingBasic.Symbol, network, "usd")
			if err != nil {
				continue
			}
			aEmissionUSD := types.ToFloat64(poolInfo.AEmissionPerSecond) * constants.SecondsPerYear * chainTokenPrice * math.Pow10(-18)
			if lendingPool.AToken.Basic == nil {
				fmt.Println("atoken of", underlyingAddress, "not found")
				continue
			}
			aSupply, err := lendingPool.AToken.Basic.TotalSupply()
			if err != nil {
				continue
			}
			aSupplyUSD := aSupply * underlyingPriceUSD
			lendingPool.AToken.ApyInfo.AprIncentive = aEmissionUSD / aSupplyUSD

			vEmissionUSD := types.ToFloat64(poolInfo.VEmissionPerSecond) * constants.SecondsPerYear * chainTokenPrice * math.Pow10(-18)
			if lendingPool.VToken.Basic == nil {
				fmt.Println("vtoken of", underlyingAddress, "not found")
				continue
			}
			vSupply, err := lendingPool.VToken.Basic.TotalSupply()
			if err != nil {
				continue
			}
			vSupplyUSD := vSupply * underlyingPriceUSD
			lendingPool.VToken.ApyInfo.AprIncentive = vEmissionUSD / vSupplyUSD

			// apr 2 apy
			lendingPool.AToken.ApyInfo.Apy = utils.Apr2Apy(lendingPool.AToken.ApyInfo.Apr)
			lendingPool.VToken.ApyInfo.Apy = utils.Apr2Apy(lendingPool.VToken.ApyInfo.Apr)
			lendingPool.SToken.ApyInfo.Apy = utils.Apr2Apy(lendingPool.VToken.ApyInfo.Apr)
			lendingPool.AToken.ApyInfo.ApyIncentive = utils.Apr2Apy(lendingPool.AToken.ApyInfo.AprIncentive)
			lendingPool.VToken.ApyInfo.ApyIncentive = utils.Apr2Apy(lendingPool.AToken.ApyInfo.AprIncentive)

			// status todo
			lendingPool.Status.TotalSupply = types.ToFloat64(poolInfo.TotalDeposits) * math.Pow10(-poolInfo.Decimals)
			lendingPool.Status.TotalVBorrow = types.ToFloat64(poolInfo.TotalCurrentVariableDebt) * math.Pow10(-poolInfo.Decimals)
			lendingPool.Status.TotalSBorrow = types.ToFloat64(poolInfo.TotalPrincipalStableDebt) * math.Pow10(-poolInfo.Decimals)
			lendingPool.Status.UtilizationRate = types.ToFloat64(poolInfo.UtilizationRate)
		}
	}
	return nil
}

// Internal use only! No protocol regular check!
//
// Update aave v3 lend pools by underlyings.
func (prot *Protocol) updateAaveV3LendingApy(underlyings []string) error {
	network := prot.ProtocolBasic.Network
	// polygon has different abi
	if network == chainId.PolygonChainName {
		return prot.updateAaveV3LendingApyPolygon(underlyings)
	}

	supportedNetworks := []string{
		chainId.AvalancheChainName,
	}
	if !utils.ContainInArrayX(network, supportedNetworks) {
		fmt.Println("Aave lend V3", network, "not supported.")
		return nil
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

			// status
			supplyCap := types.ToFloat64(assetInfo.SupplyCap)
			suppliedCap := types.ToFloat64(assetInfo.AvailableLiquidity) * math.Pow10(-*lendingPool.UnderlyingBasic.Decimals)
			supplyCapRemain := supplyCap - suppliedCap
			if supplyCap == 0 {
				supplyCapRemain = 0
			}

			totalVBorrow := types.ToFloat64(assetInfo.TotalScaledVariableDebt) * math.Pow10(-types.ToInt(assetInfo.Decimals))
			totalSBorrow := types.ToFloat64(assetInfo.TotalPrincipalStableDebt) * math.Pow10(-types.ToInt(assetInfo.Decimals))

			lendingPool.Status.CollateralFactor = types.ToFloat64(assetInfo.BaseLTVasCollateral) / 100
			lendingPool.Status.LiquidationLimit = types.ToFloat64(assetInfo.ReserveLiquidationThreshold) / 100
			lendingPool.Status.AllowBorrow = assetInfo.BorrowingEnabled
			lendingPool.Status.AllowCollateral = assetInfo.UsageAsCollateralEnabled

			lendingPool.Status.SupplyLimit = supplyCap
			lendingPool.Status.SupplyCapacity = supplyCapRemain
			lendingPool.Status.TotalSupply = suppliedCap + totalVBorrow + totalSBorrow
			lendingPool.Status.TotalVBorrow = totalVBorrow
			lendingPool.Status.TotalSBorrow = totalSBorrow
			lendingPool.Status.UtilizationRate = (totalVBorrow + totalSBorrow) / (suppliedCap + totalVBorrow + totalSBorrow)

			lendingPool.Status.EModeCategoryId = int(assetInfo.EModeCategoryId)
			lendingPool.Status.EModeCollateralFactor = types.ToFloat64(assetInfo.EModeLtv) / 100
			lendingPool.Status.EModeLiquidationLimit = types.ToFloat64(assetInfo.EModeLiquidationThreshold) / 100

			lendingPool.Status.BorrowableInIsolation = assetInfo.BorrowableInIsolation
		}
	}
	return nil
}

// Internal use only! No protocol regular check!
//
// Update aave v3 lend pools by underlyings.
func (prot *Protocol) updateAaveV3LendingApyPolygon(underlyings []string) error {
	network := prot.ProtocolBasic.Network
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

			// status
			supplyCap := types.ToFloat64(assetInfo.SupplyCap)
			suppliedCap := types.ToFloat64(assetInfo.AvailableLiquidity) * math.Pow10(-*lendingPool.UnderlyingBasic.Decimals)
			supplyCapRemain := supplyCap - suppliedCap
			if supplyCap == 0 {
				supplyCapRemain = 0
			}

			totalVBorrow := types.ToFloat64(assetInfo.TotalScaledVariableDebt) * math.Pow10(-types.ToInt(assetInfo.Decimals))
			totalSBorrow := types.ToFloat64(assetInfo.TotalPrincipalStableDebt) * math.Pow10(-types.ToInt(assetInfo.Decimals))

			lendingPool.Status.CollateralFactor = types.ToFloat64(assetInfo.BaseLTVasCollateral) / 100
			lendingPool.Status.LiquidationLimit = types.ToFloat64(assetInfo.ReserveLiquidationThreshold) / 100
			lendingPool.Status.AllowBorrow = assetInfo.BorrowingEnabled
			lendingPool.Status.AllowCollateral = assetInfo.UsageAsCollateralEnabled

			lendingPool.Status.SupplyLimit = supplyCap
			lendingPool.Status.SupplyCapacity = supplyCapRemain
			lendingPool.Status.TotalSupply = suppliedCap + totalVBorrow + totalSBorrow
			lendingPool.Status.TotalVBorrow = totalVBorrow
			lendingPool.Status.TotalSBorrow = totalSBorrow
			lendingPool.Status.UtilizationRate = (totalVBorrow + totalSBorrow) / (suppliedCap + totalVBorrow + totalSBorrow)

			lendingPool.Status.EModeCategoryId = int(assetInfo.EModeCategoryId)
			lendingPool.Status.EModeCollateralFactor = types.ToFloat64(assetInfo.EModeLtv) / 100
			lendingPool.Status.EModeLiquidationLimit = types.ToFloat64(assetInfo.EModeLiquidationThreshold) / 100

			lendingPool.Status.BorrowableInIsolation = assetInfo.BorrowableInIsolation
		}
	}
	return nil
}

// Internal use only! No protocol regular check!
//
// Update benqi lend pools by underlyings.
func (prot *Protocol) updateBenqiLendingApy(underlyings []string) error {
	network := prot.ProtocolBasic.Network
	if !utils.ContainInArrayX(network, []string{chainId.AvalancheChainName}) {
		fmt.Println("Benqi lend", network, "not supported.")
		return nil
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
			lendingPool.CToken.DepositApyInfo.Apr = types.ToFloat64(supplyRatePerSecond) * math.Pow10(-18) * constants.SecondsPerYear
			lendingPool.CToken.DepositApyInfo.Apy = utils.Apr2Apy(lendingPool.CToken.DepositApyInfo.Apr)
			// borrow apy
			borrowRatePerSecond, err := qitoken.BorrowRatePerTimestamp(nil)
			if err != nil {
				continue
			}
			lendingPool.CToken.BorrowApyInfo.Apr = types.ToFloat64(borrowRatePerSecond) * math.Pow10(-18) * constants.SecondsPerYear
			lendingPool.CToken.BorrowApyInfo.Apy = utils.Apr2Apy(lendingPool.CToken.BorrowApyInfo.Apr)
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
			totalBorrows := types.ToFloat64(totalBorrow) * math.Pow10(-*lendingPool.UnderlyingBasic.Decimals)
			totalSupply := types.ToFloat64(totalBorrow)*math.Pow10(-*lendingPool.UnderlyingBasic.Decimals) + types.ToFloat64(cash)*math.Pow10(-*lendingPool.UnderlyingBasic.Decimals)
			supplyAprIncentive := supplyRewardsPerYear / totalSupply / underlyingPriceUSD
			borrowAprIncentive := borrowRewardsPerYear / totalBorrows / underlyingPriceUSD
			if totalSupply == 0 {
				supplyAprIncentive = 0
			}
			if totalBorrows == 0 {
				borrowAprIncentive = 0
			}
			lendingPool.CToken.DepositApyInfo.AprIncentive = supplyAprIncentive
			lendingPool.CToken.BorrowApyInfo.AprIncentive = borrowAprIncentive
			lendingPool.CToken.DepositApyInfo.ApyIncentive = utils.Apr2Apy(lendingPool.CToken.DepositApyInfo.AprIncentive)
			lendingPool.CToken.BorrowApyInfo.ApyIncentive = utils.Apr2Apy(lendingPool.CToken.BorrowApyInfo.AprIncentive)

			// status
			lendingPool.Status.TotalSupply = totalSupply
			lendingPool.Status.TotalCBorrow = totalBorrows
			lendingPool.Status.UtilizationRate = totalBorrows / totalSupply
		}
	}
	return nil
}

// Internal use only! No protocol regular check!
//
// Update compound-like lend pools by underlyings.
func (prot *Protocol) updateTraderjoeLendingApy(underlyings []string) error {
	network := prot.ProtocolBasic.Network
	if !utils.ContainInArrayX(network, []string{chainId.AvalancheChainName}) {
		fmt.Println("Traderjoe lend", network, "not supported.")
		return nil
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
			lendingPool.CToken.DepositApyInfo.Apr = types.ToFloat64(supplyRatePerSecond) * math.Pow10(-18) * constants.SecondsPerYear
			lendingPool.CToken.DepositApyInfo.Apy = utils.Apr2Apy(lendingPool.CToken.DepositApyInfo.Apr)
			// borrow apy
			borrowRatePerSecond, err := cToken.BorrowRatePerSecond(nil)
			if err != nil {
				continue
			}
			lendingPool.CToken.BorrowApyInfo.Apr = types.ToFloat64(borrowRatePerSecond) * math.Pow10(-18) * constants.SecondsPerYear
			lendingPool.CToken.BorrowApyInfo.Apy = utils.Apr2Apy(lendingPool.CToken.BorrowApyInfo.Apr)

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
			totalBorrows := types.ToFloat64(totalBorrow) * math.Pow10(-*lendingPool.UnderlyingBasic.Decimals)
			totalSupply := types.ToFloat64(totalBorrow)*math.Pow10(-*lendingPool.UnderlyingBasic.Decimals) + types.ToFloat64(cash)*math.Pow10(-*lendingPool.UnderlyingBasic.Decimals)
			supplyAprIncentive := supplyRewardsPerYear / totalSupply / underlyingPriceUSD
			borrowAprIncentive := borrowRewardsPerYear / totalBorrows / underlyingPriceUSD
			if totalSupply == 0 {
				supplyAprIncentive = 0
			}
			if totalBorrows == 0 {
				borrowAprIncentive = 0
			}
			lendingPool.CToken.DepositApyInfo.AprIncentive = supplyAprIncentive
			lendingPool.CToken.BorrowApyInfo.AprIncentive = borrowAprIncentive
			lendingPool.CToken.DepositApyInfo.ApyIncentive = utils.Apr2Apy(lendingPool.CToken.DepositApyInfo.AprIncentive)
			lendingPool.CToken.BorrowApyInfo.ApyIncentive = utils.Apr2Apy(lendingPool.CToken.BorrowApyInfo.AprIncentive)

			// status
			lendingPool.Status.TotalSupply = totalSupply
			lendingPool.Status.TotalCBorrow = totalBorrows
			lendingPool.Status.UtilizationRate = totalBorrows / totalSupply
		}
	}
	return nil
}
