package model

import (
	"github.com/0xVanfer/coingecko"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type ProtocolBasic struct {
	Network      string
	ProtocolName string
	Client       *bind.ContractBackend
	Gecko        *coingecko.Gecko
}

// Return if protocol basic pass regular check.
func (p *ProtocolBasic) Regularcheck() bool {
	if p == nil {
		return false
	}
	return (p.Network != "") && (p.ProtocolName != "")
}
