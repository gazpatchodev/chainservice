package simple

import (
	"github.com/gazpatchodev/chainservice/models"
	"github.com/gazpatchodev/chainservice/utils"
)

// Simple comment
type Simple struct{}

// New comment
func New() *Simple {
	return &Simple{}
}

// Parse comment
func (s *Simple) Parse(buf []byte) (bool, string, string, *[]models.Part) {
	if buf[0] != 0x6a {
		return false, "", "", nil
	}

	buf = buf[1:]

	var parts []models.Part

	for len(buf) > 0 {
		var d []byte
		d, buf = utils.ReadPushData(buf)
		var p models.Part
		p.MimeType = "text/plain; charset=utf-8"
		p.Data = string(d)
		parts = append(parts, p)
	}

	return true, "OP_RETURN", "", &parts
}
