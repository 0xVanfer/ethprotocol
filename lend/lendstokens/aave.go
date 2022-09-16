package lendstokens

import (
	"strings"

	"github.com/0xVanfer/ethaddr"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// Read aave v2 stokens from ethaddr.
func GetAaveV2STokenUnderlying(address string, network string, client bind.ContractBackend) string {
	for underlying, stoken := range ethaddr.AaveSTokenV2List[network] {
		if strings.EqualFold(address, stoken) {
			return underlying
		}
	}
	return ""
}

// Read aave v3 vtokens from ethaddr.
func GetAaveV3STokenUnderlying(address string, network string, client bind.ContractBackend) string {
	for underlying, stoken := range ethaddr.AaveSTokenV3List[network] {
		if strings.EqualFold(address, stoken) {
			return underlying
		}
	}
	return ""
}
