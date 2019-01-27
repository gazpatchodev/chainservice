package parser

import (
	"encoding/hex"
	"io/ioutil"
	"testing"
)

func TestATokenizedScript(t *testing.T) {
	buf, _ := hex.DecodeString("6a4cd20000002043310041424320436f20436f6d6d6f6e2045712041677265656d656e74000000000000000000000000000000000000000000000000000000000000000000000000000047425200004742520000000001a984ad75190000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000041424320436f6d70616e79000000000000546f6b656e697a656400000000000000007f4e000000000000000000000000000000")

	s, _, _ := Parse(buf)

	if s != "Tokenized" {
		t.Errorf("Expected %q, got %q", "Tokenized", s)
	}
}

func TestSimpleScript(t *testing.T) {
	buf, _ := hex.DecodeString("6a13636861726c6579206c6f766573206865696469")
	s, st, p := Parse(buf)

	if s != "OP_RETURN" {
		t.Errorf("Expected %q, got %q", "OP_RETURN", s)
	}

	if st != "" {
		t.Errorf("Expected %q, got %q", "", st)
	}

	if (*p)[0].UTF8 != "charley loves heidi" {
		t.Errorf("Expected %q, got %q", "charley loves heidi", (*p)[0].UTF8)
	}
}

func TestVideo(t *testing.T) {
	h, err := ioutil.ReadFile("video.hex")
	if err != nil {
		t.Error(err)
	}
	buf, err := hex.DecodeString(string(h))
	if err != nil {
		t.Error(err)
	}

	s, st, _ := Parse(buf)

	if s != "Base64" {
		t.Errorf("Expected %q, got %q", "Base64", s)
	}

	if st != "video/mp4" {
		t.Errorf("Expected %q, got %q", "video/mp4", st)
	}
}

func TestYoursScript(t *testing.T) {
	buf, _ := hex.DecodeString("6a09796f7572732e6f7267")
	s, st, p := Parse(buf)

	if s != "yours.org" {
		t.Errorf("Expected %q, got %q", "yours.org", s)
	}

	if st != "" {
		t.Errorf("Expected %q, got %q", "", st)
	}

	if p != nil {
		t.Errorf("Expected nil, got %v", p)
	}

}
func TestMemoCashTipScript(t *testing.T) {
	buf, _ := hex.DecodeString("6a026d0420f9fed7ac794200388a1e1ef4b84d5df7bb49dd297a0f71d6d0e5ddecf0d545dd")
	s, st, p := Parse(buf)

	if s != "memo.cash" {
		t.Errorf("Expected %q, got %q", "yours.org", s)
	}

	if st != "Like / tip memo" {
		t.Errorf("Expected %q, got %q", "Like / tip memo", st)
	}

	if len((*p)[0].Hex) != 64 {
		t.Errorf("Expected 32 byte hash, got %v", (*p)[0].Hex)
	}

}