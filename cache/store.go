package cache

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/ordishs/gocore"

	"../bitcoin"
)

var (
	logger      = gocore.Log("ChainService:store")
	address, _  = gocore.Config().Get("bitcoinAddress")
	user, _     = gocore.Config().Get("bitcoinUser")
	password, _ = gocore.Config().Get("bitcoinPassword")
	bitcoind, _ = bitcoin.New(address, user, password)
	cacheDir, _ = gocore.Config().Get("cacheDir")
)

func init() {
	err := os.MkdirAll(cacheDir, os.ModePerm)
	if err != nil {
		logger.Fatalf("Unable to make directory %s: %+v", cacheDir, err)
	}
}

// OPReturnData comment
type OPReturnData struct {
	TxID      string      `json:"txid"`
	Vout      uint16      `json:"vout"`
	BlockHash string      `json:"blockHash"`
	BlockTime string      `json:"blockTime"`
	Value     uint16      `json:"value"`
	Hex       string      `json:"hex,omitempty"`
	Type      string      `json:"type"`
	SubType   string      `json:"subType,omitempty"`
	Text      string      `json:"text,omitempty"`
	Parts     []Part      `json:"parts,omitempty"`
	Err       interface{} `json:"error,omitempty"`
}

// Part struct
type Part struct {
	Hex    string `json:"hex,omitempty"`
	UTF8   string `json:"utf8,omitempty"`
	BASE64 string `json:"base64,omitempty"`
}

// GetOPReturnData comment
func GetOPReturnData(txid string, vout uint16) (opr *OPReturnData, err error) {
	opr, err = get(txid, vout)
	if err != nil {
		if os.IsNotExist(err) {
			opr, err = GetOPReturnDataFromBitcoin(txid, vout)
			if err == nil {
				put(opr)
			}
		}
	}

	return
}

// GetOPReturnDataFromBitcoin comment
func GetOPReturnDataFromBitcoin(txid string, vout uint16) (*OPReturnData, error) {
	log.Println("BITCOIN")
	tx, err := bitcoind.GetRawTransaction(txid)
	if err != nil {
		return nil, err
	}

	if int(vout) > len(tx.Vout) {
		return nil, fmt.Errorf("vout %d does not exist in this transaction", vout)
	}

	script, err := hex.DecodeString(tx.Vout[vout].ScriptPubKey.Hex)
	if err != nil {
		return nil, fmt.Errorf("Failed to Decode script: %+v", err)
	}

	// Make sure this is an OP_RETURN script...
	if script[0] != 0x6a {
		return nil, fmt.Errorf("vout %d is not an OP_RETURN script", vout)
	}

	bt := "N/A"
	if tx.Blocktime != 0 {
		bt = time.Unix(0, int64(tx.Blocktime*1000*int64(time.Millisecond))).Format(time.RFC3339)
	}

	parts := readPushDatas(script[1:])

	opr := &OPReturnData{
		TxID:      txid,
		Vout:      vout,
		BlockHash: tx.BlockHash,
		BlockTime: bt,
		Value:     uint16(tx.Vout[vout].Value),
		Hex:       tx.Vout[vout].ScriptPubKey.Hex,
		Parts:     parts,
	}

	return opr, nil
}

func readPushDatas(buf []byte) []Part {
	var parts []Part

	var part Part
	switch buf[0] {
	case 0x4c:
		part.Hex = hex.EncodeToString(buf[2:])
		part.UTF8 = string(buf[2:])
	case 0x4d:
		part.Hex = hex.EncodeToString(buf[3:])
		part.UTF8 = string(buf[3:])
	case 0x4e:
		part.Hex = hex.EncodeToString(buf[5:])
		part.UTF8 = string(buf[5:])
	default:
		part.Hex = hex.EncodeToString(buf[1:])
		part.UTF8 = string(buf[1:])
	}

	parts = append(parts, part)
	return parts
}

func buildFilename(txid string, vout uint16) (string, error) {
	if len(txid) != 64 {
		return "", errors.New("txid must be 32 bytes")
	}

	// The path is made up of abcdefgh....:
	// base folder/ab/abcd/abcdefgh.....
	folder := path.Join(cacheDir, txid[0:2], txid[0:4])
	err := os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		return "", err
	}

	return path.Join(folder, fmt.Sprintf("%s.%d.json", txid, vout)), nil
}

func prettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}

func put(opr *OPReturnData) (err error) {
	filename, err := buildFilename(opr.TxID, opr.Vout)
	if err != nil {
		return
	}

	j, err := json.Marshal(opr)
	if err != nil {
		return
	}

	pj, err := prettyprint(j)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(filename, pj, 0644)
	return
}

func get(txid string, vout uint16) (opr *OPReturnData, err error) {
	filename, err := buildFilename(txid, vout)
	if err != nil {
		return
	}

	j, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	err = json.Unmarshal(j, &opr)
	return
}
