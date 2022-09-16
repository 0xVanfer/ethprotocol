package model

import "github.com/ethereum/go-ethereum/accounts/abi/bind"

type ProtocolBase struct {
	Network      string
	ProtocolName string
	Client       bind.ContractBackend
}
