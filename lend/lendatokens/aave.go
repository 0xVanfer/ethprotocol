package lendatokens

import (
	"strings"

	"github.com/0xVanfer/abigen/aave/aaveATokenV2"
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// Read aave v2 atokens from ethaddr, if not exist, try reading from contracts.
func GetAaveV2ATokenUnderlying(address string, network string, client bind.ContractBackend) string {
	for underlying, atoken := range ethaddr.AaveATokenV2List[network] {
		if strings.EqualFold(address, atoken) {
			return underlying
		}
	}
	atoken, err := aaveATokenV2.NewAaveATokenV2(types.ToAddress(address), client)
	if err != nil {
		return ""
	}
	underlying, err := atoken.UNDERLYINGASSETADDRESS(nil)
	underlyingStr := types.ToLowerString(underlying)
	if (err != nil) || (underlyingStr == ethaddr.ZEROAddress) {
		return ""
	}
	return underlyingStr
}

// Read aave v3 atokens from ethaddr, if not exist, try reading from contracts.
func GetAaveV3ATokenUnderlying(address string, network string, client bind.ContractBackend) string {
	for underlying, atoken := range ethaddr.AaveATokenV3List[network] {
		if strings.EqualFold(address, atoken) {
			return underlying
		}
	}
	// v3 and v2 use the same function "UNDERLYINGASSETADDRESS"
	atoken, err := aaveATokenV2.NewAaveATokenV2(types.ToAddress(address), client)
	if err != nil {
		return ""
	}
	underlying, err := atoken.UNDERLYINGASSETADDRESS(nil)
	underlyingStr := types.ToLowerString(underlying)
	if (err != nil) || (underlyingStr == ethaddr.ZEROAddress) {
		return ""
	}
	return underlyingStr
}
