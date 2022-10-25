package ethprotocol

import (
	"errors"

	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/lending"
)

// Update lending pool tokens.
func (prot *Protocol) UpdateLendingPoolTokens() error {
	// protocol basic must be initialized
	if !prot.ProtocolBasic.Regularcheck() {
		return errors.New("protocol basic must be initialized")
	}
	tokenLists := ethaddr.LendingTokenListsMap[prot.ProtocolBasic.ProtocolName]
	if !tokenLists.RegularCheck() {
		return errors.New("protocol not supported lending")
	}

	// use a token map or c token map to get the underlying list
	var tokenListMap map[string]map[string]string

	if tokenLists.ATokenList != nil {
		tokenListMap = *ethaddr.AaveLikeATokenListMap[prot.ProtocolBasic.ProtocolName]
	} else if tokenLists.CTokenList != nil {
		tokenListMap = *ethaddr.CompoundLikeCTokenListMap[prot.ProtocolBasic.ProtocolName]
	} else {
		return errors.New("protocol not supported")
	}
	// use underlying to update tokens
	for underlying := range tokenListMap[prot.ProtocolBasic.Network] {
		var newPool lending.LendingPool
		// will not return err
		_ = newPool.Init(*prot.ProtocolBasic)
		_ = newPool.UpdateTokensByUnderlying(underlying)
		prot.LendingPools = append(prot.LendingPools, &newPool)
	}
	return nil
}

// Update some of the protocol's lend pools apys by given underlying addresses.
//
// If "underlyings" is empty, update all pools.
func (prot *Protocol) UpdateLending(underlyings ...string) error {
	// protocol basic must be initialized
	if !prot.ProtocolBasic.Regularcheck() {
		return errors.New("protocol basic must be initialized")
	}
	switch prot.ProtocolBasic.ProtocolName {
	// aave v2
	case ethaddr.AaveV2Protocol:
		return prot.updateLendingAaveV2(underlyings)
	// aave v3
	case ethaddr.AaveV3Protocol:
		return prot.updateLendingAaveV3(underlyings)
	// benqi
	case ethaddr.BenqiProtocol:
		return prot.updateLendingBenqi(underlyings)
	// tradejoe
	case ethaddr.TraderJoeProtocol:
		return prot.updateLendingTraderjoe(underlyings)
	default:
		return errors.New(prot.ProtocolBasic.ProtocolName + " lend pools not supported")
	}
}
