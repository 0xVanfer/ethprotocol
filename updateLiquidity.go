package ethprotocol

import (
	"errors"
	"fmt"
	"strings"

	"github.com/0xVanfer/erc"
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/internal/common"
	"github.com/0xVanfer/ethprotocol/internal/constants"
	"github.com/0xVanfer/ethprotocol/internal/requests"
	"github.com/0xVanfer/ethprotocol/liquidity"
	"github.com/0xVanfer/ethprotocol/model"
	"github.com/0xVanfer/types"
	"github.com/0xVanfer/utils"
	"github.com/shopspring/decimal"
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
	// sushi
	case ethaddr.SushiProtocol:
		return prot.updateSushiLiquidity()
	// pangolin
	case ethaddr.PangolinProtocol:
		return prot.updatePangolinLiquidity()
	// axial
	case ethaddr.AxialProtocol:
		return prot.updateAxialLiquidity()
	default:
		return errors.New(prot.ProtocolBasic.ProtocolName + " liquidity pools not supported")
	}
}

func (prot *Protocol) updateTraderjoeLiquidity() error {
	// check network
	network := prot.ProtocolBasic.Network
	err := prot.CheckNetwork()
	if err != nil {
		return err
	}

	// all pools
	poolList, err := requests.ReqJoeAvaxPools()
	if err != nil {
		return err
	}
	for _, pool := range poolList {
		// ignore symbol
		if common.ContainIgnoreSymbol(pool.Token0.Symbol, pool.Token1.Symbol) {
			fmt.Println(pool.Name, "has ignore symbol")
			continue
		}
		// volume 24
		var volume24 decimal.Decimal
		for _, data := range pool.HourData {
			volume24 = volume24.Add(types.ToDecimal(data.VolumeUSD))
		}
		if volume24.IsZero() {
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
			Reserve:    types.ToDecimal(pool.Token0.Volume).Div(constants.WEIUnit),
			ReserveUSD: types.ToDecimal(pool.Token0.VolumeUSD).Div(constants.WEIUnit),
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
			Reserve:    types.ToDecimal(pool.Token1.Volume).Div(constants.WEIUnit),
			ReserveUSD: types.ToDecimal(pool.Token1.VolumeUSD).Div(constants.WEIUnit),
		}

		// apy
		dailyProfit := types.ToDecimal(pool.ReserveUSD).Mul(decimal.NewFromFloat(0.0025)).Mul(volume24)
		apr := dailyProfit.Mul(decimal.NewFromInt(365))
		apyInfo := model.ApyInfo{
			Apr: apr,
			Apy: types.ToDecimal(utils.Apr2Apy(apr)),
		}

		newPool := liquidity.LiquidityPool{
			ProtocolBasic: prot.ProtocolBasic,
			PoolName:      pool.Name,
			LpToken:       lp,
			Tokens:        []*liquidity.TokenOfLp{&token0OfLp, &token1OfLp},
			ApyInfo:       &apyInfo,
			Reserve:       types.ToDecimal(pool.TotalSupply),
			ReserveUSD:    types.ToDecimal(pool.ReserveUSD),
			Volume24:      volume24,
		}
		prot.LiquidityPools = append(prot.LiquidityPools, &newPool)
	}
	return nil
}

func (prot *Protocol) updateSushiLiquidity() error {
	// check network
	network := prot.ProtocolBasic.Network
	err := prot.CheckNetwork()
	if err != nil {
		return err
	}

	poolsInfo, err := requests.ReqSushiPairs(network)
	if err != nil {
		return err
	}
	for _, pool := range poolsInfo {
		name := pool.Pair.Token0.Symbol + " - " + pool.Pair.Token1.Symbol
		// ignore symbol
		if common.ContainIgnoreSymbol(pool.Pair.Token0.Symbol, pool.Pair.Token1.Symbol) {
			fmt.Println(name, "has ignore symbol")
			continue
		}
		// skip 0 volume
		if pool.Volume1D == 0 {
			fmt.Println(name, "0 volume")
			continue
		}
		// lp
		lp, err := erc.NewErc20(pool.Pair.ID, network, *prot.ProtocolBasic.Client)
		if err != nil {
			fmt.Println(name, err)
			continue
		}
		// tokens
		token0Decimals := types.ToInt(pool.Pair.Token0.Decimals)
		token0 := erc.ERC20Info{
			Network:  &network,
			Address:  &pool.Pair.Token0.ID,
			Symbol:   &pool.Pair.Token0.Symbol,
			Decimals: &token0Decimals,
		}
		token0OfLp := liquidity.TokenOfLp{
			Basic:      &token0,
			Underlying: &token0,
		}
		token1Decimals := types.ToInt(pool.Pair.Token1.Decimals)
		token1 := erc.ERC20Info{
			Network:  &network,
			Address:  &pool.Pair.Token1.ID,
			Symbol:   &pool.Pair.Token1.Symbol,
			Decimals: &token1Decimals,
		}
		token1OfLp := liquidity.TokenOfLp{
			Basic:      &token1,
			Underlying: &token1,
		}

		// apy
		apy := types.ToDecimal(strings.ReplaceAll(pool.Apy, "%", "")).Div(decimal.NewFromInt(100))
		apyInfo := model.ApyInfo{
			Apr: types.ToDecimal(utils.Apy2Apr(apy)),
			Apy: apy,
		}

		newPool := liquidity.LiquidityPool{
			ProtocolBasic: prot.ProtocolBasic,
			PoolName:      name,
			LpToken:       lp,
			Tokens:        []*liquidity.TokenOfLp{&token0OfLp, &token1OfLp},
			ApyInfo:       &apyInfo,
			ReserveUSD:    types.ToDecimal(pool.Liquidity),
			Volume24:      types.ToDecimal(pool.Volume1D),
		}
		prot.LiquidityPools = append(prot.LiquidityPools, &newPool)
	}
	return nil
}

