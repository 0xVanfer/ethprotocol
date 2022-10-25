package ethprotocol

import (
	"errors"

	"github.com/0xVanfer/chainId"
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/utils"
)

func (prot *Protocol) CheckNetwork() error {
	protName := prot.ProtocolBasic.ProtocolName
	network := prot.ProtocolBasic.Network
	var supportedNetworks []string
	switch protName {
	// aave v2
	case ethaddr.AaveV2Protocol:
		supportedNetworks = []string{chainId.AvalancheChainName, chainId.EthereumChainName}
	// aave v3
	case ethaddr.AaveV3Protocol:
		supportedNetworks = []string{chainId.AvalancheChainName, chainId.PolygonChainName}
	// axial
	case ethaddr.AxialProtocol:
		supportedNetworks = []string{chainId.AvalancheChainName}
	// benqi
	case ethaddr.BenqiProtocol:
		supportedNetworks = []string{chainId.AvalancheChainName}
	// curve
	case ethaddr.CurveProtocol:
		// curve supports almost every chain
		return nil
	// pangolin
	case ethaddr.PangolinProtocol:
		supportedNetworks = []string{chainId.AvalancheChainName}
	// platypus
	case ethaddr.PlatypusProtocol:
		supportedNetworks = []string{chainId.AvalancheChainName}
	// sushi
	case ethaddr.SushiProtocol:
		// sushi supports almost every chain
		return nil
	// traderjoe
	case ethaddr.TraderJoeProtocol:
		supportedNetworks = []string{chainId.AvalancheChainName}
	}
	// not supported
	if !utils.ContainInArrayX(network, supportedNetworks) {
		return errors.New(protName + " not supported " + network)
	}
	return nil
}
