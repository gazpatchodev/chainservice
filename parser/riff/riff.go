package riff

import (
	"encoding/hex"
	"strings"

	"../../models"
	"../../utils"
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
	walletPart.Hex = hex.EncodeToString(wallet)
	walletPart.UTF8 = string(wallet)
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
	dataPart.Hex = hex.EncodeToString(data)
	parts = append(parts, dataPart)

	proto, buf := utils.ReadPushData(buf)
	var protoPart models.Part
	protoPart.Hex = hex.EncodeToString(proto)
	protoPart.UTF8 = string(proto)
	parts = append(parts, protoPart)

	return true, "RIFF", string(proto), &parts
}
