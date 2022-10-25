package lendingtoken

import (
	"errors"
	"strings"

	"github.com/0xVanfer/erc"
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/model"
)

// Aave like a token.
type AToken struct {
	ProtocolBasic   *model.ProtocolBasic // basic info of the protocol
	Basic           *erc.ERC20Info       // basic info of the token
	UnderlyingBasic *erc.ERC20Info       // basic info of the underlying token
	ApyInfo         *model.ApyInfo       // deposit apy info
}

func GetATokenUnderlyingAddress(protocol, network, atoken string) (string, error) {
	aList := ethaddr.AaveLikeATokenListMap[protocol]
	if aList == nil {
		return "", errors.New("not supported aave like protocol " + protocol)
	}
	for underlying, atokenAddress := range (*aList)[network] {
		if strings.EqualFold(atokenAddress, atoken) {
			return underlying, nil
		}
	}
	return "", errors.New("underlying token not found by a token " + atoken + " in " + protocol + ", " + network)
}

func (t *AToken) GetUnderlyingAddress() (string, error) {
	if !t.ProtocolBasic.Regularcheck() {
		return "", errors.New("a token protocol basic must be initialized")
	}
	return GetATokenUnderlyingAddress(t.ProtocolBasic.ProtocolName, t.ProtocolBasic.Network, *t.Basic.Address)
}

// Use underlying address to update a token info.
func (t *AToken) UpdateATokenByUnderlying(underlying string) error {
	if !t.ProtocolBasic.Regularcheck() {
		return errors.New("a token protocol basic must be initialized")
	}
	atoken := (*ethaddr.AaveLikeATokenListMap[t.ProtocolBasic.ProtocolName])[t.ProtocolBasic.Network][underlying]
	newBasic, err := erc.NewErc20(atoken, t.ProtocolBasic.Network, *t.ProtocolBasic.Client)
	if err != nil {
		return err
	}
	t.Basic = newBasic
	return nil
}
