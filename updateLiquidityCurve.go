package ethprotocol

import "github.com/0xVanfer/ethprotocol/internal/requests"

// todo
func (prot *Protocol) updateLiquidityCurve() error {
	// check network
	network := prot.ProtocolBasic.Network
	err := prot.CheckNetwork()
	if err != nil {
		return err
	}
	// volume and apy
	volumeReq, err := requests.ReqCurveVolumeApr(network)
	if err != nil {
		panic(err)
	}
	// main pools info, type:"usd"
	mainPools, err := requests.ReqCurvePoolMain(network)
	if err != nil {
		panic(err)
	}

	// crypto pools info
	cryptoPool, err := requests.ReqCurveCrypto(network)
	if err != nil {
		panic(err)
	}
	allPools := append(mainPools.Data.PoolData, cryptoPool.Data.PoolData...)
	_ = volumeReq
	_ = allPools

	return nil
}
