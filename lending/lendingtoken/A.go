package lendingtoken

import (
	"errors"
	"strings"

	"github.com/0xVanfer/erc"
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/model"
)

type AToken struct {
	ProtocolBasic   *model.ProtocolBasic
	Basic           *erc.ERC20Info // basic info of the token
	UnderlyingBasic *erc.ERC20Info // basic info of the underlying token
	ApyInfo         *model.ApyInfo // deposit apy info
}

// Use a token address to get underlying address.
func (t *AToken) GetUnderlyingAddress(atoken string) (string, error) {
	if !t.ProtocolBasic.Regularcheck() {
		return "", errors.New("a token protocol basic must be initialized")
	}
	aList := *ethaddr.AaveLikeATokenListMap[t.ProtocolBasic.ProtocolName]
	for underlying, atokenAddress := range aList[t.ProtocolBasic.Network] {
		if strings.EqualFold(atokenAddress, atoken) {
			return underlying, nil
		}
	}
	return "", errors.New("underlying token not found by a token " + atoken)
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
