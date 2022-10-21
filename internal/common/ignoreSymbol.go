package common

import (
	"github.com/0xVanfer/ethprotocol/internal/constants"
	"github.com/0xVanfer/utils"
)

func ContainIgnoreSymbol(symbols ...string) bool {
	for _, symbol := range symbols {
		if utils.ContainInArrayX(symbol, constants.IgnoreSymbols) {
			return true
		}
	}
	return false
}
