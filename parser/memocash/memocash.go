package memocash

import (
	"encoding/hex"

	"../../cache"
	"../../utils"
)

// MemoCash comment
type MemoCash struct{}

var memoPrefixes = map[int]string{
	0x01: "Set name",
	0x02: "Post memo",
	0x03: "Reply to memo",
	0x04: "Like / tip memo",
	0x05: "Set profile text",
	0x06: "Follow user",
	0x07: "Unfollow user",
	0x0a: "Set profile picture",
	0x0b: "Repost memo", // Planned
	0x0c: "Post topic message",
	0x0d: "Topic follow",
	0x0e: "Topic unfollow",
	0x10: "Create poll",
	0x13: "Add poll option",
	0x14: "Poll vote",
	0x24: "Send money", // Planned
}

// New comment
func New() *MemoCash {
	return &MemoCash{}
}

// Parse comment
func (t *MemoCash) Parse(buf []byte) (bool, string, string, *[]cache.Part) {
	if buf[0] != 0x6a {
		return false, "", "", nil
	}

	res, buf := utils.ReadPushData(buf[1:])

	if len(res) == 2 && res[0] == 0x6d {
		// Looks like a memo cash transaction (https://memo.cash/protocol)
		var parts []cache.Part

		switch res[1] {
		case 0x04: // Like / Tip memo - txhash(32)
			fallthrough
		case 0x06: // Follow user - address(35)
			var p cache.Part
			var d []byte
			d, _ = utils.ReadPushData(buf)
			p.Hex = hex.EncodeToString(d)
			parts = append(parts, p)

		case 0x03: // Reply to memo
			fallthrough
		case 0x0b: // Repost memo
			fallthrough
		case 0x13: // Add poll option
			fallthrough
		case 0x14: // Poll vote
			// These guys have 32 bytes for the hash
			// The rest is ASCII
			// Ex: 6a 02 6d03 20 ae2102f5fba5e3464447cb65bec61deedf6c1df2f919bb4190335a06c4834796 4ca47468652072657374617572616e74206275696c7420616e64207461626c657320736574206265666f72652074686520637573746f6d6572732077616c6b20696e2074686520646f6f722e20416e6420697427732062657474657220746f2068617665207468652062696767657374206368616e676573207768656e20746865206665776573742070656f706c6520617265207573696e6720746865206e6574776f726b2c
			var p1 cache.Part
			var d1 []byte
			d1, buf = utils.ReadPushData(buf)
			p1.Hex = hex.EncodeToString(d1)
			parts = append(parts, p1)

			var p2 cache.Part
			var d2 []byte
			d2, buf = utils.ReadPushData(buf)
			p2.UTF8 = string(d2)
			parts = append(parts, p2)
		}
		return true, "memo.cash", memoPrefixes[int(res[1])], &parts
	}
	return false, "", "", nil
}
