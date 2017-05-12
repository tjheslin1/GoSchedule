package database

import (
	"log"
	"time"

	"github.com/lib/pq"
)

type taskListener struct {
	logger *log.Logger
}

func (lstr *taskListener) listenStatement() string {
	listener := pq.NewListener(connectionInfo, 10*time.Second, time.Minute, lstr.listenCallback)
	err := listener.Listen("watch_jobs")
	check(err, lstr.logger)

	return ""
}

func (lstr *taskListener) listenCallback(eventType pq.ListenerEventType, err error) {
	if err != nil {
		lstr.logger.Fatal(err.Error())
	}
}
