package ethprotocol

import (
	"fmt"

	"github.com/0xVanfer/ethprotocol/lend"
	"github.com/0xVanfer/ethprotocol/liquidity"
	"github.com/0xVanfer/ethprotocol/model"
	"github.com/0xVanfer/ethprotocol/stake"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type Protocol struct {
	ProtocolBasic  *model.ProtocolBasic
	LiquidityPools *[]liquidity.LiquidityPool
	StakePools     *[]stake.StakePool
	LendPools      *[]lend.LendPool
}

func New(network string, protocolName string, client *bind.ContractBackend) (*Protocol, error) {
	if client == nil {
		fmt.Println("You do not have a client, most functions can not be used.")
	}
	ProtocolBasic := model.ProtocolBasic{
		Network:      network,
		ProtocolName: protocolName,
		Client:       client,
	}
	protocol := Protocol{
		ProtocolBasic: &ProtocolBasic,
	}
	return &protocol, nil
}
