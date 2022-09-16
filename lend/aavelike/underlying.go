package aavelike

import "strings"

// Use a token to update underlying.
func (t *AToken) UpdateUnderlying() {
	aList := ATokenListMap[t.ProtocolBase.ProtocolName]
	for underlying, atoken := range aList[t.ProtocolBase.Network] {
		if strings.EqualFold(atoken, t.Basic.Address) {
			t.UnderlyingBasic.Init(underlying, t.ProtocolBase.Network, t.ProtocolBase.Client)
		}
	}
}

// Use v token to update underlying.
func (t *VToken) UpdateUnderlying() {
	vList := VTokenListMap[t.ProtocolBase.ProtocolName]
	for underlying, vtoken := range vList[t.ProtocolBase.Network] {
		if strings.EqualFold(vtoken, t.Basic.Address) {
			t.UnderlyingBasic.Init(underlying, t.ProtocolBase.Network, t.ProtocolBase.Client)
		}
	}
}

// Use s token to update underlying.
func (t *SToken) UpdateUnderlying() {
	sList := STokenListMap[t.ProtocolBase.ProtocolName]
	for underlying, stoken := range sList[t.ProtocolBase.Network] {
		if strings.EqualFold(stoken, t.Basic.Address) {
			t.UnderlyingBasic.Init(underlying, t.ProtocolBase.Network, t.ProtocolBase.Client)
		}
	}
}
