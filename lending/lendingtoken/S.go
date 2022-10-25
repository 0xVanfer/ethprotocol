package lendingtoken

import (
	"errors"
	"strings"

	"github.com/0xVanfer/erc"
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/model"
)

type SToken struct {
	ProtocolBasic   *model.ProtocolBasic
	Basic           *erc.ERC20Info // basic info of the token
	UnderlyingBasic *erc.ERC20Info // basic info of the underlying token
	ApyInfo         *model.ApyInfo // borrow stable apy info
}

// Use s token address to get underlying address.
func (t *SToken) GetUnderlyingAddress(stoken string) (string, error) {
	if !t.ProtocolBasic.Regularcheck() {
		return "", errors.New("s token protocol basic must be initialized")
	}
	sList := *ethaddr.AaveLikeSTokenListMap[t.ProtocolBasic.ProtocolName]
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
	stoken := (*ethaddr.AaveLikeSTokenListMap[t.ProtocolBasic.ProtocolName])[t.ProtocolBasic.Network][underlying]
	newBasic, err := erc.NewErc20(stoken, t.ProtocolBasic.Network, *t.ProtocolBasic.Client)
	if err != nil {
		return err
	}
	t.Basic = newBasic
	return nil
}
