package requests

import (
	"strings"

	"github.com/0xVanfer/types"
	"github.com/imroc/req"
)

type pangolinApr2Req struct {
	SwapFeeApr  int `json:"swapFeeApr"`  // swapping apr
	StakingApr  int `json:"stakingApr"`  // staking apr
	CombinedApr int `json:"combinedApr"` // both swapping apr and staking apr
}

func ReqPangolinApr2[T types.Integer | string](pid T) (v pangolinApr2Req, err error) {
	url := "https://api.pangolin.exchange/pangolin/apr2/" + types.ToString(pid)
	r, _ := req.Get(url)
	err = r.ToJSON(&v)
	return
}

type pangolinAllInfoReq struct {
	Data struct {
		Minichefs []struct {
			ID                string             `json:"id"`                // chef address
			TotalAllocPoint   string             `json:"totalAllocPoint"`   // total alloc point
			RewardPerSecond   string             `json:"rewardPerSecond"`   // reward per second
			RewardsExpiration string             `json:"rewardsExpiration"` //
			Farms             []pangolinFarmInfo `json:"farms"`             // farms
		} `json:"minichefs"`
	} `json:"data"`
}

type pangolinFarmInfo struct {
	ID              string `json:"id"`              // farm address
	Pid             string `json:"pid"`             // pool id
	Tvl             string `json:"tvl"`             // total reserve usd
	AllocPoint      string `json:"allocPoint"`      // farm alloc point
	RewarderAddress string `json:"rewarderAddress"` // rewarder address
	ChefAddress     string `json:"chefAddress"`     // chef address
	PairAddress     string `json:"pairAddress"`     // pair address
	Rewarder        struct {
		ID      string `json:"id"` // 0x???-0x???
		Rewards []struct {
			ID         string        `json:"id"`         // 0x???-0x???
			Token      pangolinToken `json:"token"`      // token info
			Multiplier string        `json:"multiplier"` // int, 1
		} `json:"rewards"` // all rewards
	} `json:"rewarder"` // rewarder info
	Pair pangolinPairInfo `json:"pair"` // pair info
}

type pangolinToken struct {
	ID         string `json:"id"`         // token address
	Symbol     string `json:"symbol"`     // token symbol
	DerivedUSD string `json:"derivedUSD"` // token price
	Name       string `json:"name"`       // token name
	Decimals   string `json:"decimals"`   // token decimals
}

type pangolinPairInfo struct {
	ID          string        `json:"id"`          // address
	Reserve0    string        `json:"reserve0"`    // token0 reserve
	Reserve1    string        `json:"reserve1"`    // token1 reserve
	TotalSupply string        `json:"totalSupply"` // lp total supply
	Token0      pangolinToken `json:"token0"`      // token 0
	Token1      pangolinToken `json:"token1"`      // token1
}

func ReqPangolinAllInfo() (v pangolinAllInfoReq, err error) {
	url := `https://api.thegraph.com/subgraphs/name/sarjuhansaliya/minichefv2-dummy`
	payload := strings.NewReader(`{"query":"query minichefs($where: Minichef_filter) {\n  minichefs(where: $where) {\n    id\n    totalAllocPoint\n    rewardPerSecond\n    rewardsExpiration\n    farms(first: 1000) {\n      id\n      pid\n      tvl\n      allocPoint\n      rewarderAddress\n      chefAddress\n      pairAddress\n      rewarder {\n        id\n        rewards {\n          id\n          token {\n            id\n            symbol\n            derivedUSD\n            name\n            decimals\n          }\n          multiplier\n        }\n      }\n      pair {\n        id\n        reserve0\n        reserve1\n        totalSupply\n        token0 {\n          id\n          symbol\n          derivedUSD\n          name\n          decimals\n        }\n        token1 {\n          id\n          symbol\n          derivedUSD\n          name\n          decimals\n        }\n      }\n   }\n  }\n}\n","operationName":"minichefs"}`)
	r, _ := req.Post(url, payload)
	err = r.ToJSON(&v)
	return
}
