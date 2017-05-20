package main

import (
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/tjheslin1/GoSchedule/database"
	"github.com/tjheslin1/GoSchedule/server"
)

// http://coussej.github.io/2015/09/15/Listening-to-generic-JSON-notifications-from-PostgreSQL-in-Go/
func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println("GoSchedule is running!")

	quit := make(chan bool)

	go server.Start(logger, quit)

	dbClient := database.PostgresDBClient{Logger: logger}
	database.SetUpSchema(&dbClient, logger)

	taskListener := database.JobListener{Logger: logger}
	go taskListener.Run()

	<-quit
}
