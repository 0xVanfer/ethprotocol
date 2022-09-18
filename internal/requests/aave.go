package requests

import (
	"errors"
	"strings"

	"github.com/0xVanfer/chainId"
	"github.com/imroc/req"
)

type aaveV2LendPoolsReq struct {
	Data struct {
		Reserves []aaveV2LendPoolInfo `json:"reserves"`
	} `json:"data"`
}
type aaveAToken struct {
	ID string `json:"id"`
}

type aaveV2LendPoolInfo struct {
	Symbol             string     `json:"symbol"`
	Decimals           int        `json:"decimals"`
	UnderlyingAsset    string     `json:"underlyingAsset"`
	AToken             aaveAToken `json:"aToken"`
	VToken             aaveAToken `json:"vToken"`
	SToken             aaveAToken `json:"sToken"`
	LiquidityRate      string     `json:"liquidityRate"`
	StableBorrowRate   string     `json:"stableBorrowRate"`
	VariableBorrowRate string     `json:"variableBorrowRate"`
	AEmissionPerSecond string     `json:"aEmissionPerSecond"`
	VEmissionPerSecond string     `json:"vEmissionPerSecond"`
	SEmissionPerSecond string     `json:"sEmissionPerSecond"`
}

// Return aava v2 lend pools info.
func ReqAaveV2LendPools(network string) ([]aaveV2LendPoolInfo, error) {
	var url string
	switch network {
	case chainId.AvalancheChainName:
		url = "https://api.thegraph.com/subgraphs/name/aave/protocol-v2-avalanche"
	case chainId.EthereumChainName:
		url = "https://api.thegraph.com/subgraphs/name/aave/protocol-v2"
	default:
		errString := "not supported network:" + network
		err := errors.New(errString)
		return nil, err
	}
	payload := strings.NewReader(`{"query":"{\n reserves {\n symbol\n decimals\n underlyingAsset\n aToken{id}\n vToken{id}\n sToken{id}\n liquidityRate\n stableBorrowRate\n variableBorrowRate\n aEmissionPerSecond\n vEmissionPerSecond\n sEmissionPerSecond\n}\n}\n\n"}`)
	r, _ := req.Post(url, payload)
	var v aaveV2LendPoolsReq
	err := r.ToJSON(&v)
	return v.Data.Reserves, err
}
