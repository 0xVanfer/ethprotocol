package lend

// import (
// 	"errors"
// 	"math"
// 	"strings"

// 	"github.com/0xVanfer/abigen/aave/aaveUiIncentiveDataProviderV3"
// 	"github.com/0xVanfer/abigen/aave/aaveUiPoolDataProviderV3"
// 	"github.com/0xVanfer/chainId"
// 	"github.com/0xVanfer/coingecko"
// 	"github.com/0xVanfer/ethaddr"
// 	"github.com/0xVanfer/ethprotocol/internal/apy"
// 	"github.com/0xVanfer/ethprotocol/internal/constants"
// 	"github.com/0xVanfer/ethprotocol/internal/requests"
// 	"github.com/0xVanfer/types"
// )

// // BETTER NOT use this for a list of lend pools.
// func (p *LendPool) GetApys(gecko *coingecko.Gecko) error {
// 	chaintokenSymbol := chainId.ChainTokenSymbolList[p.Network]
// 	chainTokenPrice, err := gecko.GetPriceBySymbol(chaintokenSymbol, p.Network, "usd")
// 	if err != nil {
// 		return err
// 	}
// 	if p.PoolType.IsAaveLike {
// 		switch p.ProtocolName {

// 		// aave v2
// 		case ethaddr.AaveV2Protocol:
// 			aaveGraphs, err := requests.ReqAaveV2Pools(p.Network)
// 			if err != nil {
// 				return err
// 			}
// 			for _, reserves := range aaveGraphs {
// 				// skip ethereum Amm pools
// 				if strings.EqualFold(reserves.Symbol[:3], "Amm") {
// 					continue
// 				}
// 				// choose the correct lending pool
// 				if !strings.EqualFold(p.UnderlyingBasic.Address, reserves.UnderlyingAsset) {
// 					continue
// 				}
// 				p.AToken.ApyInfo.Apr = types.ToFloat64(reserves.LiquidityRate) * math.Pow10(-27)
// 				p.VToken.ApyInfo.Apr = types.ToFloat64(reserves.VariableBorrowRate) * math.Pow10(-27)
// 				p.SToken.ApyInfo.Apr = types.ToFloat64(reserves.StableBorrowRate) * math.Pow10(-27)

// 				p.AToken.ApyInfo.Apy = apy.Apr2Apy(p.AToken.ApyInfo.Apr)
// 				p.VToken.ApyInfo.Apy = apy.Apr2Apy(p.VToken.ApyInfo.Apr)
// 				p.SToken.ApyInfo.Apy = apy.Apr2Apy(p.SToken.ApyInfo.Apr)

// 				// atoken incentive
// 				aEmissionUSD := types.ToFloat64(reserves.AEmissionPerSecond) * constants.SecondsPerYear * chainTokenPrice * math.Pow10(-18)
// 				aTotalSupply, err := p.AToken.Basic.TotalSupply()
// 				if err != nil {
// 					return err
// 				}
// 				aPrice, err := p.AToken.UnderlyingBasic.PriceUSD(gecko)
// 				if err != nil {
// 					return err
// 				}
// 				aTotalSupplyUSD := aTotalSupply * aPrice
// 				p.AToken.ApyInfo.AprIncentive = aEmissionUSD / aTotalSupplyUSD
// 				p.AToken.ApyInfo.ApyIncentive = apy.Apr2Apy(p.AToken.ApyInfo.AprIncentive)

// 				// vtoken incentive
// 				vEmissionUSD := types.ToFloat64(reserves.VEmissionPerSecond) * constants.SecondsPerYear * chainTokenPrice * math.Pow10(-18)
// 				vTotalSupply, err := p.VToken.Basic.TotalSupply()
// 				if err != nil {
// 					return err
// 				}
// 				vPrice, err := p.VToken.UnderlyingBasic.PriceUSD(gecko)
// 				if err != nil {
// 					return err
// 				}
// 				vTotalSupplyUSD := vTotalSupply * vPrice
// 				p.VToken.ApyInfo.AprIncentive = vEmissionUSD / vTotalSupplyUSD
// 				p.VToken.ApyInfo.ApyIncentive = apy.Apr2Apy(p.VToken.ApyInfo.AprIncentive)
// 			}
// 			return nil
// 		// aave v3
// 		case ethaddr.AaveV3Protocol:
// 			// address provider, used in contracts
// 			addressProviderAddress := types.ToAddress(ethaddr.AavePoolAddressProviderV3List[p.Network])
// 			// get the basic info and rewards of a pool
// 			uiPoolDataProvider, err := aaveUiPoolDataProviderV3.NewAaveUiPoolDataProviderV3(types.ToAddress(ethaddr.AaveUiPoolDataProveiderV3List[p.Network]), p.Client)
// 			if err != nil {
// 				return err
// 			}
// 			allInfo, _, err := uiPoolDataProvider.GetReservesData(nil, addressProviderAddress)
// 			if err != nil {
// 				return err
// 			}
// 			// get the incentive rewards info of a pool
// 			uiIncentiveDataProvider, err := aaveUiIncentiveDataProviderV3.NewAaveUiIncentiveDataProviderV3(types.ToAddress(ethaddr.AaveUiIncentiveDataProveiderV3List[p.Network]), p.Client)
// 			if err != nil {
// 				return err
// 			}
// 			incentiveInfo, err := uiIncentiveDataProvider.GetReservesIncentivesData(nil, addressProviderAddress)
// 			if err != nil {
// 				return err
// 			}
// 			// base apy
// 			for _, assetInfo := range allInfo {
// 				if !strings.EqualFold(types.ToLowerString(assetInfo.UnderlyingAsset), p.UnderlyingBasic.Address) {
// 					continue
// 				}
// 				p.AToken.ApyInfo.Apy = types.ToFloat64(assetInfo.LiquidityIndex) * types.ToFloat64(assetInfo.LiquidityRate) * math.Pow10(-54)
// 				p.VToken.ApyInfo.Apy = types.ToFloat64(assetInfo.VariableBorrowIndex) * types.ToFloat64(assetInfo.VariableBorrowRate) * math.Pow10(-54)

