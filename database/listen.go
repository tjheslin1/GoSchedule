package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
)

type taskListener struct {
	logger   *log.Logger
	dbClient DBClient
}

func (lstr *taskListener) run() {
	listener := pq.NewListener(connectionInfo, 10*time.Second, time.Minute, lstr.reportProblemCallback)
	err := listener.Listen("watch_tasks")
	check(err, lstr.logger)

	for {
		lstr.watch()
		lstr.waitForNotification(listener)
	}
}

func (lstr *taskListener) watch() {
	for {
		var work sql.NullInt64
		err := lstr.dbClient.Connection().QueryRow("SELECT watch_tasks()").Scan(&work)
		if err != nil {
			lstr.logger.Println("call to watch_tasks() failed: ", err)
			time.Sleep(10 * time.Second)
			continue
		}
		if !work.Valid {
			// no more work to do
			lstr.logger.Println("ran out of work")
			return
		}
		lstr.logger.Println("starting work on ", work.Int64)
		// go processTask()
	}
}

func (lstr *taskListener) processTask() {

}

func (lstr *taskListener) reportProblemCallback(eventType pq.ListenerEventType, err error) {
	if err != nil {
		lstr.logger.Fatal(err.Error())
	}
}

// waitForNotification is from: https://github.com/lib/pq/blob/master/listen_example/doc.go
func (lstr *taskListener) waitForNotification(l *pq.Listener) {
	select {
	case <-l.Notify:
		fmt.Println("received notification, new work available")
	case <-time.After(90 * time.Second):
		go l.Ping()
		// Check if there's more work available, just in case it takes
		// a while for the Listener to notice connection loss and
		// reconnect.
		fmt.Println("received no work for 90 seconds, checking for new work")
	}
}
