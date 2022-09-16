package apy

import (
	"math"

	"github.com/0xVanfer/types"
)

type ApyInfo struct {
	Apy          float64
	Apr          float64
	ApyIncentive float64
	AprIncentive float64
}

func Apr2Apy[T types.Number](apr T) (apy float64) {
	return math.Pow((1+types.ToFloat64(apr)/365), 365) - 1
}

func Apy2Apr[T types.Number](apy T) (apr float64) {
	return (math.Pow(1+types.ToFloat64(apy), 1.0/365) - 1) * 365
}
