package moneybutton

import (
	"regexp"

	"../../cache"
	"../../utils"
)

// MoneyButton comment
type MoneyButton struct{}

var (
	// With "<video controls><source type='video/mp4' src="data:video/webm;base64,AAAAI..."

	// Look for src="data:......;base64". If this regexp finds a match, we
	// have a data uri.  Use the type embedded in data uri.
	base64RE = regexp.MustCompile(`(?i)\Wsrc\s*=\s*['"]data:(.+?);base64,`)

	// Look for the type="something" which will then override the type from previous regexp
	typeRE = regexp.MustCompile(`(?i)\Wtype\s*=\s*['"](.+?)['"]`)

	// Now extact the src= from the data uri
	sourceRE = regexp.MustCompile(`(?i)(?m)src\s*=\s*['"](.+)['"]`)
)

// New comment
func New() *MoneyButton {
	return &MoneyButton{}
}

// Parse comment
func (t *MoneyButton) Parse(buf []byte) (bool, string, string, *[]cache.Part) {
	if buf[0] != 0x6a {
		return false, "", "", nil
	}

	prefix, buf := utils.ReadPushData(buf[1:])

	if string(prefix) == "moneybutton.com" {
		var encoding []byte
		encoding, buf = utils.ReadPushData(buf)

		if string(encoding) == "utf8" {
			var a []byte
			a, buf = utils.ReadPushData(buf)

			audio := a

			var p cache.Part
			p.URI = string(audio)
			var parts []cache.Part
			parts = append(parts, p)
			return true, "moneybutton.com", "URI", &parts
		}
	}

	return false, "", "", nil

}
