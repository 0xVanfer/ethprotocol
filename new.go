package ethprotocol

import (
	"fmt"

	"github.com/0xVanfer/coingecko"
	"github.com/0xVanfer/ethprotocol/model"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// Create a new protocol, with basic info and empty pools.
//
// "client" and "coingeckoApiKey" can be empty, but some functions may fail.
func New(network string, protocolName string, client bind.ContractBackend, coingeckoApiKey string) (*Protocol, error) {
	if client == nil {
		fmt.Println("You do not have a contract backend, most functions will not work properly.")
	}
	if coingeckoApiKey == "" {
		fmt.Println("You do not have a coingecko api key, some calculations related to token price will not work properly.")
	}
	gecko, err := coingecko.New(coingeckoApiKey)
	if err != nil {
		return nil, err
	}
	ProtocolBasic := model.ProtocolBasic{
		Network:      network,
		ProtocolName: protocolName,
		Client:       &client,
		Gecko:        gecko,
	}
	protocol := Protocol{
		ProtocolBasic: &ProtocolBasic,
	}
	return &protocol, nil
}
