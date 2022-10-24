package requests

import (
	"github.com/0xVanfer/chainId"
	"github.com/0xVanfer/types"
	"github.com/imroc/req"
)

type sushiPair struct {
	Token0 sushiToken `json:"token0"` // token 0 info
	Token1 sushiToken `json:"token1"` // token 1 info
	ID     string     `json:"id"`     // lp address
}

type sushiToken struct {
	ID          string `json:"id"`          // token address
	Name        string `json:"name"`        // token name
	Symbol      string `json:"symbol"`      // token symbol
	Decimals    string `json:"decimals"`    // token decimals
	TotalSupply string `json:"totalSupply"` // token total supply
	DerivedETH  string `json:"derivedETH"`  // token price
}

type sushiPairInfo struct {
	Pair                     sushiPair `json:"pair"`
	Liquidity                string    `json:"liquidity"`                // total reserve usd
	Liquidity1D              string    `json:"liquidity1d"`              //
	Liquidity1DChange        float64   `json:"liquidity1dChange"`        //
	Liquidity1DChangePercent float64   `json:"liquidity1dChangePercent"` //
	Volume1D                 float64   `json:"volume1d"`                 // volume 24
	Volume1DChange           float64   `json:"volume1dChange"`           //
	Volume1DChangePercent    float64   `json:"volume1dChangePercent"`    //
	Volume1W                 float64   `json:"volume1w"`                 //
	Fees1D                   float64   `json:"fees1d"`                   // fees collected
	Fees1W                   float64   `json:"fees1w"`                   //
	Utilisation1D            float64   `json:"utilisation1d"`            //
	Utilisation2D            float64   `json:"utilisation2d"`            //
	Utilisation1DChange      float64   `json:"utilisation1dChange"`      //
	Tx1D                     int       `json:"tx1d"`                     //
	Tx2D                     int       `json:"tx2d"`                     //
	Tx1DChange               float64   `json:"tx1dChange"`               //
	Apy                      string    `json:"apy"`                      // apy
}

func ReqSushiPairs(network string) (v []sushiPairInfo, err error) {
	chainNumber := chainId.ChainName2IdMap[network]
	url := "https://app.sushi.com/api/analytics/pairs/" + types.ToLowerString(chainNumber)
	r, _ := req.Get(url)
	err = r.ToJSON(&v)
	return
}
