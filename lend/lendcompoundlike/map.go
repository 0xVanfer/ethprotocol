package lendcompoundlike

import "github.com/0xVanfer/ethaddr"

// map[protocol name][network][underlying] = c token.
var CTokenListMap = map[string]map[string]map[string]string{
	ethaddr.CompoundProtocol:  ethaddr.CompoundCTokenList,
	ethaddr.BenqiProtocol:     ethaddr.BenqiCTokenList,
	ethaddr.TraderJoeProtocol: ethaddr.TraderjoeCTokenList,
}
