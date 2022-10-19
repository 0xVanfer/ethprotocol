package ethprotocol

import (
	"fmt"

	"github.com/0xVanfer/coingecko"
	"github.com/0xVanfer/ethprotocol/model"
)

// Create a new protocol, with basic info and empty pools.
func New(input ProtocolInput) (*Protocol, error) {
	if input.Client == nil {
		fmt.Println("You do not have a contract backend, most functions will not work properly.")
	}
	var gecko *coingecko.Gecko

	if input.Coingecko.TokenList != nil {
		// already initiated
		gecko = &input.Coingecko
	} else {
		//  use key to initiate
		if input.Coingecko.ApiKey == "" {
			fmt.Println("You do not have a coingecko api key, some calculations related to token price will not work properly.")
		}
		var err error
		gecko, err = coingecko.New(input.Coingecko.ApiKey)
		if err != nil {
			return nil, err
		}
	}
	ProtocolBasic := model.ProtocolBasic{
		Network:      input.Network,
		ProtocolName: input.Name,
		Client:       &input.Client,
		Gecko:        gecko,
	}
	protocol := Protocol{
		ProtocolBasic: &ProtocolBasic,
	}
	// lending tokens
	err := protocol.UpdateLendingPoolTokens()
	if err != nil {
		return &protocol, err
	}
	return &protocol, nil
}
