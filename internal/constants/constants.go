package constants

import "github.com/shopspring/decimal"

var SecondsPerDay = decimal.NewFromInt(86400)                   // Seconds in a day.
var SecondsPerYear = SecondsPerDay.Mul(decimal.NewFromInt(365)) // Seconds in a year.

var RAYUnit = decimal.New(1, 27) // 1e27
var WEIUnit = decimal.New(1, 18) // 1e18

// Some annoying symbols that represent a meaningless token,
// which will be ignored by the program.
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
