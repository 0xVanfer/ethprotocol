package lend

import (
	"errors"

	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/erc"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type AVSCToken struct {
	ProtocolName    string
	Network         string
	Basic           erc.ERC20 // basic info of the token
	UnderlyingBasic erc.ERC20 // basic info of the underlying token
	IsAToken        bool
	IsVToken        bool
	IsSToken        bool
	IsCToken        bool
}

func (t *AVSCToken) InitWithUnderlying(underlying string, network string, protocolName string, client bind.ContractBackend) error {
	// todo
	return nil
}

func (t *AVSCToken) InitWithAVSCToken(address string, network string, protocolName string, client bind.ContractBackend) error {
	t.Network = network
	t.ProtocolName = protocolName
	err := t.Basic.Init(address, network, client)
	if err != nil {
		return err
	}
	var underlying string
	switch protocolName {
	// aavev2
	case ethaddr.AaveV2Protocol:
		underlying = GetAaveV2ATokenUnderlying(address, network, client)
		if underlying != "" {
			t.IsAToken = true
			break
		}
		underlying = GetAaveV2VTokenUnderlying(address, network, client)
		if underlying != "" {
			t.IsVToken = true
			break
		}
	// aavev3
	case ethaddr.AaveV3Protocol:
		underlying = GetAaveV3ATokenUnderlying(address, network, client)
		if underlying != "" {
			t.IsAToken = true
			break
		}
		underlying = GetAaveV3VTokenUnderlying(address, network, client)
		if underlying != "" {
			t.IsVToken = true
			break
		}
	// benqi
	case ethaddr.BenqiProtocol:
		underlying = GetBenqiCTokenUnderlying(address, network, client)
		if underlying != "" {
			t.IsCToken = true
			break
		}
	// compound
	case ethaddr.CompoundProtocol:
		underlying = GetCompoundCTokenUnderlying(address, network, client)
		if underlying != "" {
			t.IsCToken = true
			break
		}
	// traderjoe
	case ethaddr.TraderJoeProtocol:
		underlying = GetTraderjoeCTokenUnderlying(address, network, client)
		if underlying != "" {
			t.IsCToken = true
			break
		}
	// network not defined, try every possibility
	case "":
		underlying = GetAaveV2ATokenUnderlying(address, network, client)
		if underlying != "" {
			t.IsAToken = true
			t.ProtocolName = ethaddr.AaveV2Protocol
			break
		}
		underlying = GetAaveV2VTokenUnderlying(address, network, client)
		if underlying != "" {
			t.IsVToken = true
			t.ProtocolName = ethaddr.AaveV2Protocol
			break
		}
		underlying = GetAaveV3ATokenUnderlying(address, network, client)
		if underlying != "" {
			t.IsAToken = true
			t.ProtocolName = ethaddr.AaveV3Protocol
			break
		}
		underlying = GetAaveV3VTokenUnderlying(address, network, client)
		if underlying != "" {
			t.IsVToken = true
			t.ProtocolName = ethaddr.AaveV3Protocol
			break
		}
		underlying = GetBenqiCTokenUnderlying(address, network, client)
		if underlying != "" {
			t.IsCToken = true
			t.ProtocolName = ethaddr.BenqiProtocol
			break
		}
		underlying = GetCompoundCTokenUnderlying(address, network, client)
		if underlying != "" {
			t.IsCToken = true
			t.ProtocolName = ethaddr.CompoundProtocol
			break
		}
		underlying = GetTraderjoeCTokenUnderlying(address, network, client)
		if underlying != "" {
			t.IsCToken = true
			t.ProtocolName = ethaddr.TraderJoeProtocol
			break
		}
	}
	if underlying == "" {
		return errors.New("init failed: not supported a/v/s/c token")
	}
	err = t.UnderlyingBasic.Init(underlying, network, client)
	if err != nil {
		return err
	}
	return nil
}
