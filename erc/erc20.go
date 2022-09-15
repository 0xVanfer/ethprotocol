package erc

import (
	"errors"
	"math"
	"strings"

	"github.com/0xVanfer/abigen/erc20"
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// Basic info of ERC20 token.
type ERC20 struct {
	Network     string
	Address     string
	Symbol      string
	Decimals    int
	TotalSupply float64
	// Contract    *erc20.Erc20
}

func (t *ERC20) Init(address string, network string, client bind.ContractBackend) error {
	if strings.EqualFold(types.ToLowerString(address), ethaddr.ZEROAddress) {
		return t.Init(ethaddr.WrappedChainTokenList[network], network, client)
	}
	if len(address) != 42 {
		return errors.New("address length must be 42, and not 0x0")
	}
	token, err := erc20.NewErc20(types.ToAddress(address), client)
	if err != nil {
		return err
	}
	decimals, err := token.Decimals(nil)
	if err != nil {
		return err
	}
	symbol, err := token.Symbol(nil)
	if err != nil {
		return err
	}
	totalSupply, err := token.TotalSupply(nil)
	if err != nil {
		return err
	}

	t = &ERC20{
		Network:     network,
		Address:     address,
		Symbol:      symbol,
		Decimals:    int(decimals.Int64()),
		TotalSupply: types.ToFloat64(totalSupply) * math.Pow10(-int(decimals.Int64())),
		// Contract:    token,
	}
	return nil
}
