package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tjheslin1/GoSchedule/database"
)

// Port is the http port the server is started on.
var Port = 6060

// Start starts up the http rest server.
func Start(logger *log.Logger, close chan<- bool) {
	muxRouter := mux.NewRouter()

	muxRouter.HandleFunc("/ready", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(204)
	}).Methods("GET")

	muxRouter.HandleFunc("/close", func(w http.ResponseWriter, req *http.Request) {
		logger.Println("Closing server.")
		close <- true
	}).Methods("POST")

	dbClient := database.PostgresDBClient{logger}
	submitJob := SubmitJob{"/submit", logger, dbClient}
	muxRouter.HandleFunc(submitJob.urlPath, submitJob.Handler).Methods("POST")

	startServer(muxRouter, logger)
}

// startServer sets up the HTTP server in a goroutine and waits for it to exit
func startServer(handler http.Handler, logger *log.Logger) {
	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(Port), handler)
		if err != nil {
			logger.Println(err)
			panic(err)
		}
	}()

	logger.Printf("Server started on port: %v\n", strconv.Itoa(Port))
}
