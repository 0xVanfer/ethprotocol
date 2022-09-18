package model

import "github.com/ethereum/go-ethereum/accounts/abi/bind"

type ProtocolBasic struct {
	Network      string
	ProtocolName string
	Client       *bind.ContractBackend
}

// Return if protocol basic pass regular check.
func (p *ProtocolBasic) Regularcheck() bool {
	return (p.Network != "") && (p.ProtocolName != "")
}
