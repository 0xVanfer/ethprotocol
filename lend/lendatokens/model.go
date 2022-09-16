package lendatokens

import (
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/apy"
	"github.com/0xVanfer/ethprotocol/erc"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type AToken struct {
	ProtocolName    string
	Network         string
	Basic           erc.ERC20   // basic info of the token
	UnderlyingBasic erc.ERC20   // basic info of the underlying token
	ApyInfo         apy.ApyInfo // deposit apy info
}

// Initialize the AToken by using atoken address.
// Return: if found(bool); err(error).
func (t *AToken) InitWithAToken(aAddress string, network string, protocolName string, client bind.ContractBackend) (bool, error) {
	var underlying string
	switch protocolName {
	// aavev2
	case ethaddr.AaveV2Protocol:
		underlying = GetAaveV2ATokenUnderlying(aAddress, network, client)
	// aavev3
	case ethaddr.AaveV3Protocol:
		underlying = GetAaveV3ATokenUnderlying(aAddress, network, client)
	}
	// not found in ethaddr
	if underlying == "" {
		return false, nil
	}
	// set t value
	t.ProtocolName = protocolName
	t.Network = network
	err := t.Basic.Init(aAddress, network, client)
	if err != nil {
		return false, err
	}
	err = t.UnderlyingBasic.Init(underlying, network, client)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Initialize the AToken by using underlying address.
// Return: if found(bool); err(error).
func (t *AToken) InitWithUnderlying(underlying string, network string, protocolName string, client bind.ContractBackend) (bool, error) {
	var aAddress string
	switch protocolName {
	// aavev2
	case ethaddr.AaveV2Protocol:
		aAddress = ethaddr.AaveATokenV2List[network][underlying]
	// aavev3
	case ethaddr.AaveV3Protocol:
		aAddress = ethaddr.AaveATokenV3List[network][underlying]
	}
	// not found in ethaddr
	if aAddress == "" {
		return false, nil
	}
	// set t value
	t.ProtocolName = protocolName
	t.Network = network
	err := t.Basic.Init(aAddress, network, client)
	if err != nil {
		return false, err
	}
	err = t.UnderlyingBasic.Init(underlying, network, client)
	if err != nil {
		return false, err
	}
	return true, nil
}
