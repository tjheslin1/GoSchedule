package main

import (
	"log"
	"os"

	"github.com/tjheslin1/GoSchedule/server"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println("GoSchedule is running!")

	close := make(chan bool)
	server.Start(logger, close)

	<-close
}
