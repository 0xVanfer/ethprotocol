package lend

import (
	"strings"

	"github.com/0xVanfer/abigen/benqi/benqiCToken"
	"github.com/0xVanfer/abigen/compound/compoundCToken"
	"github.com/0xVanfer/abigen/traderjoe/traderjoeCToken"
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

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
