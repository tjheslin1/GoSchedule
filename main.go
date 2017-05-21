package main

import (
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/tjheslin1/GoSchedule/clock"
	"github.com/tjheslin1/GoSchedule/database"
	"github.com/tjheslin1/GoSchedule/model"
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

	jobsSubmitted := make(chan model.SubmitJob)
	jobListener := database.JobListener{
		JobSubmittedReceiver: jobsSubmitted,
		ConnectionCheckTime:  60,
		Logger:               logger,
	}

	go jobListener.Run()

	for {
		select {
		case job := <-jobsSubmitted:
			logger.Printf("Job has been submitted to database.\n"+
				"Can now spawn a job runner to listen to changes to this row.\n%v\n", job)
			jobRunner := model.JobRunner{Job: job, Logger: logger}

			clck := clock.Clock{Now: time.Now()}
			go jobRunner.Start(&clck)
		case <-quit:
			logger.Println("Closing.")
			os.Exit(0)
		}
	}
}
