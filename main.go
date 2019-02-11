package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/ordishs/gocore"

	"./cache"
)

var (
	logger           = gocore.Log("ChainService")
	listenAddress, _ = gocore.Config().Get("listenAddress")
)

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

	r.HandleFunc("/api/v1/bsv/{txhash}", getTransactionOutputs).Methods("GET")
	// r.HandleFunc("/api/v1/bsv/{txhash}/{vout:[0-9]+}", getTransactionOutput).Methods("GET")
	r.Handle("/", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(errorHandler)))

	// Wrap our server with our gzip handler to gzip compress all responses.
	log.Fatal(http.ListenAndServe(listenAddress, handlers.CompressHandler(r)))
}

func getTransactionOutputs(w http.ResponseWriter, r *http.Request) {

	txhash := mux.Vars(r)["txhash"]

	// voutStr, ok := mux.Vars(r)["vout"]
	// vout := -1
	// var err error

	// if ok {
	// 	vout, err = strconv.Atoi(voutStr)
	// 	if err != nil {
	// 		// This shouldn't happen because the mux routing should not allow non integer characters through.
	// 		logger.Errorf("vout parameter must be a positive integer: %+v", err)
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		io.WriteString(w, err.Error())
	// 		return
	// 	}
	// }

	opr, err := cache.GetOPReturnData(txhash) //, int16(vout))
	if err != nil {
		logger.Errorf("Error getting transaction: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	e := json.NewEncoder(w)

	e.Encode(opr)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, "Not found.")
	// io.WriteString(w, "fmt.Sprintf("%+v", r)")
}
