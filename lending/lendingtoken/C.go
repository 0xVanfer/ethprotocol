package lendingtoken

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
	SupplyApyInfo   *model.ApyInfo // supply apy info
	BorrowApyInfo   *model.ApyInfo // borrow apy info
}

// Use c token address to get underlying address.
func (t *CToken) GetUnderlyingAddress(ctoken string) (string, error) {
	if !t.ProtocolBasic.Regularcheck() {
		return "", errors.New("c token protocol basic must be initialized")
	}
	cList := *ethaddr.CompoundLikeCTokenListMap[t.ProtocolBasic.ProtocolName]
	if cList == nil {
		return "", errors.New(t.ProtocolBasic.ProtocolName + " not supported on " + t.ProtocolBasic.Network)
	}
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
	ctoken := (*ethaddr.CompoundLikeCTokenListMap[t.ProtocolBasic.ProtocolName])[t.ProtocolBasic.Network][underlying]
	newBasic, err := erc.NewErc20(ctoken, t.ProtocolBasic.Network, *t.ProtocolBasic.Client)
	if err != nil {
		return err
	}
	t.Basic = newBasic
	return nil
}
