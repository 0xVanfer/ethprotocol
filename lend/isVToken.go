package lend

import (
	"strings"

	"github.com/0xVanfer/ethaddr"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func GetAaveV2VTokenUnderlying(address string, network string, client bind.ContractBackend) string {
	for underlying, vtoken := range ethaddr.AaveVTokenV2List[network] {
		if strings.EqualFold(address, vtoken) {
			return underlying
		}
	}
	return ""
}

func GetAaveV3VTokenUnderlying(address string, network string, client bind.ContractBackend) string {
	for underlying, vtoken := range ethaddr.AaveVTokenV3List[network] {
		if strings.EqualFold(address, vtoken) {
			return underlying
		}
	}
	return ""
}
