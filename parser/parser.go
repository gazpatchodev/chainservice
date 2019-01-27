package parser

import (
	"../cache"
	"./base64"
	"./memocash"
	"./moneybutton"
	"./simple"
	"./stresstest"
	"./tokenized"
	"./yours"
)

// Parser interface
type Parser interface {
	Parse(buf []byte) (bool, string, string, *[]cache.Part)
}

var parsers []Parser

func init() {
	// The order of these IS important.
	parsers = append(parsers, yours.New())
	parsers = append(parsers, stresstest.New())
	parsers = append(parsers, tokenized.New())
	parsers = append(parsers, moneybutton.New())
	parsers = append(parsers, base64.New())
	parsers = append(parsers, memocash.New())
	parsers = append(parsers, simple.New())
}

// Parse comment
func Parse(buf []byte) (string, string, *[]cache.Part) {
	for _, p := range parsers {
		match, t, st, parts := p.Parse(buf)
		if match {
			return t, st, parts
		}
	}
	return "UNKNOWN", "", nil
}
