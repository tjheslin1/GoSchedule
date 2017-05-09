package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/tjheslin1/GoSchedule/database"
)

func TestSubmitJobHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "http://localhost:6060/submit",
		bytes.NewBufferString(`{"name":"testJob", "interval":1000, "url":"http:localhost:6060/ready"}`))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	respWriter := httptest.NewRecorder()

	logger := log.New(new(bytes.Buffer), "", 0)
	dummyDBClient := DummyDBClient{}
	submitJob := SubmitJob{"/submit", logger, &dummyDBClient}

	handler := http.HandlerFunc(submitJob.Handler)
	handler.ServeHTTP(respWriter, req)

	var expectedTableEntry = database.TableEntry{
		Name: "jobs",
		Data: map[string]interface{}{
			"job_id":   1,
			"name":     "testJob",
			"url":      "http:localhost:6060/ready",
			"interval": json.Number("1000"),
		},
	}
	if !reflect.DeepEqual(dummyDBClient.capturedTableEntry, expectedTableEntry) {
		t.Errorf("Expected\n'%v'\nto equal\n'%v'\n", dummyDBClient.capturedTableEntry, expectedTableEntry)
	}

	if status := respWriter.Code; status != http.StatusAccepted {
		t.Errorf("Handler returned wrong status code, wanted:\n'%v'\nbut got:\n'%v\n",
			http.StatusAccepted, status)
	}

	expected := "Job submitted."
	if respWriter.Body.String() != expected {
		t.Errorf("Handler returned unexpected body, wanted:\n'%v'\nbut got:\n'%v'\n",
			expected, respWriter.Body.String())
	}
}

func TestSubmitJobHandlerBadRequest(t *testing.T) {
	req, err := http.NewRequest("POST", "http://localhost:6060/submit",
		bytes.NewBufferString(`{"name":"testJob", "interval":1000`))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	respWriter := httptest.NewRecorder()

	logger := log.New(new(bytes.Buffer), "", 0)
	submitJob := SubmitJob{"/submit", logger, nil}

	handler := http.HandlerFunc(submitJob.Handler)
	handler.ServeHTTP(respWriter, req)

	if status := respWriter.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code, wanted:\n'%v'\nbut got:\n'%v\n",
			http.StatusInternalServerError, status)
	}

	expected := "Error occured handling request to submit job:\nunexpected end of JSON input\n\n"
	if respWriter.Body.String() != expected {
		t.Errorf("Handler returned unexpected body, wanted:\n'%v'\nbut got:\n'%v'\n",
			expected, respWriter.Body.String())
	}
}

type DummyDBClient struct {
	capturedTableEntry database.TableEntry
}

func (dummyClient *DummyDBClient) Connection() *sql.DB {
	return nil
}

func (dummyClient *DummyDBClient) SubmitJob(jobEntry database.TableEntry) error {
	dummyClient.capturedTableEntry = jobEntry
	return nil
}
