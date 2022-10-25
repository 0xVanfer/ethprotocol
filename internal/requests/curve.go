package requests

import (
	"github.com/imroc/req"
	"github.com/shopspring/decimal"
)

type curveAvaxPoolMainListReq struct {
	Success bool `json:"success"`
	Data    struct {
		PoolData        []curvePoolData `json:"poolData"`        // pool data info
		TvlAll          decimal.Decimal `json:"tvlAll"`          // total main pool tvl
		Tvl             decimal.Decimal `json:"tvl"`             // same as tvl all, ethereum has this
		GeneratedTimeMs int64           `json:"generatedTimeMs"` //
	} `json:"data"`
}

type curvePoolData struct {
	ID                        string           `json:"id"`                        // sort of pid
	Address                   string           `json:"address"`                   // router address
	CoinsAddresses            []string         `json:"coinsAddresses"`            // tokens to make up the lp
	Decimals                  []string         `json:"decimals"`                  // lp token decimals
	UnderlyingDecimals        []string         `json:"underlyingDecimals"`        // underlying decimals
	AssetType                 string           `json:"assetType"`                 // 0
	TotalSupply               interface{}      `json:"totalSupply"`               // total supply, big int
	LpTokenAddress            string           `json:"lpTokenAddress,omitempty"`  // lp address
	Name                      string           `json:"name,omitempty"`            // pool name
	Symbol                    string           `json:"symbol,omitempty"`          // lp symbol
	PriceOracle               interface{}      `json:"priceOracle,omitempty"`     // oracle address
	Implementation            string           `json:"implementation"`            // ""
	AssetTypeName             string           `json:"assetTypeName"`             // "usd"
	Coins                     []curvePoolCoins `json:"coins"`                     // tokens info
	UsdTotal                  decimal.Decimal  `json:"usdTotal"`                  // total usd
	GaugeAddress              string           `json:"gaugeAddress"`              // gauge address, ethereum has this
	UsdTotalExcludingBasePool decimal.Decimal  `json:"usdTotalExcludingBasePool"` //
	AmplificationCoefficient  string           `json:"amplificationCoefficient"`  //
}

type curvePoolCoins struct {
	Address           string          `json:"address"`           // token address
	UsdPrice          decimal.Decimal `json:"usdPrice"`          // token price in usd
	Decimals          interface{}     `json:"decimals"`          // token decimals
	IsBasePoolLpToken bool            `json:"isBasePoolLpToken"` // whether it is base pool lp
	Symbol            string          `json:"symbol"`            // token symbol
	PoolBalance       string          `json:"poolBalance"`       // reserve amount(not usd)
}

// Curve main pools.
func ReqCurvePoolMain(network string) (v curveAvaxPoolMainListReq, err error) {
	url := "https://api.curve.fi/api/getPools/" + network + "/main"
	r, _ := req.Get(url)
	err = r.ToJSON(&v)
	return
}

// Curve crypto pools.
func ReqCurveCrypto(network string) (v curveAvaxPoolMainListReq, err error) {
	url := "https://api.curve.fi/api/getPools/" + network + "/crypto"
	r, _ := req.Get(url)
	err = r.ToJSON(&v)
	return
}

type curvePoolVolumeAprReq struct {
	Success bool `json:"success"`
	Data    struct {
		PoolList        []curvePoolVolumeApr `json:"poolList"`        // pools
		TotalVolume     float64              `json:"totalVolume"`     // total volume
		CryptoVolume    float64              `json:"cryptoVolume"`    // crypto volume
		CryptoShare     float64              `json:"cryptoShare"`     //
		GeneratedTimeMs int64                `json:"generatedTimeMs"` //
	} `json:"data"`
}

type curvePoolVolumeApr struct {
	Type            string      `json:"type"`            // type "main"/"stable-factory"/"crypto"
	Address         string      `json:"address"`         // router
	VolumeUSD       float64     `json:"volumeUSD"`       // volume in usd
	RawVolume       float64     `json:"rawVolume"`       // volume in amount
	LatestDailyApy  interface{} `json:"latestDailyApy"`  // apy
	LatestWeeklyApy interface{} `json:"latestWeeklyApy"` // apy
	VirtualPrice    interface{} `json:"virtualPrice"`    // virtual price
}

func ReqCurveVolumeApr(network string) (v curvePoolVolumeAprReq, err error) {
	url := "https://api.curve.fi/api/getSubgraphData/" + network
	r, _ := req.Get(url)
	err = r.ToJSON(&v)
	return
}

type curveAvaxStakeApyReq struct {
	Success bool `json:"success"`
	Data    struct {
		SideChainGaugesApys []curveSideChainGaugeApy `json:"sideChainGaugesApys"`
		GeneratedTimeMs     int64                    `json:"generatedTimeMs"`
	} `json:"data"`
}

type curveSideChainGaugeApy struct {
	Address                    string      `json:"address"`
	Name                       string      `json:"name"`
	Apy                        interface{} `json:"apy"`
	AreCrvRewardsStuckInBridge bool        `json:"areCrvRewardsStuckInBridge"`
}

func ReqCurveStakeApy(network string) (v curveAvaxStakeApyReq, err error) {
	url := "https://api.curve.fi/api/getFactoGaugesCrvRewards/" + network
	r, _ := req.Get(url)
	err = r.ToJSON(&v)
	return
}
