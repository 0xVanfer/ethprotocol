package lendctokens

import (
	"strings"

	"github.com/0xVanfer/abigen/traderjoe/traderjoeCToken"
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// Read traderjoe ctokens from ethaddr, if not exist, try reading from contracts.
func GetTraderjoeCTokenUnderlying(address string, network string, client bind.ContractBackend) string {
	for underlying, ctoken := range ethaddr.TraderjoeCTokenList[network] {
		if strings.EqualFold(address, ctoken) {
			return underlying
		}
	}
	jToken, err := traderjoeCToken.NewTraderjoeCToken(types.ToAddress(address), client)
	if err != nil {
		return ""
	}
	underlying, err := jToken.Underlying(nil)
	underlyingStr := types.ToLowerString(underlying)
	if (err != nil) || (underlyingStr == ethaddr.ZEROAddress) {
		return ""
	}
	return underlyingStr
}
