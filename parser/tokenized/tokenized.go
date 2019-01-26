package tokenized

import (
	"encoding/hex"

	"../../cache"
	"../../utils"
)

// Tokenized comment
type Tokenized struct{}

var tokenizedActions = map[string]string{
	"C1": "Contract Offer",
	"C2": "Contract Formation",
	"C3": "Contract Ammendment",

	"A1": "Asset Definition",
	"A2": "Asset Creation",
	"A3": "Asset Modification",

	"T1": "Send",
	"T2": "Exchange",
	"T3": "Swap",
	"T4": "Settlement",

	"G1": "Initiative",
	"G2": "Referendum",
	"G3": "Vote",
	"G4": "Ballot Cast",
	"G5": "Ballot Counted",
	"G6": "Result",

	"E1": "Order",
	"E2": "Freeze",
	"E3": "Thaw",
	"E4": "Confiscation",
	"E5": "Reconcilliation",

	"M1": "Message",
	"M2": "Rejection",

	"R1": "Establishment",
	"R2": "Addition",
	"R3": "Alteration",
	"R4": "Removal",
}

// New comment
func New() *Tokenized {
	return &Tokenized{}
}

// Parse comment
func (t *Tokenized) Parse(buf []byte) (bool, string, string, *cache.Part) {
	if buf[0] != 0x6a {
		return false, "", "", nil
	}

	d := utils.ReadPushData(buf[1:])

	if len(d) > 4 {
		if d[0] == 0x00 && d[1] == 0x00 && d[2] == 0x00 && d[3] == 0x20 {
			d := utils.ReadPushData(buf[4:])
			var p cache.Part
			p.Hex = hex.EncodeToString(d)
			p.UTF8 = string(d)
			return true, "TOKENIZED", "", &p
		}
	}

	return false, "", "", nil

	// // func (t *Tokenized) Parse(script string) (action string, text string) {
	// strAsciiScript := string(buf)
	// trimmedAsciiScript := strings.TrimSpace(string(strAsciiScript))
	// prefix := trimmedAsciiScript[4:6]

	// action, ok := tokenizedActions[prefix]
	// if !ok {
	// 	return false, "", "", ""
	// }
	// text = trimmedAsciiScript[6:]
	// return
}
