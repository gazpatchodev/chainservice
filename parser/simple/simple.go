package simple

import (
	"encoding/hex"

	"../../cache"
	"../../utils"
)

// Simple comment
type Simple struct{}

// New comment
func New() *Simple {
	return &Simple{}
}

// Parse comment
func (s *Simple) Parse(buf []byte) (bool, string, string, *cache.Part) {
	if buf[0] != 0x6a {
		return false, "", "", nil
	}

	d := utils.ReadPushData(buf[1:])
	var p cache.Part
	p.Hex = hex.EncodeToString(d)
	p.UTF8 = string(d)

	return true, "SIMPLE", "", &p
}
