package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tjheslin1/GoSchedule/database"
)

// SubmitJob represents the path to submit jobs to the database.
type SubmitJob struct {
	urlPath  string
	logger   *log.Logger
	dbClient database.DBClient
}

// SubmitJobRequest represents an incooming request to persist a job.
type SubmitJobRequest struct {
	Name     string      `json:"name"`
	Interval json.Number `json:"interval"`
	URL      string      `json:"url"`
}

// Handler handles request to insert jobs into the database to be watched.
func (submitJob *SubmitJob) Handler(respWriter http.ResponseWriter, req *http.Request) {
	submitJob.logger.Printf("REQUEST: \n%v\n::::::\n", req)

	body, err := ioutil.ReadAll(req.Body)
	if checkErrOccured(err, respWriter, submitJob.logger) {
		return
	}
	defer req.Body.Close()

	var submitJobReq SubmitJobRequest
	err = json.Unmarshal(body, &submitJobReq)
	if checkErrOccured(err, respWriter, submitJob.logger) {
		return
	}

	respWriter.WriteHeader(http.StatusAccepted)
	io.WriteString(respWriter, "Job submitted.")
}

func checkErrOccured(err error, respWriter http.ResponseWriter, logger *log.Logger) bool {
	if err != nil {
		logger.Printf("Error occured.\n'%v\n", err)

		errMsg := fmt.Sprintf("Error occured handling request to submit job:\n%v\n", err)
		http.Error(respWriter, errMsg, http.StatusInternalServerError)

		return true
	}

	return false
}
