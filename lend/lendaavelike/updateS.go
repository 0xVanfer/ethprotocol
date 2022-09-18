package lendaavelike

import (
	"errors"
	"strings"
)

// Use s token address to get underlying address.
func (t *SToken) GetUnderlyingAddress(stoken string) (string, error) {
	if !t.ProtocolBasic.Regularcheck() {
		return "", errors.New("s token protocol basic must be initialized")
	}
	sList := STokenListMap[t.ProtocolBasic.ProtocolName]
	for underlying, stokenAddress := range sList[t.ProtocolBasic.Network] {
		if strings.EqualFold(stokenAddress, stoken) {
			return underlying, nil
		}
	}
	return "", errors.New("underlying token not found by s token " + stoken)
}

// Use underlying address to update s token info.
func (t *SToken) UpdateSTokenByUnderlying(underlying string) error {
	if !t.ProtocolBasic.Regularcheck() {
		return errors.New("s token protocol basic must be initialized")
	}
	stoken := STokenListMap[t.ProtocolBasic.ProtocolName][t.ProtocolBasic.Network][underlying]
	err := t.Basic.Init(stoken, t.ProtocolBasic.Network, *t.ProtocolBasic.Client)
	if err != nil {
		return err
	}
	return nil
}
