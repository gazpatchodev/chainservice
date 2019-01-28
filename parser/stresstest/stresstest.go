package stresstest

import (
	"../../models"
	"../../utils"
)

// StressTest comment
type StressTest struct{}

// New comment
func New() *StressTest {
	return &StressTest{}
}

// Parse comment
func (t *StressTest) Parse(buf []byte) (bool, string, string, *[]models.Part) {
	if buf[0] != 0x6a {
		return false, "", "", nil
	}

	d, _ := utils.ReadPushData(buf[1:])

	s := string(d)

	if s == "stresstestbitcoin.cash" {
		return true, "stresstestbitcoin.cash", "", nil
	}

	return false, "", "", nil

}
