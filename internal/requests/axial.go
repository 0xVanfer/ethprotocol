package requests

import (
	"time"

	"github.com/imroc/req"
)

type axialToken struct {
	ID          string    `json:"id"`           // ignore
	Name        string    `json:"name"`         // token name
	Symbol      string    `json:"symbol"`       // token symbol
	Address     string    `json:"address"`      // token address
	Decimals    string    `json:"decimals"`     // token decimals
	TotalSupply string    `json:"total_supply"` // 0, should not be used
	MaxSupply   string    `json:"max_supply"`   // 0, should not be used
	CircSupply  string    `json:"circ_supply"`  // 0, should not be used
	Modified    time.Time `json:"modified"`     // last modified
	Created     time.Time `json:"created"`      // created
	TokenIndex  string    `json:"token_index"`  // token index in this pool
}

type axialRewardToken struct {
	Address    string `json:"address"`     // token address
	AvgApr     string `json:"avg_apr"`     // average apr
	BoostedApr string `json:"boosted_apr"` // boosted apr
	Decimals   string `json:"decimals"`    // token decimals
	Symbol     string `json:"symbol"`      // token symbol
}

type AxialPool struct {
	ID                    string             `json:"id"`                       // ignore
	Symbol                string             `json:"symbol"`                   // lp token symvol
	TokenAddress          string             `json:"tokenaddress"`             // lp address
	SwapAddress           string             `json:"swapaddress"`              // router address
	LastApr               string             `json:"last_apr"`                 // last apr
	LastVol               string             `json:"last_vol"`                 // last volume
	Deprecated            bool               `json:"deprecated"`               // if deprecated
	Modified              time.Time          `json:"modified"`                 // last modified
	Created               time.Time          `json:"created"`                  // last created
	LastRewardsApr        [][]string         `json:"last_rewards_apr"`         // last reward apr info
	LastPercGaugeAlloc    string             `json:"last_perc_gauge_alloc"`    //
	LastDailyAxialAlloc   string             `json:"last_daily_axial_alloc"`   //
	LastGaugeAxialBalance string             `json:"last_gauge_axial_balance"` //
	LastGaugeWeight       string             `json:"last_gauge_weight"`        //
	LastTvl               string             `json:"last_tvl"`                 // last reserve usd
	LastTokenPrice        string             `json:"last_token_price"`         // last token price
	LastSwapApr           string             `json:"last_swap_apr"`            // last swap apr
	GaugeAddress          string             `json:"gauge_address"`            // gauge address
	MedianBoost           string             `json:"median_boost"`             // median boost
	Metapool              bool               `json:"metapool"`                 // if is metapool
	Tokens                []axialToken       `json:"tokens"`                   // tokens
	RewardTokens          []axialRewardToken `json:"rewardTokens"`             // reward token
}

func ReqAxialAvaxPools() (v []AxialPool, err error) {
	url := "https://axial-api.snowapi.net/pools"
	r, _ := req.Get(url)
	err = r.ToJSON(&v)
	return
}
