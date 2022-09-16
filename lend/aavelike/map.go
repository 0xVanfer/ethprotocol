package aavelike

import "github.com/0xVanfer/ethaddr"

// map[protocol name] = a token list.
var ATokenListMap = map[string]map[string]map[string]string{
	ethaddr.AaveV2Protocol: ethaddr.AaveATokenV2List,
	ethaddr.AaveV3Protocol: ethaddr.AaveATokenV3List,
}

// map[protocol name] = v token list.
var VTokenListMap = map[string]map[string]map[string]string{
	ethaddr.AaveV2Protocol: ethaddr.AaveVTokenV2List,
	ethaddr.AaveV3Protocol: ethaddr.AaveVTokenV3List,
}

// map[protocol name] = s token list.
var STokenListMap = map[string]map[string]map[string]string{
	ethaddr.AaveV2Protocol: ethaddr.AaveSTokenV2List,
	ethaddr.AaveV3Protocol: ethaddr.AaveSTokenV3List,
}
