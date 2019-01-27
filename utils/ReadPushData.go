package utils

import "encoding/binary"

const (
	opRETURN    = 0x6a
	opPUSHDATA1 = 0x4c
	opPUSHDATA2 = 0x4d
	opPUSHDATA4 = 0x4e
)

// ReadPushData comment
func ReadPushData(buf []byte) ([]byte, []byte) {
	if len(buf) <= 1 {
		return []byte{}, []byte{}
	}

	switch buf[0] {
	case 0x00:
		return []byte{}, buf[1:]

	case opPUSHDATA1: // 04c L0 D0 D1...
		if len(buf) < 2 {
			return []byte{}, []byte{}
		}

		l := int(buf[1])
		if len(buf) < l+2 {
			return buf[2:], []byte{}
		}

		return buf[2 : l+2], buf[l+2:]

	case opPUSHDATA2: // 04d L0 L1 D0 D1...
		if len(buf) < 3 {
			return []byte{}, []byte{}
		}

		l := binary.LittleEndian.Uint16(buf[1:])
		if uint16(len(buf)) < l+3 {
			return buf[3:], []byte{}
		}

		return buf[3 : l+3], buf[l+3:]

	case opPUSHDATA4: // 04e L0 L1 L2 L3 D0 D1...
		if len(buf) < 5 {
			return []byte{}, []byte{}
		}

		l := binary.LittleEndian.Uint32(buf[1:])
		if uint32(len(buf)) < l+5 {
			return buf[5:], []byte{}
		}

		return buf[5 : l+5], buf[l+5:]

	default:
		l := int(buf[0])
		if len(buf) < l+1 {
			return buf[1:], []byte{}
		}

		return buf[1 : l+1], buf[l+1:]
	}
}
