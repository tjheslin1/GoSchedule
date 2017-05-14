package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tjheslin1/GoSchedule/database"
	"github.com/tjheslin1/GoSchedule/model"
)

// SubmitJobHandler represents the path to submit jobs to the database.
type SubmitJobHandler struct {
	urlPath  string
	logger   *log.Logger
	dbClient database.DBClient
}

// Handler handles request to insert jobs into the database to be watched.
func (submitJob *SubmitJobHandler) Handler(respWriter http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if checkErrOccured(err, respWriter, submitJob.logger) {
		return
	}

	var submitJobReq model.SubmitJob
	err = json.Unmarshal(body, &submitJobReq)
	if checkErrOccured(err, respWriter, submitJob.logger) {
		return
	}

	submitJobEntry := database.TableEntry{
		Name: "jobs",
		Data: map[string]database.TableCell{
			"name":       database.StringCell{Value: submitJobReq.Name},
			"url":        database.StringCell{Value: submitJobReq.URL},
			"start_time": database.Int64Cell{Value: submitJobReq.StartTime},
			"interval":   database.Int64Cell{Value: submitJobReq.Interval},
		},
	}
	err = submitJob.dbClient.SubmitEntry(submitJobEntry)
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
