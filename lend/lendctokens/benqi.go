package lendctokens

import (
	"strings"

	"github.com/0xVanfer/abigen/benqi/benqiCToken"
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// Read benqi ctokens from ethaddr, if not exist, try reading from contracts.
func GetBenqiCTokenUnderlying(address string, network string, client bind.ContractBackend) string {
	for underlying, ctoken := range ethaddr.BenqiCTokenList[network] {
		if strings.EqualFold(address, ctoken) {
			return underlying
		}
	}
	qiToken, err := benqiCToken.NewBenqiCToken(types.ToAddress(address), client)
	if err != nil {
		return ""
	}
	underlying, err := qiToken.Underlying(nil)
	underlyingStr := types.ToLowerString(underlying)
	if (err != nil) || (underlyingStr == ethaddr.ZEROAddress) {
		return ""
	}
	return underlyingStr
}
