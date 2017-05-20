package database

import (
	"encoding/json"
	"log"
	"time"

	"github.com/lib/pq"
	"github.com/tjheslin1/GoSchedule/model"
)

// JobListener listens over the postgresql database for INSERTs into the jobs table.
type JobListener struct {
	JobSubmitted        chan<- model.SubmitJob
	ConnectionCheckTime int
	Logger              *log.Logger
}

// Run creates the github.com/lib/pq.Listener to listen on the `watch_tasks`
// channel for new jobs created.
func (lstr *JobListener) Run() {
	listener := pq.NewListener(connectionInfo, 10*time.Second, time.Minute, lstr.reportProblemCallback)
	defer listener.Close()

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
	case notification := <-l.Notify:
		passNotitification(notification, lstr)
		return
	case <-time.After(time.Duration(lstr.ConnectionCheckTime) * time.Second):
		lstr.Logger.Printf("Received no events for %d seconds, checking connection\n", lstr.ConnectionCheckTime)
		go func() {
			l.Ping()
		}()
		return
	}
}

func passNotitification(notification *pq.Notification, lstr *JobListener) {
	submittedJob := lstr.unmarshallJSONJobNotification([]byte(notification.Extra))
	lstr.JobSubmitted <- submittedJob
}

func (lstr *JobListener) unmarshallJSONJobNotification(jsonData []byte) model.SubmitJob {
	var notification JobWatchNotification
	err := json.Unmarshal(jsonData, &notification)
	check(err, lstr.Logger)

	return notification.Data
}

// JobWatchNotification represents the JSON returned by the postgresql notificaiton.
// The struct and its fields are expored to conform to the Go standards when using
// the JSON tag.
type JobWatchNotification struct {
	Table  string          `json:"table"`
	Action string          `json:"action"`
	Data   model.SubmitJob `json:"data"`
}
