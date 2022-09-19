package ethprotocol

import (
	"errors"
	"fmt"
	"math"
	"strings"

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
	chainTokenPrice, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol(chainId.ChainTokenSymbolList[prot.ProtocolBasic.Network], prot.ProtocolBasic.Network, "usd")
	if err != nil {
		return err
	}
	switch prot.ProtocolBasic.ProtocolName {
	// aave v2
	case ethaddr.AaveV2Protocol:
		poolsInfo, err := requests.ReqAaveV2LendPools(prot.ProtocolBasic.Network)
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
			lendPool.UpdateTokensByUnderlying(underlyingAddress)

			lendPool.AToken.ApyInfo.Apr = types.ToFloat64(poolInfo.LiquidityRate) * math.Pow10(-27)
			lendPool.VToken.ApyInfo.Apr = types.ToFloat64(poolInfo.VariableBorrowRate) * math.Pow10(-27)
			lendPool.SToken.ApyInfo.Apr = types.ToFloat64(poolInfo.StableBorrowRate) * math.Pow10(-27)

			underlyingPriceUSD, err := prot.ProtocolBasic.Gecko.GetPriceBySymbol(*lendPool.UnderlyingBasic.Symbol, prot.ProtocolBasic.Network, "usd")
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
	case ethaddr.AaveV3Protocol:
		// todo
		return nil
	case ethaddr.BenqiProtocol:
		// todo
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
