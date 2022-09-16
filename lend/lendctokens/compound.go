package lendctokens

import (
	"strings"

	"github.com/0xVanfer/abigen/compound/compoundCToken"
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// Read compound ctokens from ethaddr, if not exist, try reading from contracts.
func GetCompoundCTokenUnderlying(address string, network string, client bind.ContractBackend) string {
	for underlying, ctoken := range ethaddr.CompoundCTokenList[network] {
		if strings.EqualFold(address, ctoken) {
			return underlying
		}
	}
	cToken, err := compoundCToken.NewCompoundCToken(types.ToAddress(address), client)
	if err != nil {
		return ""
	}
	underlying, err := cToken.Underlying(nil)
	underlyingStr := types.ToLowerString(underlying)
	if (err != nil) || (underlyingStr == ethaddr.ZEROAddress) {
		return ""
	}
	return underlyingStr
}
