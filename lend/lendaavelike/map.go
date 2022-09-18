package lendaavelike

import "github.com/0xVanfer/ethaddr"

// map[protocol name][network][underlying] = a token.
var ATokenListMap = map[string]map[string]map[string]string{
	ethaddr.AaveV2Protocol: ethaddr.AaveATokenV2List,
	ethaddr.AaveV3Protocol: ethaddr.AaveATokenV3List,
}

// map[protocol name][network][underlying] = v token.
var VTokenListMap = map[string]map[string]map[string]string{
	ethaddr.AaveV2Protocol: ethaddr.AaveVTokenV2List,
	ethaddr.AaveV3Protocol: ethaddr.AaveVTokenV3List,
}

// map[protocol name][network][underlying] = s token.
var STokenListMap = map[string]map[string]map[string]string{
	ethaddr.AaveV2Protocol: ethaddr.AaveSTokenV2List,
	ethaddr.AaveV3Protocol: ethaddr.AaveSTokenV3List,
}
