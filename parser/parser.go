package parser

import (
	"github.com/gazpatchodev/chainservice/models"
	"github.com/gazpatchodev/chainservice/parser/base64"
	"github.com/gazpatchodev/chainservice/parser/memocash"
	"github.com/gazpatchodev/chainservice/parser/moneybutton"
	"github.com/gazpatchodev/chainservice/parser/riff"
	"github.com/gazpatchodev/chainservice/parser/simple"
	"github.com/gazpatchodev/chainservice/parser/stresstest"
	"github.com/gazpatchodev/chainservice/parser/tokenized"
	"github.com/gazpatchodev/chainservice/parser/yours"
)

// Parser interface
type Parser interface {
	Parse(buf []byte) (bool, string, string, *[]models.Part)
}

var parsers []Parser

func init() {
	// The order of these IS important.
	parsers = append(parsers, riff.New())
	parsers = append(parsers, yours.New())
	parsers = append(parsers, stresstest.New())
	parsers = append(parsers, tokenized.New())
	parsers = append(parsers, moneybutton.New())
	parsers = append(parsers, base64.New())
	parsers = append(parsers, memocash.New())
	parsers = append(parsers, simple.New())
}

// Parse comment
func Parse(buf []byte) (string, string, *[]models.Part) {
	for _, p := range parsers {
		match, t, st, parts := p.Parse(buf)
		if match {
			return t, st, parts
		}
	}
	return "UNKNOWN", "", nil
}
