package parser

import (
	"encoding/hex"
	"io/ioutil"
	"testing"
)

func TestATokenizedScript(t *testing.T) {
	buf, _ := hex.DecodeString("6a4cd20000002043310041424320436f20436f6d6d6f6e2045712041677265656d656e74000000000000000000000000000000000000000000000000000000000000000000000000000047425200004742520000000001a984ad75190000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000041424320436f6d70616e79000000000000546f6b656e697a656400000000000000007f4e000000000000000000000000000000")

	s, _, _ := Parse(buf)

	if s != "TOKENIZED" {
		t.Errorf("Expected %q, got %q", "TOKENIZED", s)
	}
}

func TestSimpleScript(t *testing.T) {
	buf, _ := hex.DecodeString("6a13636861726c6579206c6f766573206865696469")
	s, st, p := Parse(buf)

	if s != "SIMPLE" {
		t.Errorf("Expected %q, got %q", "SIMPLE", s)
	}

	if st != "" {
		t.Errorf("Expected %q, got %q", "", st)
	}

	if p.UTF8 != "charley loves heidi" {
		t.Errorf("Expected %q, got %q", "charley loves heidi", p.UTF8)
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

	if s != "BASE64" {
		t.Errorf("Expected %q, got %q", "BASE64", s)
	}

	if st != "video/mp4" {
		t.Errorf("Expected %q, got %q", "video/mp4", st)
	}
}
