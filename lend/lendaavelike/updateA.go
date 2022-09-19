package lendaavelike

import (
	"errors"
	"strings"

	"github.com/0xVanfer/erc"
)

// Use a token address to get underlying address.
func (t *AToken) GetUnderlyingAddress(atoken string) (string, error) {
	if !t.ProtocolBasic.Regularcheck() {
		return "", errors.New("a token protocol basic must be initialized")
	}
	aList := ATokenListMap[t.ProtocolBasic.ProtocolName]
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
	atoken := ATokenListMap[t.ProtocolBasic.ProtocolName][t.ProtocolBasic.Network][underlying]
	var newBasic erc.ERC20Info
	err := newBasic.Init(atoken, t.ProtocolBasic.Network, *t.ProtocolBasic.Client)
	if err != nil {
		return err
	}
	t.Basic = &newBasic
	return nil
}
