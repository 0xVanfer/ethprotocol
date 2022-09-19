package ethprotocol

import (
	"fmt"

	"github.com/0xVanfer/coingecko"
	"github.com/0xVanfer/ethprotocol/lend"
	"github.com/0xVanfer/ethprotocol/liquidity"
	"github.com/0xVanfer/ethprotocol/model"
	"github.com/0xVanfer/ethprotocol/stake"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type Protocol struct {
	ProtocolBasic  model.ProtocolBasic
	LiquidityPools []liquidity.LiquidityPool
	StakePools     []stake.StakePool
	LendPools      []lend.LendPool
}

func New(network string, protocolName string, client bind.ContractBackend, coingeckoApiKey string) (*Protocol, error) {
	if client == nil {
		fmt.Println("You do not have a client, most functions will not work properly.")
	}
	if coingeckoApiKey == "" {
		fmt.Println("You do not have a coingecko api key, some apy calculations will not work properly.")
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
		ProtocolBasic: ProtocolBasic,
	}
	return &protocol, nil
}