func (prot *Protocol) updatePangolinLiquidity() error {
	// check network
	network := prot.ProtocolBasic.Network
	err := prot.CheckNetwork()
	if err != nil {
		return err
	}
	allInfo, err := requests.ReqPangolinAllInfo()
	if err != nil {
		return err
	}

	for _, pool := range allInfo.Data.Minichefs[0].Farms {
		name := pool.Pair.Token0.Symbol + " - " + pool.Pair.Token1.Symbol

		// ignore symbol
		if common.ContainIgnoreSymbol(pool.Pair.Token0.Symbol, pool.Pair.Token1.Symbol) {
			fmt.Println(name, "has ignore symbol")
			continue
		}
		// lp
		lp, err := erc.NewErc20(pool.Pair.ID, network, *prot.ProtocolBasic.Client)
		if err != nil {
			fmt.Println(name, err)
			continue
		}
		// tokens
		token0Decimals := types.ToInt(pool.Pair.Token0.Decimals)
		token0 := erc.ERC20Info{
			Network:  &network,
			Address:  &pool.Pair.Token0.ID,
			Symbol:   &pool.Pair.Token0.Symbol,
			Decimals: &token0Decimals,
		}
		token0OfLp := liquidity.TokenOfLp{
			Basic:      &token0,
			Underlying: &token0,
			Reserve:    types.ToDecimal(pool.Pair.Reserve0),
			ReserveUSD: types.ToDecimal(pool.Pair.Token0.DerivedUSD).Mul(types.ToDecimal(pool.Pair.Reserve0)),
		}
		token1Decimals := types.ToInt(pool.Pair.Token1.Decimals)
		token1 := erc.ERC20Info{
			Network:  &network,
			Address:  &pool.Pair.Token1.ID,
			Symbol:   &pool.Pair.Token1.Symbol,
			Decimals: &token1Decimals,
		}
		token1OfLp := liquidity.TokenOfLp{
			Basic:      &token1,
			Underlying: &token1,
			Reserve:    types.ToDecimal(pool.Pair.Reserve1),
			ReserveUSD: types.ToDecimal(pool.Pair.Token1.DerivedUSD).Mul(types.ToDecimal(pool.Pair.Reserve1)),
		}

		// apy
		apys, err := requests.ReqPangolinApr2(pool.Pid)
		if err != nil {
			fmt.Println(name, err)
			continue
		}
		apr := types.ToDecimal(apys.SwapFeeApr).Div(decimal.NewFromInt(100))
		apyInfo := model.ApyInfo{
			Apr: apr,
			Apy: types.ToDecimal(utils.Apr2Apy(apr)),
		}

		newPool := liquidity.LiquidityPool{
			ProtocolBasic: prot.ProtocolBasic,
			PoolName:      name,
			LpToken:       lp,
			Tokens:        []*liquidity.TokenOfLp{&token0OfLp, &token1OfLp},
			ApyInfo:       &apyInfo,
			Reserve:       types.ToDecimal(pool.Pair.TotalSupply),
			ReserveUSD:    token0OfLp.ReserveUSD.Add(token1OfLp.ReserveUSD),
		}
		prot.LiquidityPools = append(prot.LiquidityPools, &newPool)
	}
	return nil
}

func (prot *Protocol) updateAxialLiquidity() error {
	// check network
	network := prot.ProtocolBasic.Network
	err := prot.CheckNetwork()
	if err != nil {
		return err
	}
	pools, err := requests.ReqAxialAvaxPools()
	if err != nil {
		return err
	}

	for _, pool := range pools {
		// ignore symbol
		if common.ContainIgnoreSymbol(pool.Tokens[0].Symbol, pool.Tokens[1].Symbol) {
			fmt.Println(pool.Symbol, "has ignore symbol")
			continue
		}
		// lp
		lp, err := erc.NewErc20(pool.TokenAddress, network, *prot.ProtocolBasic.Client)
		if err != nil {
			fmt.Println(pool.Symbol, err)
			continue
		}
		// tokens
		var tokens []*liquidity.TokenOfLp
		for _, tokenInfo := range pool.Tokens {
			decimals := types.ToInt(tokenInfo.Decimals)
			token := erc.ERC20Info{
				Network:  &network,
				Address:  &tokenInfo.Address,
				Symbol:   &tokenInfo.Symbol,
				Decimals: &decimals,
			}
			tokenOfLp := liquidity.TokenOfLp{
				Basic:      &token,
				Underlying: &token,
				Reserve:    decimal.Zero, // todo
				ReserveUSD: decimal.Zero, // todo
			}
			tokens = append(tokens, &tokenOfLp)
		}

		// apy
		apr := types.ToDecimal(pool.LastSwapApr).Div(decimal.NewFromInt(100))
		apyInfo := model.ApyInfo{
			Apr: apr,
			Apy: types.ToDecimal(utils.Apr2Apy(apr)),
		}

		newPool := liquidity.LiquidityPool{
			ProtocolBasic: prot.ProtocolBasic,
			PoolName:      pool.Symbol,
			LpToken:       lp,
			Tokens:        tokens,
			ApyInfo:       &apyInfo,
			Reserve:       types.ToDecimal(pool.LastTvl).Div(types.ToDecimal(pool.LastTokenPrice)),
			ReserveUSD:    types.ToDecimal(pool.LastTvl),
		}
		prot.LiquidityPools = append(prot.LiquidityPools, &newPool)
	}
	return nil
}
