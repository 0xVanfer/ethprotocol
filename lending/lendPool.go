package lending

import (
	"errors"

	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/model"
	"github.com/0xVanfer/utils"
)

// Initialize the lend pool protocol basic and pool type.
func (p *LendingPool) Init(protocolBasic model.ProtocolBasic) error {
	if utils.ContainInArrayX(protocolBasic.ProtocolName, ethaddr.AaveLikeProtocols) {
		p.PoolType = LendingPoolType{IsAaveLike: true}
	} else if utils.ContainInArrayX(protocolBasic.ProtocolName, ethaddr.CompoundLikeProtocols) {
		p.PoolType = LendingPoolType{IsCompoundLike: true}
	} else {
		return errors.New("protocol not supported lend pool")
	}
	if protocolBasic.Network == "" {
		return errors.New("network must not be empty")
	}
	p.ProtocolBasic = &protocolBasic
	return nil
}