// 			}
// 			// incentive apy
// 			for _, incentiveReward := range incentiveInfo {
// 				if !strings.EqualFold(types.ToLowerString(incentiveReward.UnderlyingAsset), p.UnderlyingBasic.Address) {
// 					continue
// 				}
// 				aTotalSupply, err := p.AToken.Basic.TotalSupply()
// 				if err != nil {
// 					return err
// 				}
// 				aPrice, err := p.AToken.UnderlyingBasic.PriceUSD(gecko)
// 				if err != nil {
// 					return err
// 				}
// 				aTotalSupplyUSD := aTotalSupply * aPrice
// 				aRewardTokens := incentiveReward.AIncentiveData.RewardsTokenInformation
// 				for _, aRewardToken := range aRewardTokens {
// 					if aTotalSupplyUSD == 0 {
// 						break
// 					}
// 					rewardTokenPrice, err := gecko.GetPriceBySymbol(aRewardToken.RewardTokenSymbol, p.Network, "usd")
// 					if err != nil {
// 						return err
// 					}
// 					rewardPerYearUSD := types.ToFloat64(aRewardToken.EmissionPerSecond) * constants.SecondsPerYear * math.Pow10(-types.ToInt(aRewardToken.RewardTokenDecimals)) * rewardTokenPrice
// 					apy := rewardPerYearUSD / aTotalSupplyUSD
// 					p.AToken.ApyInfo.ApyIncentive += apy
// 				}
// 				p.AToken.ApyInfo.AprIncentive = apy.Apy2Apr(p.AToken.ApyInfo.AprIncentive)

// 				vTotalSupply, err := p.VToken.Basic.TotalSupply()
// 				if err != nil {
// 					return err
// 				}
// 				vPrice, err := p.VToken.UnderlyingBasic.PriceUSD(gecko)
// 				if err != nil {
// 					return err
// 				}
// 				vTotalSupplyUSD := vTotalSupply * vPrice
// 				vRewardTokens := incentiveReward.VIncentiveData.RewardsTokenInformation
// 				for _, vRewardToken := range vRewardTokens {
// 					if vTotalSupplyUSD == 0 {
// 						break
// 					}
// 					rewardTokenPrice, err := gecko.GetPriceBySymbol(vRewardToken.RewardTokenSymbol, p.Network, "usd")
// 					if err != nil {
// 						return err
// 					}
// 					rewardPerYearUSD := types.ToFloat64(vRewardToken.EmissionPerSecond) * constants.SecondsPerYear * math.Pow10(-types.ToInt(vRewardToken.RewardTokenDecimals)) * rewardTokenPrice
// 					apy := rewardPerYearUSD / vTotalSupplyUSD
// 					p.VToken.ApyInfo.ApyIncentive += apy
// 				}
// 				p.VToken.ApyInfo.AprIncentive = apy.Apy2Apr(p.VToken.ApyInfo.AprIncentive)
// 			}
// 			return nil
// 		default:
// 			return errors.New("not supported protocol " + p.ProtocolName)
// 		}
// 	}
// 	if p.PoolType.IsCompoundLike {
// 		switch p.ProtocolName {
// 		// compound
// 		case ethaddr.CompoundProtocol:
// 			// todo
// 			return nil
// 		// benqi
// 		case ethaddr.BenqiProtocol:
// 			// todo
// 			return nil
// 		// traderjoe
// 		case ethaddr.TraderJoeProtocol:
// 			// todo
// 			return nil
// 		default:
// 			return errors.New("not supported protocol " + p.ProtocolName)
// 		}
// 	}
// 	return errors.New("must be either aave-like or compound-like")
// }
