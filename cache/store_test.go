package cache

import (
	"os"
	"testing"
)

func TestPut(t *testing.T) {
	d := &OPReturnData{
		TxID:      "bbae12b5f4c088d940733dcd1455efc6a3a69cf9340e17a981286d3778615684",
		BlockHash: "1234567890abcdef",
	}

	err := put(d)
	if err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	d := &OPReturnData{
		TxID:      "bbae12b5f4c088d940733dcd1455efc6a3a69cf9340e17a981286d3778615684",
		BlockHash: "1234567890abcdef",
	}

	err := put(d)
	if err != nil {
		t.Error(err)
	}

	d2, err := get("bbae12b5f4c088d940733dcd1455efc6a3a69cf9340e17a981286d3778615684")
	if err != nil {
		t.Error(err)
	}

	t.Logf("%+v", d2)
}

func TestYoursFailing(t *testing.T) {
	d, err := GetOPReturnDataFromBitcoin("6ee609d1770e0e9dcad0fd216bba78109aa66715c2a8ebb9f6fdc0302b86ddc2")
	if err != nil {
		t.Error(err)
	}

	t.Logf("%+v", d)
}

func TestGetRealExample(t *testing.T) {
	d, err := GetOPReturnData("8bae12b5f4c088d940733dcd1455efc6a3a69cf9340e17a981286d3778615684")
	if err != nil {
		t.Error(err)
	}

	t.Logf("%+v", d)
}

func TestGetNoFile(t *testing.T) {
	d, err := get("NOFILE_NOFILE_NOFILE_NOFILE_NOFILE_NOFILE_NOFILE_NOFILE_NOFILE_X")
	if err != nil {
		if os.IsNotExist(err) {
			t.Log("Cache miss")
		} else {
			t.Error(err)
		}
	} else {
		t.Logf("%+v", d)
	}
}
