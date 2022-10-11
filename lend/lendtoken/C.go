package lendtoken

import (
	"errors"
	"strings"

	"github.com/0xVanfer/erc"
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/model"
)

type CToken struct {
	ProtocolBasic   *model.ProtocolBasic
	Basic           *erc.ERC20Info // basic info of the token
	UnderlyingBasic *erc.ERC20Info // basic info of the underlying token
	DepositApyInfo  model.ApyInfo  // deposit apy info
	BorrowApyInfo   model.ApyInfo  // deposit apy info
}

// Use c token address to get underlying address.
func (t *CToken) GetUnderlyingAddress(ctoken string) (string, error) {
	if !t.ProtocolBasic.Regularcheck() {
		return "", errors.New("c token protocol basic must be initialized")
	}
	cList := ethaddr.CompoundLikeCTokenListMap[t.ProtocolBasic.ProtocolName]
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
	ctoken := ethaddr.CompoundLikeCTokenListMap[t.ProtocolBasic.ProtocolName][t.ProtocolBasic.Network][underlying]
	var newBasic erc.ERC20Info
	err := newBasic.Init(ctoken, t.ProtocolBasic.Network, *t.ProtocolBasic.Client)
	if err != nil {
		return err
	}
	t.Basic = &newBasic
	return nil
}
