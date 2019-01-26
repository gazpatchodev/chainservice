package utils

const (
	opRETURN    = 0x6a
	opPUSHDATA1 = 0x4c
	opPUSHDATA2 = 0x4d
	opPUSHDATA4 = 0x4e
)

// ReadPushData comment
func ReadPushData(buf []byte) []byte {
	switch buf[0] {
	case opPUSHDATA1:
		return buf[2:]
	case opPUSHDATA2:
		return buf[3:]
	case opPUSHDATA4:
		return buf[5:]
	default:
		return buf[1:]
	}
}
