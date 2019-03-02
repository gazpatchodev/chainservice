package tokenized

import (
	"github.com/gazpatchodev/chainservice/models"
	"github.com/gazpatchodev/chainservice/utils"
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
func (t *Tokenized) Parse(buf []byte) (bool, string, string, *[]models.Part) {
	if buf[0] != 0x6a {
		return false, "", "", nil
	}

	buf = buf[1:]
	var res []byte

	for len(buf) > 0 {
		var d []byte
		d, buf = utils.ReadPushData(buf)
		res = append(res, d...)
	}

	if len(res) > 6 {
		if res[0] == 0x00 && res[1] == 0x00 && res[2] == 0x00 && res[3] == 0x20 {
			// The next 2 bytes are the tokenizer action
			code := string(res[4:6])
			action, ok := tokenizedActions[code]
			if !ok {
				return false, "", "", nil
			}

			var p models.Part
			p.MimeType = "text/plain; charset=utf-8"
			// p.Hex = hex.EncodeToString(res[4:])
			p.Data = string(res[7:]) // Skip the 0x00
			var parts []models.Part
			parts = append(parts, p)
			return true, "Tokenized", action, &parts
		}
	}

	return false, "", "", nil

	// // func (t *Tokenized) Parse(script string) (action string, text string) {
	// strAsciiScript := string(buf)
	// trimmedAsciiScript := strings.TrimSpace(string(strAsciiScript))
	// prefix := trimmeAsciiScript[4:6]

	// action, ok := tokenizedActions[prefix]
	// if !ok {
	// 	return false, "", "", ""
	// }
	// text = trimmedAsciiScript[6:]
	// return
}
