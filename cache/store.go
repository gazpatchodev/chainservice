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
	"../models"
	"../parser"
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

// Vout comment
type Vout struct {
	N       int           `json:"vout"`
	Value   uint16        `json:"value"`
	Hex     string        `json:"hex,omitempty"`
	Type    string        `json:"type"`
	SubType string        `json:"subType,omitempty"`
	Text    string        `json:"text,omitempty"`
	Parts   []models.Part `json:"parts,omitempty"`
}

// OPReturnData comment
type OPReturnData struct {
	TxID      string `json:"txid"`
	BlockHash string `json:"blockHash"`
	BlockTime string `json:"blockTime"`
	Vouts     []Vout `json:"vouts"`
}

// GetOPReturnData comment
func GetOPReturnData(txid string) (opr *OPReturnData, err error) {
	opr, err = get(txid)
	if err != nil {
		if os.IsNotExist(err) {
			opr, err = GetOPReturnDataFromBitcoin(txid)
			if err != nil {
				logger.Infof("ERROR /%s (%+v)", txid, err)
			} else {
				put(opr)
				logger.Infof("BITCOIN /%s", txid)
			}
		}
	} else {
		logger.Infof("CACHE   /%s", txid)
	}

	return
}

// GetOPReturnDataFromBitcoin comment
func GetOPReturnDataFromBitcoin(txid string) (*OPReturnData, error) {
	log.Println("BITCOIN")
	tx, err := bitcoind.GetRawTransaction(txid)
	if err != nil {
		return nil, err
	}

	bt := "N/A"
	if tx.Blocktime != 0 {
		bt = time.Unix(0, int64(tx.Blocktime*1000*int64(time.Millisecond))).Format(time.RFC3339)
	}

	opr := &OPReturnData{
		TxID:      txid,
		BlockHash: tx.BlockHash,
		BlockTime: bt,
	}

	for _, vo := range tx.Vout {
		script, err := hex.DecodeString(vo.ScriptPubKey.Hex)
		if err != nil {
			return nil, fmt.Errorf("Failed to Decode script: %+v", err)
		}

		s, st, parts := parser.Parse(script)

		if parts != nil {
			vout := Vout{
				N:       vo.N,
				Type:    s,
				SubType: st,
				Value:   uint16(vo.Value),
				Hex:     vo.ScriptPubKey.Hex,
				Parts:   *parts,
			}

			opr.Vouts = append(opr.Vouts, vout)
		}
	}

	return opr, nil
}

func buildFilename(txid string) (string, error) {
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

	return path.Join(folder, fmt.Sprintf("%s.opr.json", txid)), nil
}

func prettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}

func put(opr *OPReturnData) (err error) {
	filename, err := buildFilename(opr.TxID)
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

func get(txid string) (opr *OPReturnData, err error) {
	filename, err := buildFilename(txid)
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
