package lendaavelike

import (
	"errors"
	"strings"
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
	err := t.Basic.Init(atoken, t.ProtocolBasic.Network, *t.ProtocolBasic.Client)
	if err != nil {
		return err
	}
	return nil
}

// // Use a token address to update underlying and a token info.
// func (t *AToken) UpdateTokensByAToken(atoken string) error {
// 	if !t.ProtocolBasic.Regularcheck() {
// 		return errors.New("a token protocol basic must be initialized")
// 	}
// 	err := t.Basic.Init(atoken, t.ProtocolBasic.Network, *t.ProtocolBasic.Client)
// 	if err != nil {
// 		return err
// 	}
// 	underlyingAddress, err := t.GetUnderlyingAddress(atoken)
// 	if err != nil {
// 		return err
// 	}
// 	err = t.UnderlyingBasic.Init(underlyingAddress, t.ProtocolBasic.Network, *t.ProtocolBasic.Client)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
