package lending

import (
	"errors"

	"github.com/0xVanfer/chainId"
	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/model"
)

// Initialize the lend pool protocol basic and pool type.
func (p *LendingPool) Init(protocolBasic model.ProtocolBasic) error {
	// regular check
	if !protocolBasic.Regularcheck() {
		return errors.New("protocol basic must not be initialized")
	}
	p.ProtocolBasic = &protocolBasic

	// // decide the pool type
	// if utils.ContainInArrayX(protocolBasic.ProtocolName, ethaddr.AaveLikeProtocols) {
	// 	p.PoolType = LendingPoolType{IsAaveLike: true}
	// } else if utils.ContainInArrayX(protocolBasic.ProtocolName, ethaddr.CompoundLikeProtocols) {
	// 	p.PoolType = LendingPoolType{IsCompoundLike: true}
	// } else {
	// 	return errors.New("protocol not supported lend pool")
	// }

	// decide network
	switch protocolBasic.ProtocolName {
	// aave v2
	case ethaddr.AaveV2Protocol:
		p.PoolType = LendingPoolType{IsAaveLike: true}
		p.SupportedNetworks = []string{
			chainId.AvalancheChainName,
			chainId.EthereumChainName,
		}
	// aave v3
	case ethaddr.AaveV3Protocol:
		p.PoolType = LendingPoolType{IsAaveLike: true}
		p.SupportedNetworks = []string{
			chainId.AvalancheChainName,
			chainId.PolygonChainName,
		}
	// benqi
	case ethaddr.BenqiProtocol:
		p.PoolType = LendingPoolType{IsCompoundLike: true}
		p.SupportedNetworks = []string{
			chainId.AvalancheChainName,
		}
	// tradejoe
	case ethaddr.TraderJoeProtocol:
		p.PoolType = LendingPoolType{IsCompoundLike: true}
		p.SupportedNetworks = []string{
			chainId.AvalancheChainName,
		}
	default:
		return errors.New(protocolBasic.ProtocolName + " lend pools not supporteds")
	}
	return nil
}
