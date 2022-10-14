package ethprotocol

import (
	"errors"

	"github.com/0xVanfer/ethaddr"
)

func (prot *Protocol) UpdateLendingPoolTokens() error {
	// protocol basic must be initialized
	if !prot.ProtocolBasic.Regularcheck() {
		return errors.New("protocol basic must be initialized")
	}
	tokenLists := ethaddr.LendingTokenListsMap[prot.ProtocolBasic.ProtocolName]
	if !tokenLists.RegularCheck() {
		return errors.New("protocol not supported")
	}
	if len(*tokenLists.ATokenList) > len(*tokenLists.CTokenList) {
		// use a token
		_ = 1
	} else {
		// use c token
		_ = 1
	}
	return nil
}
