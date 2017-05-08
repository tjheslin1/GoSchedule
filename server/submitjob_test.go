package server

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSubmitJobHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "http://localhost:6060/submit",
		bytes.NewBufferString(`{"name":"testJob", "interval":1000, "url":"http:localhost:6060/ready"}`))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	respWriter := httptest.NewRecorder()

	logBuffer := new(bytes.Buffer)
	logger := log.New(logBuffer, "", 0)
	submitJob := SubmitJob{"/submit", logger, nil} // TODO

	handler := http.HandlerFunc(submitJob.Handler)
	handler.ServeHTTP(respWriter, req)

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

	logBuffer := new(bytes.Buffer)
	logger := log.New(logBuffer, "", 0)
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
