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

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/ordishs/gocore"

	"github.com/gazpatchodev/chainservice/cache"
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
	r.HandleFunc("/api/v1/bsv/{txhash}/{vout:[0-9]+}", getTransactionOutputs).Methods("GET")
	r.Handle("/", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(errorHandler)))

	// Wrap our server with our gzip handler to gzip compress all responses.
	log.Fatal(http.ListenAndServe(listenAddress, handlers.CompressHandler(r)))
}

func getTransactionOutputs(w http.ResponseWriter, r *http.Request) {

	txhash := mux.Vars(r)["txhash"]

	opr, err := cache.GetOPReturnData(txhash)
	if err != nil {
		logger.Errorf("Error getting transaction: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}

	voutStr, ok := mux.Vars(r)["vout"]

	if ok {
		voutRequested, err := strconv.Atoi(voutStr)
		if err != nil {
			// This shouldn't happen because the mux routing should not allow non integer characters through.
			logger.Errorf("vout parameter must be a positive integer: %+v", err)
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, err.Error())
			return
		}

		// Go through the opr and remove all vouts that are not required
		i := 0 // output index
		for _, vout := range opr.Vouts {
			if vout.N == voutRequested {
				// copy and increment index
				opr.Vouts[i] = vout
				i++
			}
		}

		if i == 0 {
			// No matching vouts...
			m := fmt.Sprintf("vout %d does not exist in transaction", voutRequested)
			logger.Error(m)
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, m)
			return
		}

		opr.Vouts = opr.Vouts[:i]
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
