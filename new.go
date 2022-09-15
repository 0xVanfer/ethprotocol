package ethprotocol

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type Protocol struct {
	Network      string
	ProtocolName string
	Client       bind.ContractBackend
}

func New(network string, protocolName string, client bind.ContractBackend) (*Protocol, error) {
	if client == nil {
		fmt.Println("You do not have a client, most functions can not be used.")
	}
	protocol := Protocol{
		Network:      network,
		ProtocolName: protocolName,
		Client:       client,
	}
	return &protocol, nil
}
