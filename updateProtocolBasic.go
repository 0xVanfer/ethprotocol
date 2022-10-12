package ethprotocol

import (
	"time"

	"github.com/0xVanfer/coingecko"
)

// Update the protocol basic info.
//
// Will not update the empty params.
func (prot *Protocol) UpdateProtocolBasic(input ProtocolInput) {
	if input.Network != "" {
		prot.ProtocolBasic.Network = input.Network
	}
	if input.Name != "" {
		prot.ProtocolBasic.ProtocolName = input.Name
	}
	if input.Client != nil {
		prot.ProtocolBasic.Client = &input.Client
	}
	// update gecko
	if input.Coingecko.TokenList != nil {
		if input.Coingecko.ApiKey != "" {
			prot.ProtocolBasic.Gecko.ApiKey = input.Coingecko.ApiKey
		}
		prot.ProtocolBasic.Gecko.TokenList = input.Coingecko.TokenList
		prot.ProtocolBasic.Gecko.UpdatedAt = time.Now()
	} else if input.Coingecko.ApiKey != "" {
		newGecko, err := coingecko.New(input.Coingecko.ApiKey)
		if err != nil {
			return
		}
		prot.ProtocolBasic.Gecko = newGecko
	}
}
