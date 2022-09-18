package lendcompoundlike

import (
	"errors"
	"strings"
)

// Use c token address to get underlying address.
func (t *CToken) GetUnderlyingAddress(ctoken string) (string, error) {
	if !t.ProtocolBasic.Regularcheck() {
		return "", errors.New("c token protocol basic must be initialized")
	}
	cList := CTokenListMap[t.ProtocolBasic.ProtocolName]
	for underlying, ctokenAddress := range cList[t.ProtocolBasic.Network] {
		if strings.EqualFold(ctokenAddress, ctoken) {
			return underlying, nil
		}
	}
	return "", errors.New("underlying token not found by c token " + ctoken)
}

// Use underlying address to update c token info.
func (t *CToken) UpdateCTokenByUnderlying(underlying string) error {
	if !t.ProtocolBasic.Regularcheck() {
		return errors.New("c token protocol basic must be initialized")
	}
	ctoken := CTokenListMap[t.ProtocolBasic.ProtocolName][t.ProtocolBasic.Network][underlying]
	err := t.Basic.Init(ctoken, t.ProtocolBasic.Network, *t.ProtocolBasic.Client)
	if err != nil {
		return err
	}
	return nil
}
