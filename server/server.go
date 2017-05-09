package server

import (
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tjheslin1/GoSchedule/database"
)

// Port is the http port the server is started on.
var Port = 6060

// Start starts up the http rest server.
func Start(logger *log.Logger, quit chan<- bool) {
	muxRouter := mux.NewRouter()

	var readyHandler http.HandlerFunc = func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(204)
	}
	muxRouter.Handle("/ready", logRequestResponse(logger, readyHandler)).Methods("GET")

	var closeHandler http.HandlerFunc = func(w http.ResponseWriter, req *http.Request) {
		logger.Println("Closing server.")
		quit <- true
	}
	muxRouter.Handle("/close", logRequestResponse(logger, closeHandler)).Methods("POST")

	dbClient := database.PostgresDBClient{Logger: logger}
	submitJob := SubmitJob{"/submit", logger, &dbClient}
	var submitJobHandler http.HandlerFunc = submitJob.Handler
	muxRouter.Handle(submitJob.urlPath, logRequestResponse(logger, submitJobHandler)).Methods("POST")

	startServer(muxRouter, logger)
}

func logRequestResponse(logger *log.Logger, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, req *http.Request) {
		loggingRespWriter := &loggingResponseWriter{
			ResponseWriter: respWriter,
		}
		requestDump, err := httputil.DumpRequest(req, true)
		check(err, logger)

		logger.Printf("REQUEST:\n%v\n::::::\n", string(requestDump))
		handler.ServeHTTP(loggingRespWriter, req)
		logger.Printf("RESPONSE:\n%d\n%s\n::::::\n", loggingRespWriter.status, string(loggingRespWriter.body))
	})
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

func check(err error, logger *log.Logger) {
	if err != nil {
		logger.Printf("Error occured handling request:\n'%v", err)
	}
}
