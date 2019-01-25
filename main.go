package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"./bitcoin"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/ordishs/gocore"
)

var (
	logger           = gocore.Log("ChainService")
	listenAddress, _ = gocore.Config().Get("listenAddress")
	address, _       = gocore.Config().Get("bitcoinAddress")
	user, _          = gocore.Config().Get("bitcoinUser")
	password, _      = gocore.Config().Get("bitcoinPassword")
	bitcoind, _      = bitcoin.New(address, user, password, false)
)

type part struct {
	Hex  string `json:"hex,omitempty"`
	UTF8 string `json:"utf8,omitempty"`
}

type response struct {
	TxID      string      `json:"txid"`
	BlockHash string      `json:"blockHash"`
	BlockTime string      `json:"blockTime"`
	Value     uint16      `json:"value"`
	Hex       string      `json:"hex,omitempty"`
	Type      string      `json:"type"`
	SubType   string      `json:"subType,omitempty"`
	Parts     []part      `json:"parts,omitempty"`
	Err       interface{} `json:"error,omitempty"`
}

func main() {
	stats := gocore.Config().Stats()
	logger.Infof("STATS\n%s\nVERSION\n-------\n%s (%s)\n\n", stats, version, commit)

	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt)

	go func() {
		<-signalChan

		appCleanup()
		os.Exit(1)
	}()

	start()
}

func appCleanup() {
	logger.Infof("Chain Service shutting dowm...")
}

func start() {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/bsv/{txhash}/{vout:[0-9]+}", getTransactionOutput).Methods("GET")
	r.Handle("/", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(errorHandler)))
	// r.HandleFunc("/articles/{category}/{id:[0-9]+}", ArticleHandler)
	// Wrap our server with our gzip handler to gzip compress all responses.
	log.Fatal(http.ListenAndServe(listenAddress, handlers.CompressHandler(r)))
}

func getTransactionOutput(w http.ResponseWriter, r *http.Request) {
	// Get the transaction...
	txhash := mux.Vars(r)["txhash"]
	tx, err := bitcoind.GetRawTransaction(txhash)
	if err != nil {
		logger.Errorf("Error getting transaction: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}

	// Now look for the output
	vout, err := strconv.Atoi(mux.Vars(r)["vout"])
	if err != nil {
		// This shouldn't happen because the mux routing should not allow non integer characters through.
		logger.Errorf("vout parameter must be a positive integer: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	if vout > len(tx.Vout) {
		logger.Errorf("vout %d does not exist in this transaction", vout)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("vout %d does not exist in this transaction", vout))
		return
	}

	output := tx.Vout[vout]

	res := response{
		TxID:      txhash,
		BlockHash: tx.BlockHash,
		BlockTime: time.Unix(0, int64(tx.Blocktime*1000*int64(time.Millisecond))).Format(time.RFC3339),
		Value:     uint16(output.Value),
		Hex:       output.ScriptPubKey.Hex,
	}

	w.Header().Set("Content-Type", "application/json")
	e := json.NewEncoder(w)
	e.Encode(res)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, fmt.Sprintf("%+v", r))

}
