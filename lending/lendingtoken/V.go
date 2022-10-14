package lendingtoken

import (
	"errors"
	"strings"

	"github.com/0xVanfer/erc"
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/model"
)

type VToken struct {
	ProtocolBasic   *model.ProtocolBasic
	Basic           *erc.ERC20Info // basic info of the token
	UnderlyingBasic *erc.ERC20Info // basic info of the underlying token
	ApyInfo         model.ApyInfo  // borrow variable apy info
}

// Use v token address to get underlying address.
func (t *VToken) GetUnderlyingAddress(vtoken string) (string, error) {
	if !t.ProtocolBasic.Regularcheck() {
		return "", errors.New("v token protocol basic must be initialized")
	}
	vList := *ethaddr.AaveLikeVTokenListMap[t.ProtocolBasic.ProtocolName]
	for underlying, vtokenAddress := range vList[t.ProtocolBasic.Network] {
		if strings.EqualFold(vtokenAddress, vtoken) {
			return underlying, nil
		}
	}
	return "", errors.New("underlying token not found by v token " + vtoken)
}

// Use underlying address to update v token info.
func (t *VToken) UpdateVTokenByUnderlying(underlying string) error {
	if !t.ProtocolBasic.Regularcheck() {
		return errors.New("v token protocol basic must be initialized")
	}
	vtoken := (*ethaddr.AaveLikeVTokenListMap[t.ProtocolBasic.ProtocolName])[t.ProtocolBasic.Network][underlying]
	var newBasic erc.ERC20Info
	err := newBasic.Init(vtoken, t.ProtocolBasic.Network, *t.ProtocolBasic.Client)
	if err != nil {
		return err
	}
	t.Basic = &newBasic
	return nil
}
