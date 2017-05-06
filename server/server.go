package server

import (
	"log"
	"net/http"
	"strconv"
)

// Port is the http port the server is started on.
var Port = 6060

// Start starts up the http rest server.
func Start(logger *log.Logger, close chan<- bool) {
	http.HandleFunc("/ready", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(204)
	})
	http.HandleFunc("/close", func(w http.ResponseWriter, req *http.Request) {
		logger.Println("Closing server.")
		close <- true
	})
	startServer(logger)
}

// startServer sets up the HTTP server in a goroutine and waits for it to exit
func startServer(logger *log.Logger) {
	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(Port), nil)
		if err != nil {
			logger.Println(err)
			panic(err)
		}
	}()

	logger.Printf("Server started on port: %v\n", strconv.Itoa(Port))
}
