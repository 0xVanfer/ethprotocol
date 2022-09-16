package lendctokens

import (
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/apy"
	"github.com/0xVanfer/ethprotocol/erc"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type CToken struct {
	ProtocolName    string
	Network         string
	Basic           erc.ERC20   // basic info of the token
	UnderlyingBasic erc.ERC20   // basic info of the underlying token
	DepositApyInfo  apy.ApyInfo // deposit apy info
	BorrowApyInfo   apy.ApyInfo // borrow apy info
}

// Initialize the CToken by using ctoken address.
// Return: if found(bool); err(error).
func (t *CToken) InitWithCToken(cAddress string, network string, protocolName string, client bind.ContractBackend) (bool, error) {
	var underlying string
	switch protocolName {
	// compound
	case ethaddr.CompoundProtocol:
		underlying = GetCompoundCTokenUnderlying(cAddress, network, client)
	// benqi
	case ethaddr.BenqiProtocol:
		underlying = GetBenqiCTokenUnderlying(cAddress, network, client)
		// traderjoe
	case ethaddr.TraderJoeProtocol:
		underlying = GetTraderjoeCTokenUnderlying(cAddress, network, client)
	}
	// not found in ethaddr
	if underlying == "" {
		return false, nil
	}
	// set t value
	t.ProtocolName = protocolName
	t.Network = network
	err := t.Basic.Init(cAddress, network, client)
	if err != nil {
		return false, err
	}
	err = t.UnderlyingBasic.Init(underlying, network, client)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Initialize the CToken by using underlying address.
// Return: if found(bool); err(error).
func (t *CToken) InitWithUnderlying(underlying string, network string, protocolName string, client bind.ContractBackend) (bool, error) {
	var cAddress string
	switch protocolName {
	// compound
	case ethaddr.CompoundProtocol:
		cAddress = ethaddr.CompoundCTokenList[network][underlying]
	// benqi
	case ethaddr.BenqiProtocol:
		cAddress = ethaddr.BenqiCTokenList[network][underlying]
		// traderjoe
	case ethaddr.TraderJoeProtocol:
		cAddress = ethaddr.TraderjoeCTokenList[network][underlying]
	}
	// not found in ethaddr
	if cAddress == "" {
		return false, nil
	}
	// set t value
	t.ProtocolName = protocolName
	t.Network = network
	err := t.Basic.Init(cAddress, network, client)
	if err != nil {
		return false, err
	}
	err = t.UnderlyingBasic.Init(underlying, network, client)
	if err != nil {
		return false, err
	}
	return true, nil
}
