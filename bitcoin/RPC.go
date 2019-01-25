package bitcoin

import (
	"encoding/json"
)

// A Bitcoind represents a Bitcoind client
type Bitcoind struct {
	client *rpcClient
}

// New return a new bitcoind
func New(address string, user, passwd string) (*Bitcoind, error) {
	rpcClient, err := newClient(address, user, passwd)
	if err != nil {
		return nil, err
	}
	return &Bitcoind{rpcClient}, nil
}

// GetRawTransaction comment
func (b *Bitcoind) GetRawTransaction(txhash string) (tx *RawTransaction, err error) {
	r, err := b.client.call("getrawtransaction", []interface{}{txhash, 1})
	if err != nil {
		return
	}

	json.Unmarshal(r.Result, &tx)

	return
}
