package riff

import (
	"encoding/hex"
	"strings"

	"github.com/gazpatchodev/chainservice/models"
	"github.com/gazpatchodev/chainservice/utils"
)

// Riff comment
type Riff struct{}

// New comment
func New() *Riff {
	return &Riff{}
}

// Parse comment
func (s *Riff) Parse(buf []byte) (bool, string, string, *[]models.Part) {
	if buf[0] != 0x6a {
		return false, "", "", nil
	}

	var parts []models.Part

	wallet, buf := utils.ReadPushData(buf[1:])
	var walletPart models.Part
	walletPart.MimeType = "text/plain; charset=utf-8"
	walletPart.Data = string(wallet)
	parts = append(parts, walletPart)

	data, buf := utils.ReadPushData(buf)

	l := 100
	if l > len(data) {
		l = len(data)
	}
	if !strings.Contains(string(data[0:l]), "RIFF") {
		return false, "", "", nil
	}

	var dataPart models.Part
	dataPart.MimeType = "text/plain; charset=utf-8"
	dataPart.Data = hex.EncodeToString(data)
	parts = append(parts, dataPart)

	proto, buf := utils.ReadPushData(buf)
	var protoPart models.Part
	protoPart.MimeType = "audio/wav"
	protoPart.Data = string(proto)
	parts = append(parts, protoPart)

	return true, "RIFF", string(proto), &parts
}
