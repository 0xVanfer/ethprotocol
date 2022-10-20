package lending

import (
	"errors"

	"github.com/0xVanfer/ethaddr"
	"github.com/0xVanfer/ethprotocol/model"
	"github.com/0xVanfer/utils"
)

// Initialize the lend pool protocol basic and pool type.
func (p *LendingPool) Init(protocolBasic model.ProtocolBasic) error {
	// regular check
	if !protocolBasic.Regularcheck() {
		return errors.New("protocol basic must not be initialized")
	}
	// decide the pool type
	if utils.ContainInArrayX(protocolBasic.ProtocolName, ethaddr.AaveLikeProtocols) {
		p.PoolType = LendingPoolType{IsAaveLike: true}
	} else if utils.ContainInArrayX(protocolBasic.ProtocolName, ethaddr.CompoundLikeProtocols) {
		p.PoolType = LendingPoolType{IsCompoundLike: true}
	} else {
		return errors.New("protocol not supported lend pool")
	}
	p.ProtocolBasic = &protocolBasic
	return nil
}
