package ethprotocol

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/0xVanfer/abigen/aave/aaveUiIncentiveDataProviderV3"
	"github.com/0xVanfer/abigen/aave/aaveUiPoolDataProviderV3"
	"github.com/0xVanfer/abigen/benqi/benqiCToken"
	"github.com/0xVanfer/abigen/benqi/benqiComptroller"
	"github.com/0xVanfer/chainId"
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/internal/apy"
	"github.com/0xVanfer/ethprotocol/internal/constants"
	"github.com/0xVanfer/ethprotocol/internal/requests"
	"github.com/0xVanfer/ethprotocol/lend"
	"github.com/0xVanfer/types"
)

func (prot *Protocol) UpdateLendApys() error {
	if !prot.ProtocolBasic.Regularcheck() {
		return errors.New("protocol basic must be initialized")
	}
	network := prot.ProtocolBasic.Network
	chainTokenPrice, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol(chainId.ChainTokenSymbolList[network], network, "usd")
	if err != nil {
		return err
	}
	switch prot.ProtocolBasic.ProtocolName {
	// aave v2
	case ethaddr.AaveV2Protocol:
		poolsInfo, err := requests.ReqAaveV2LendPools(network)
		if err != nil {
			return err
		}
		for _, poolInfo := range poolsInfo {
			// skip ethereum Amm pools
			if strings.EqualFold(poolInfo.Symbol[:3], "Amm") {
				continue
			}
			underlyingAddress := poolInfo.UnderlyingAsset
			var lendPool lend.LendPool
			err = lendPool.Init(*prot.ProtocolBasic)
			if err != nil {
				return err // must be fatal error
			}
			err = lendPool.UpdateTokensByUnderlying(underlyingAddress)
			if err != nil {
				return err
			}
			lendPool.AToken.ApyInfo.Apr = types.ToFloat64(poolInfo.LiquidityRate) * math.Pow10(-27)
			lendPool.VToken.ApyInfo.Apr = types.ToFloat64(poolInfo.VariableBorrowRate) * math.Pow10(-27)
			lendPool.SToken.ApyInfo.Apr = types.ToFloat64(poolInfo.StableBorrowRate) * math.Pow10(-27)

			underlyingPriceUSD, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol(*lendPool.UnderlyingBasic.Symbol, network, "usd")
			if err != nil {
				fmt.Println(err)
				continue
			}
			aEmissionUSD := types.ToFloat64(poolInfo.AEmissionPerSecond) * constants.SecondsPerYear * chainTokenPrice * math.Pow10(-18)
			if lendPool.AToken.Basic == nil {
				fmt.Println("atoken of", underlyingAddress, "not found")
				continue
			}
			aSupply, err := lendPool.AToken.Basic.TotalSupply()
			if err != nil {
				fmt.Println(err)
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
				fmt.Println(err)
				continue
			}
			vSupplyUSD := vSupply * underlyingPriceUSD
			lendPool.VToken.ApyInfo.AprIncentive = vEmissionUSD / vSupplyUSD

			lendPool.AToken.ApyInfo.Apy = apy.Apr2Apy(lendPool.AToken.ApyInfo.Apr)
			lendPool.VToken.ApyInfo.Apy = apy.Apr2Apy(lendPool.VToken.ApyInfo.Apr)
			lendPool.SToken.ApyInfo.Apy = apy.Apr2Apy(lendPool.VToken.ApyInfo.Apr)
			lendPool.AToken.ApyInfo.ApyIncentive = apy.Apr2Apy(lendPool.AToken.ApyInfo.AprIncentive)
			lendPool.VToken.ApyInfo.ApyIncentive = apy.Apr2Apy(lendPool.AToken.ApyInfo.AprIncentive)

			prot.LendPools = append(prot.LendPools, &lendPool)
		}
		return nil
	// aave v3
	case ethaddr.AaveV3Protocol:
		// address provider, used in contracts
		addressProviderAddress := types.ToAddress(ethaddr.AavePoolAddressProviderV3List[network])
		// get the basic info and rewards of a pool
		uiPoolDataProvider, _ := aaveUiPoolDataProviderV3.NewAaveUiPoolDataProviderV3(types.ToAddress(ethaddr.AaveUiPoolDataProveiderV3List[network]), *prot.ProtocolBasic.Client)
		allInfo, _, _ := uiPoolDataProvider.GetReservesData(nil, addressProviderAddress)
		// get the incentive rewards info of a pool
		uiIncentiveDataProvider, _ := aaveUiIncentiveDataProviderV3.NewAaveUiIncentiveDataProviderV3(types.ToAddress(ethaddr.AaveUiIncentiveDataProveiderV3List[network]), *prot.ProtocolBasic.Client)
		incentiveInfo, _ := uiIncentiveDataProvider.GetReservesIncentivesData(nil, addressProviderAddress)
		for _, assetInfo := range allInfo {
			underlyingAddress := types.ToLowerString(assetInfo.UnderlyingAsset)

			var lendPool lend.LendPool
			err = lendPool.Init(*prot.ProtocolBasic)
			if err != nil {
				return err // must be fatal error
			}
			err = lendPool.UpdateTokensByUnderlying(underlyingAddress)
			if err != nil {
				return err
			}
			underlyingPriceUSD, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol(*lendPool.UnderlyingBasic.Symbol, network, "usd")
			if err != nil {
				fmt.Println(err)
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
				aSupply, err := lendPool.AToken.Basic.TotalSupply()
				if err != nil {
					fmt.Println(err)
					continue
				}
				aSupplyUSD := aSupply * underlyingPriceUSD
				aRewardTokens := incentiveReward.AIncentiveData.RewardsTokenInformation
				for _, aRewardToken := range aRewardTokens {
					rewardTokenPrice, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol(aRewardToken.RewardTokenSymbol, network, "usd")
					if err != nil {
						fmt.Println(err)
						continue
					}
					rewardPerYearUSD := types.ToFloat64(aRewardToken.EmissionPerSecond) * constants.SecondsPerYear * math.Pow10(-types.ToInt(aRewardToken.RewardTokenDecimals)) * rewardTokenPrice
					apy := rewardPerYearUSD / aSupplyUSD
					lendPool.AToken.ApyInfo.ApyIncentive += apy
				}

				vSupply, err := lendPool.VToken.Basic.TotalSupply()
				if err != nil {
					fmt.Println(err)
					continue
				}
				vSupplyUSD := vSupply * underlyingPriceUSD
				vRewardTokens := incentiveReward.VIncentiveData.RewardsTokenInformation
				for _, vRewardToken := range vRewardTokens {
					rewardTokenPrice, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol(vRewardToken.RewardTokenSymbol, network, "usd")
					if err != nil {
						fmt.Println(err)
						continue
					}
					rewardPerYearUSD := types.ToFloat64(vRewardToken.EmissionPerSecond) * constants.SecondsPerYear * math.Pow10(-types.ToInt(vRewardToken.RewardTokenDecimals)) * rewardTokenPrice
					apy := rewardPerYearUSD / vSupplyUSD
					lendPool.VToken.ApyInfo.ApyIncentive += apy
				}
				lendPool.AToken.ApyInfo.AprIncentive = apy.Apy2Apr(lendPool.AToken.ApyInfo.ApyIncentive)
				lendPool.VToken.ApyInfo.AprIncentive = apy.Apy2Apr(lendPool.VToken.ApyInfo.ApyIncentive)
			}
			prot.LendPools = append(prot.LendPools, &lendPool)
		}
		return nil
	case ethaddr.BenqiProtocol:
		comptroller, err := benqiComptroller.NewBenqiComptroller(types.ToAddress(ethaddr.BenqiComptrollerList[network]), *prot.ProtocolBasic.Client)
		if err != nil {
			return err
		}
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
				return err
			}
			underlyingPriceUSD, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol(*lendPool.UnderlyingBasic.Symbol, network, "usd")
			if err != nil {
				fmt.Println(err)
				continue
			}
			qitoken, err := benqiCToken.NewBenqiCToken(ctoken, *prot.ProtocolBasic.Client)
			if err != nil {
				fmt.Println(err)
				continue
			}
			// supply apy
			supplyRatePerSecond, err := qitoken.SupplyRatePerTimestamp(nil)
			if err != nil {
				fmt.Println(err)
				continue
			}
			lendPool.CToken.DepositApyInfo.Apr = types.ToFloat64(supplyRatePerSecond) * math.Pow10(-18) * constants.SecondsPerYear
			lendPool.CToken.DepositApyInfo.Apy = apy.Apr2Apy(lendPool.CToken.DepositApyInfo.Apr)
			// borrow apy
			borrowRatePerSecond, err := qitoken.BorrowRatePerTimestamp(nil)
			if err != nil {
				fmt.Println(err)
				continue
			}
			lendPool.CToken.BorrowApyInfo.Apr = types.ToFloat64(borrowRatePerSecond) * math.Pow10(-18) * constants.SecondsPerYear
			lendPool.CToken.BorrowApyInfo.Apy = apy.Apr2Apy(lendPool.CToken.BorrowApyInfo.Apr)
			// apy incentives
			supplyReward0, err := comptroller.SupplyRewardSpeeds(nil, 0, ctoken)
			if err != nil {
				fmt.Println(err)
				continue
			}
			supplyReward1, err := comptroller.SupplyRewardSpeeds(nil, 1, ctoken)
			if err != nil {
				fmt.Println(err)
				continue
			}
			supplyReward0PerDay := types.ToFloat64(supplyReward0) * 86400 * math.Pow10(-18)
			supplyReward1PerDay := types.ToFloat64(supplyReward1) * 86400 * math.Pow10(-18)
			supplyRewardsPerYear := supplyReward0PerDay*365*qiPriceUSD + supplyReward1PerDay*365*chainTokenPrice

			borrowReward0, err := comptroller.BorrowRewardSpeeds(nil, 0, ctoken)
			if err != nil {
				fmt.Println(err)
				continue
			}
			borrowReward1, err := comptroller.BorrowRewardSpeeds(nil, 1, ctoken)
			if err != nil {
				fmt.Println(err)
				continue
			}
			borrowReward0PerDay := types.ToFloat64(borrowReward0) * 86400 * math.Pow10(-18)
			borrowReward1PerDay := types.ToFloat64(borrowReward1) * 86400 * math.Pow10(-18)
			borrowRewardsPerYear := borrowReward0PerDay*365*qiPriceUSD + borrowReward1PerDay*365*chainTokenPrice

			cash, err := qitoken.GetCash(nil)
			if err != nil {
				fmt.Println(err)
				continue
			}
			totalBorrow, err := qitoken.TotalBorrows(nil)
			if err != nil {
				fmt.Println(err)
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
			lendPool.CToken.DepositApyInfo.ApyIncentive = apy.Apr2Apy(lendPool.CToken.DepositApyInfo.AprIncentive)
			lendPool.CToken.BorrowApyInfo.ApyIncentive = apy.Apr2Apy(lendPool.CToken.BorrowApyInfo.AprIncentive)

			prot.LendPools = append(prot.LendPools, &lendPool)
		}
		return nil
	case ethaddr.CompoundProtocol:
		// todo
		return nil
	case ethaddr.TraderJoeProtocol:
		// todo
		return nil
	default:
		return errors.New("protocol not supported lend pools")
	}
}
