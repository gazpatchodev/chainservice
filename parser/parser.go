package parser

import (
	"../cache"
	"./base64"
	"./simple"
	"./tokenized"
)

// Parser interface
type Parser interface {
	Parse(buf []byte) (bool, string, string, *[]cache.Part)
}

var parsers []Parser

func init() {
	parsers = append(parsers, tokenized.New())
	parsers = append(parsers, base64.New())
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
