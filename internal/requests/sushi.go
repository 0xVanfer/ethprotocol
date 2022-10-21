package requests

import (
	"github.com/0xVanfer/chainId"
	"github.com/0xVanfer/types"
	"github.com/imroc/req"
)

type SushiPairInfo struct {
	Pair struct {
		Token0 struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Symbol      string `json:"symbol"`
			Decimals    string `json:"decimals"`
			TotalSupply string `json:"totalSupply"`
			DerivedETH  string `json:"derivedETH"`
		} `json:"token0"`
		Token1 struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Symbol      string `json:"symbol"`
			Decimals    string `json:"decimals"`
			TotalSupply string `json:"totalSupply"`
			DerivedETH  string `json:"derivedETH"`
		} `json:"token1"`
		ID string `json:"id"`
	} `json:"pair"`
	Liquidity                string  `json:"liquidity"`                // total reserve usd
	Liquidity1D              string  `json:"liquidity1d"`              //
	Liquidity1DChange        float64 `json:"liquidity1dChange"`        //
	Liquidity1DChangePercent float64 `json:"liquidity1dChangePercent"` //
	Volume1D                 float64 `json:"volume1d"`                 // volume 24
	Volume1DChange           float64 `json:"volume1dChange"`           //
	Volume1DChangePercent    float64 `json:"volume1dChangePercent"`    //
	Volume1W                 float64 `json:"volume1w"`                 //
	Fees1D                   float64 `json:"fees1d"`                   // fees collected
	Fees1W                   float64 `json:"fees1w"`                   //
	Utilisation1D            float64 `json:"utilisation1d"`            //
	Utilisation2D            float64 `json:"utilisation2d"`            //
	Utilisation1DChange      float64 `json:"utilisation1dChange"`      //
	Tx1D                     int     `json:"tx1d"`                     //
	Tx2D                     int     `json:"tx2d"`                     //
	Tx1DChange               float64 `json:"tx1dChange"`               //
	Apy                      string  `json:"apy"`                      // apy
}

func ReqSushiPairs(network string) (v []SushiPairInfo, err error) {
	chainNumber := chainId.ChainName2IdMap[network]
	url := "https://app.sushi.com/api/analytics/pairs/" + types.ToLowerString(chainNumber)
	r, _ := req.Get(url)
	err = r.ToJSON(&v)
	return
}
