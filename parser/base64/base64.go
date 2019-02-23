package base64

import (
	"regexp"

	"github.com/gazpatchodev/chainservice/models"
	"github.com/gazpatchodev/chainservice/utils"
)

// Base64 comment
type Base64 struct{}

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
func New() *Base64 {
	return &Base64{}
}

// Parse comment
func (t *Base64) Parse(buf []byte) (bool, string, string, *[]models.Part) {
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

	s := string(res)

	var base64Type string

	// We only want to search the first 100 characters or the whole string if it's shorter.
	size := len(s)
	if size > 100 {
		size = 100
	}

	if base64 := base64RE.FindStringSubmatch(s[0:size]); len(base64) > 0 {
		if len(base64) > 1 {
			base64Type = base64[1]
		}

		if t := typeRE.FindStringSubmatch(s[0:size]); len(t) > 0 {
			if len(t) > 1 {
				base64Type = t[1]
			}
		}

		sr2 := sourceRE.FindStringSubmatch(s)
		_ = sr2
		if sr := sourceRE.FindStringSubmatch(s); len(sr) > 0 {
			if len(sr) > 1 {
				var p models.Part
				p.BASE64 = sr[1]
				var parts []models.Part
				parts = append(parts, p)
				return true, "Base64", base64Type, &parts
			}
		}
	}

	return false, "", "", nil

}
