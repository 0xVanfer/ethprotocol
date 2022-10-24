package requests

import (
	"errors"
	"fmt"
	"strings"

	"github.com/0xVanfer/types"
	"github.com/0xVanfer/utils"
	"github.com/imroc/req"
)

type joeLiquidityPoolToken struct {
	ID        string `json:"id"`        // token address
	Symbol    string `json:"symbol"`    // token symbol
	Decimals  string `json:"decimals"`  // token decimals
	Volume    string `json:"volume"`    // token volume
	VolumeUSD string `json:"volumeUSD"` // token volume in usd
}

type joeLiquidityPoolHourData struct {
	VolumeUSD string `json:"volumeUSD"` // hour volume in usd
	Date      int    `json:"date"`      // timestamp
}

type joeLiquidityPool struct {
	ID          string                     `json:"id"`          // lp address
	Name        string                     `json:"name"`        // pool name
	ReserveUSD  string                     `json:"reserveUSD"`  // total reserve in usd
	TotalSupply string                     `json:"totalSupply"` // total supply, big int
	Token0      joeLiquidityPoolToken      `json:"token0"`      // token 0
	Token1      joeLiquidityPoolToken      `json:"token1"`      // token 1
	VolumeUSD   string                     `json:"volumeUSD"`   // total history volume
	HourData    []joeLiquidityPoolHourData `json:"hourData"`    // pool hour data
}

type joeLiquidityPoolReq struct {
	Data struct {
		Pairs []joeLiquidityPool `json:"pairs"`
	} `json:"data"`
}

func ReqJoeAvaxPools() ([]joeLiquidityPool, error) {
	url := "https://api.thegraph.com/subgraphs/name/traderjoe-xyz/exchange"
	timestamp := utils.TimestampNow()
	startTime := fmt.Sprint(timestamp - 86400*2)

	payload := strings.NewReader(`{
        "operationName": "pairsQuery",
        "variables": {"first": 300,"skip": 0,"orderBy": "reserveUSD","orderDirection": "desc"},
        "query": "query pairsQuery($first: Int! = 300, $skip: Int! = 0, $orderBy: String! = \"reserveUSD\", $orderDirection: String! = \"desc\", $dateAfter: Int! = ` + startTime + `) {\n  pairs(first: $first, skip: $skip, orderBy: $orderBy, orderDirection: $orderDirection) {\n    id\n\n    name\n    totalSupply\n    token0 {\n      id\n      symbol\n      decimals\n    volume\n    volumeUSD\n}\n    token1 {\n      id\n      symbol\n       decimals\n    volume\n    volumeUSD\n}\n    reserveUSD\n    hourData(first: 24, where: {date_gt: $dateAfter}, orderBy: date, orderDirection: desc) {\n      volumeUSD\n      date\n      }\n}\n}\n"
      }`)

	r, _ := req.Post(url, payload)
	var v joeLiquidityPoolReq
	err := r.ToJSON(&v)
	if err != nil {
		return nil, err
	}
	return v.Data.Pairs, nil
}

type joeStakeChefReq struct {
	Data joeStakeChefReqData `json:"data"`
}

type joeStakeChefReqData struct {
	MasterChefs []joeStakeMasterChefInfo `json:"masterChefs"`
}
type joeStakeMasterChefInfo struct {
	ID              string `json:"id"`
	TotalAllocPoint string `json:"totalAllocPoint"`
	JoePerSec       string `json:"joePerSec"`
}

// version = 2 or 3, representing chefV2 and chefV3
func ReqJoeAvaxStakeChefInfo[T types.Integer | string](version_ int) (info joeStakeMasterChefInfo, err error) {
	version := types.ToString(version_)
	if version != "2" && version != "3" {
		err = errors.New("traderjoe chef version not supported")
		return
	}
	url := `https://api.thegraph.com/subgraphs/name/traderjoe-xyz/masterchefv` + version
	payload := strings.NewReader(`{"query":"{\n masterChefs{\n id\n totalAllocPoint\n joePerSec\n}\n}\n"}`)
	r, _ := req.Post(url, payload)
	var v joeStakeChefReq
	err = r.ToJSON(&v)
	info = v.Data.MasterChefs[0]
	return
}

type joeStakePoolsReq struct {
	Data joeStakePoolsReqData `json:"data"`
}

type joeStakePoolsReqData struct {
	Pools []joeStakePoolsInfo `json:"pools"`
}
type joeStakePoolsInfo struct {
	ID         string      `json:"id"`
	Pair       string      `json:"pair"`
	AllocPoint string      `json:"allocPoint"`
	JlpBalance string      `json:"jlpBalance"`
	Rewarder   interface{} `json:"rewarder"`
}

// version = 2 or 3, representing chefV2 and chefV3
func ReqJoeAvaxStakePoolsInfo[T types.Integer | string](version_ int) (infos []joeStakePoolsInfo, err error) {
	version := types.ToString(version_)
	if version != "2" && version != "3" {
		err = errors.New("traderjoe chef version not supported")
		return
	}
	url := `https://api.thegraph.com/subgraphs/name/traderjoe-xyz/masterchefv` + version
	payload := strings.NewReader(`{"query":"{\n pools{\n id\n pair\n allocPoint\n jlpBalance\n rewarder{\n id\n rewardToken\n symbol\n tokenPerSec}}\n}"}`)
	r, _ := req.Post(url, payload)
	var v joeStakePoolsReq
	err = r.ToJSON(&v)
	infos = v.Data.Pools
	return
}
