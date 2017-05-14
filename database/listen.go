package database

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/lib/pq"
)

// JobListener listens over the postgresql database for INSERTs into the jobs table.
type JobListener struct {
	Logger *log.Logger
}

// Run creates the github.com/lib/pq.Listener to listen on the `watch_tasks`
// channel for new jobs created.
func (lstr *JobListener) Run() {
	listener := pq.NewListener(connectionInfo, 10*time.Second, time.Minute, lstr.reportProblemCallback)
	err := listener.Listen("watch_tasks")
	check(err, lstr.Logger)

	for {
		lstr.waitForNotification(listener)
	}
}

// reportProblemCallback handles errors fatally in the event that the pq.listender,
// returns an error callback.
func (lstr *JobListener) reportProblemCallback(eventType pq.ListenerEventType, err error) {
	if err != nil {
		lstr.Logger.Printf("Error occured for pq.ListenerEventType: '%d'.\n%v\n", eventType, err)
	}
}

// waitForNotification is from:
// http://coussej.github.io/2015/09/15/Listening-to-generic-JSON-notifications-from-PostgreSQL-in-Go/
//
// waitForNotification listens for a notification back from the pq.Listener.
func (lstr *JobListener) waitForNotification(l *pq.Listener) {
	select {
	case n := <-l.Notify:
		lstr.Logger.Println("Received data from channel [", n.Channel, "] :")

		var prettyJSON bytes.Buffer
		err := json.Indent(&prettyJSON, []byte(n.Extra), "", "\t")
		if err != nil {
			lstr.Logger.Println("Error processing JSON: ", err)
			return
		}
		lstr.Logger.Println(string(prettyJSON.Bytes()))
		return
	case <-time.After(90 * time.Second):
		lstr.Logger.Println("Received no events for 90 seconds, checking connection")
		go func() {
			l.Ping()
		}()
		return
	}
}
