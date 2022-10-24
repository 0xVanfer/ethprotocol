package constants

import "github.com/shopspring/decimal"

var SecondsPerDay = decimal.New(86400, 0)
var SecondsPerYear = SecondsPerDay.Mul(decimal.New(365, 0))

var RAYUnit = decimal.New(1, 27)
var WEIUnit = decimal.New(1, 18)

var IgnoreSymbols = []string{
	"ot-qiusdc-28dec2023",
	"ot-jlp-28dec2023",
	"ot-qiavax-28dec2023",
	"ot-xjoe-30jun2022",
	"ot-wmemo-24feb2022",
	"ot-jlp-29dec2022",
	"testpendle",
	"💎",
	"🐆",
	"🔥",
	"jlp",
	"",
	"bog",
	"remo",
	"ot-wxbtrfly-21apr2022",
	"ot-cdai-29dec2022",
	"ot-ausdc-30dec2021",
	"aaa",
	"bbb",
	"ccc",
	"ddd",
	"aa",
	"bb",
	"cc",
	"dd",
	"420",
}
