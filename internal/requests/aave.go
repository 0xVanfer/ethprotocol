package requests

import (
	"errors"
	"strings"

	"github.com/0xVanfer/chainId"
	"github.com/imroc/req"
)

type aaveV2LendingPoolsReq struct {
	Data struct {
		Reserves []aaveV2LendingPoolInfo `json:"reserves"`
	} `json:"data"`
}
type aaveAToken struct {
	ID string `json:"id"` // token address
}

type aaveV2LendingPoolInfo struct {
	Symbol                     string     `json:"symbol"`                     // underlying symbol
	Decimals                   int        `json:"decimals"`                   // underlying decimals
	UnderlyingAsset            string     `json:"underlyingAsset"`            // underlying address
	AToken                     aaveAToken `json:"aToken"`                     // atoken
	VToken                     aaveAToken `json:"vToken"`                     // vtoken
	SToken                     aaveAToken `json:"sToken"`                     // stoken
	LiquidityRate              string     `json:"liquidityRate"`              // in RAY. a basic reward rate
	StableBorrowRate           string     `json:"stableBorrowRate"`           // in RAY. s basic borrow rate
	VariableBorrowRate         string     `json:"variableBorrowRate"`         // in RAY. v basic borrow rate
	AEmissionPerSecond         string     `json:"aEmissionPerSecond"`         // in WEI. a incentive reward rate
	VEmissionPerSecond         string     `json:"vEmissionPerSecond"`         // in WEI. s incentive reward rate
	SEmissionPerSecond         string     `json:"sEmissionPerSecond"`         // in WEI. v incentive reward rate
	ReserveFactor              string     `json:"reserveFactor"`              // reserve factor
	TotalDeposits              string     `json:"totalDeposits"`              // total deposited assets(might not accurate)
	TotalLiquidity             string     `json:"totalLiquidity"`             // total avalable assets(might not accurate)
	TotalScaledVariableDebt    string     `json:"totalScaledVariableDebt"`    // total v debt used by aave
	TotalPrincipalStableDebt   string     `json:"totalPrincipalStableDebt"`   // total s debt
	TotalCurrentVariableDebt   string     `json:"totalCurrentVariableDebt"`   // total v debt
	TotalLiquidityAsCollateral string     `json:"totalLiquidityAsCollateral"` // total liquidity as collateral
	UtilizationRate            string     `json:"utilizationRate"`            // utilization of the pool token, equals total borrowed/total deposited
}

// Return aava v2 lend pools info.
//
// Only support avalanche and ethereum now.
func ReqAaveV2LendingPools(network string) ([]aaveV2LendingPoolInfo, error) {
	var url string
	switch network {
	case chainId.AvalancheChainName:
		url = "https://api.thegraph.com/subgraphs/name/aave/protocol-v2-avalanche"
	case chainId.EthereumChainName:
		url = "https://api.thegraph.com/subgraphs/name/aave/protocol-v2"
	default:
		return nil, errors.New("not supported network:" + network)
	}
	payload := strings.NewReader(`{"query":"{\n reserves {\n symbol\n decimals\n underlyingAsset\n aToken{id}\n vToken{id}\n sToken{id}\n liquidityRate\n stableBorrowRate\n variableBorrowRate\n aEmissionPerSecond\n vEmissionPerSecond\n sEmissionPerSecond\n totalDeposits\n totalLiquidity\n totalScaledVariableDebt\n totalPrincipalStableDebt\n reserveFactor\n totalCurrentVariableDebt\n totalLiquidityAsCollateral\n utilizationRate\n}\n}\n\n"}`)
	r, _ := req.Post(url, payload)
	var v aaveV2LendingPoolsReq
	err := r.ToJSON(&v)
	return v.Data.Reserves, err
}
