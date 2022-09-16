package lendstokens

import (
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/erc"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type SToken struct {
	ProtocolName    string
	Network         string
	Basic           erc.ERC20 // basic info of the token
	UnderlyingBasic erc.ERC20 // basic info of the underlying token
}

// Initialize the SToken by using stoken address.
// Return: if found(bool); err(error).
func (t *SToken) InitWithSToken(sAddress string, network string, protocolName string, client bind.ContractBackend) (bool, error) {
	var underlying string
	switch protocolName {
	// aavev2
	case ethaddr.AaveV2Protocol:
		underlying = GetAaveV2STokenUnderlying(sAddress, network, client)
	// aavev3
	case ethaddr.AaveV3Protocol:
		underlying = GetAaveV3STokenUnderlying(sAddress, network, client)
	}
	// not found in ethaddr
	if underlying == "" {
		return false, nil
	}
	// set t value
	t.ProtocolName = protocolName
	t.Network = network
	err := t.Basic.Init(sAddress, network, client)
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
func (t *SToken) InitWithUnderlying(underlying string, network string, protocolName string, client bind.ContractBackend) (bool, error) {
	var sAddress string
	switch protocolName {
	// aavev2
	case ethaddr.AaveV2Protocol:
		sAddress = ethaddr.AaveSTokenV2List[network][underlying]
	// aavev3
	case ethaddr.AaveV3Protocol:
		sAddress = ethaddr.AaveSTokenV3List[network][underlying]
	}
	// not found in ethaddr
	if sAddress == "" {
		return false, nil
	}
	// set t value
	t.ProtocolName = protocolName
	t.Network = network
	err := t.Basic.Init(sAddress, network, client)
	if err != nil {
		return false, err
	}
	err = t.UnderlyingBasic.Init(underlying, network, client)
	if err != nil {
		return false, err
	}
	return true, nil
}
