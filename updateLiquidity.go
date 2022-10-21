package ethprotocol

import (
	"errors"
	"fmt"

	"github.com/0xVanfer/chainId"
	"github.com/0xVanfer/erc"
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/internal/constants"
	"github.com/0xVanfer/ethprotocol/internal/requests"
	"github.com/0xVanfer/ethprotocol/liquidity"
	"github.com/0xVanfer/ethprotocol/model"
	"github.com/0xVanfer/types"
	"github.com/0xVanfer/utils"
)

func (prot *Protocol) UpdateLiquidity() error {
	// protocol basic must be initialized
	if !prot.ProtocolBasic.Regularcheck() {
		return errors.New("protocol basic must be initialized")
	}
	switch prot.ProtocolBasic.ProtocolName {
	// traderjoe
	case ethaddr.TraderJoeProtocol:
		return prot.updateTraderjoeLiquidity()
	default:
		return errors.New(prot.ProtocolBasic.ProtocolName + " liquidity pools not supported")
	}
}

func (prot *Protocol) updateTraderjoeLiquidity() error {
	network := prot.ProtocolBasic.Network
	supportedNetworks := []string{
		chainId.AvalancheChainName,
	}
	if !utils.ContainInArrayX(network, supportedNetworks) {
		fmt.Println("Traderjoe", network, "not supported.")
		return nil
	}

	// all pools
	poolList, err := requests.ReqJoeAvaxPools()
	if err != nil {
		return err
	}
	for _, pool := range poolList {
		// ignore symbol
		if utils.ContainInArrayX(pool.Token0.Symbol, constants.IgnoreSymbols) || utils.ContainInArrayX(pool.Token1.Symbol, constants.IgnoreSymbols) {
			fmt.Println(pool.Name, "has ignore symbol")
			continue
		}
		// volume 24
		var volume24 float64
		for _, data := range pool.HourData {
			volume24 += types.ToFloat64(data.VolumeUSD)
		}
		if volume24 == 0 {
			fmt.Println(pool.Name, "volume is 0")
			continue
		}
		// lp
		lp, err := erc.NewErc20(pool.ID, network, *prot.ProtocolBasic.Client)
		if err != nil {
			fmt.Println(pool.Name, err)
			continue
		}
		// tokens
		token0Decimals := types.ToInt(pool.Token0.Decimals)
		token0 := erc.ERC20Info{
			Network:  &network,
			Address:  &pool.Token0.ID,
			Symbol:   &pool.Token0.Symbol,
			Decimals: &token0Decimals,
		}
		token0OfLp := liquidity.TokenOfLp{
			Basic:      &token0,
			Underlying: &token0,
			Reserve:    types.ToFloat64(pool.Token0.Volume) * 1e-18,
			ReserveUSD: types.ToFloat64(pool.Token0.VolumeUSD) * 1e-18,
		}
		token1Decimals := types.ToInt(pool.Token1.Decimals)
		token1 := erc.ERC20Info{
			Network:  &network,
			Address:  &pool.Token1.ID,
			Symbol:   &pool.Token1.Symbol,
			Decimals: &token1Decimals,
		}
		token1OfLp := liquidity.TokenOfLp{
			Basic:      &token1,
			Underlying: &token1,
			Reserve:    types.ToFloat64(pool.Token1.Volume) * 1e-18,
			ReserveUSD: types.ToFloat64(pool.Token1.VolumeUSD) * 1e-18,
		}

		// apy
		dailyProfit := 0.0025 * volume24 / types.ToFloat64(pool.ReserveUSD)
		apr := dailyProfit * 365
		apyInfo := model.ApyInfo{
			Base: &model.ApyBase{
				Apr: apr,
				Apy: utils.Apr2Apy(apr),
			},
		}

		newPool := liquidity.LiquidityPool{
			ProtocolBasic: prot.ProtocolBasic,
			PoolName:      pool.Name,
			LpToken:       lp,
			Tokens:        []*liquidity.TokenOfLp{&token0OfLp, &token1OfLp},
			ApyInfo:       &apyInfo,
			Reserve:       types.ToFloat64(pool.TotalSupply),
			ReserveUSD:    types.ToFloat64(pool.ReserveUSD),
			Volume24:      volume24,
		}
		prot.LiquidityPools = append(prot.LiquidityPools, &newPool)
	}
	return nil
}
