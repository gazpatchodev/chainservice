package yours

import (
	"github.com/gazpatchodev/chainservice/models"
	"github.com/gazpatchodev/chainservice/utils"
)

// Yours comment
type Yours struct{}

// New comment
func New() *Yours {
	return &Yours{}
}

// Parse comment
func (t *Yours) Parse(buf []byte) (bool, string, string, *[]models.Part) {
	if buf[0] != 0x6a {
		return false, "", "", nil
	}

	d, _ := utils.ReadPushData(buf[1:])

	s := string(d)

	if s == "yours.org" {
		return true, "yours.org", "", nil
	}

	return false, "", "", nil

}
