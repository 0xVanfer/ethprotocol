package lendvtokens

import (
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/erc"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type VToken struct {
	ProtocolName    string
	Network         string
	Basic           erc.ERC20 // basic info of the token
	UnderlyingBasic erc.ERC20 // basic info of the underlying token
}

// Initialize the VToken by using vtoken address.
// Return: if found(bool); err(error).
func (t *VToken) InitWithVToken(vAddress string, network string, protocolName string, client bind.ContractBackend) (bool, error) {
	var underlying string
	switch protocolName {
	// aavev2
	case ethaddr.AaveV2Protocol:
		underlying = GetAaveV2VTokenUnderlying(vAddress, network, client)
	// aavev3
	case ethaddr.AaveV3Protocol:
		underlying = GetAaveV3VTokenUnderlying(vAddress, network, client)
	}
	// not found in ethaddr
	if underlying == "" {
		return false, nil
	}
	// set t value
	t.ProtocolName = protocolName
	t.Network = network
	err := t.Basic.Init(vAddress, network, client)
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
func (t *VToken) InitWithUnderlying(underlying string, network string, protocolName string, client bind.ContractBackend) (bool, error) {
	var vAddress string
	switch protocolName {
	// aavev2
	case ethaddr.AaveV2Protocol:
		vAddress = ethaddr.AaveVTokenV2List[network][underlying]
	// aavev3
	case ethaddr.AaveV3Protocol:
		vAddress = ethaddr.AaveVTokenV3List[network][underlying]
	}
	// not found in ethaddr
	if vAddress == "" {
		return false, nil
	}
	// set t value
	t.ProtocolName = protocolName
	t.Network = network
	err := t.Basic.Init(vAddress, network, client)
	if err != nil {
		return false, err
	}
	err = t.UnderlyingBasic.Init(underlying, network, client)
	if err != nil {
		return false, err
	}
	return true, nil
}
