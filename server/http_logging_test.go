package server

import (
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestWriterHeader(t *testing.T) {
	recordedRespWriter := httptest.NewRecorder()
	loggingRespWriter := loggingResponseWriter{
		ResponseWriter: recordedRespWriter,
	}

	loggingRespWriter.WriteHeader(100)

	if loggingRespWriter.status != 100 {
		t.Errorf("Expected status to be stored as\n'%v'\nbut was\n'%v'\n", 100, loggingRespWriter.status)
	}

	if recordedRespWriter.Code != 100 {
		t.Errorf("Expected status to be written to responseWriter as\n'%v'\nbut was\n'%v'\n", 100, recordedRespWriter.Code)
	}
}

func TestWrite(t *testing.T) {
	recordedRespWriter := httptest.NewRecorder()
	loggingRespWriter := loggingResponseWriter{
		ResponseWriter: recordedRespWriter,
	}

	expectedBody := []byte("Some content!")
	loggingRespWriter.Write(expectedBody)

	if !reflect.DeepEqual(loggingRespWriter.body, expectedBody) {
		t.Errorf("Expected status to be stored as\n'%v'\nbut was\n'%v'\n", string(expectedBody), string(loggingRespWriter.body))
	}

	if !reflect.DeepEqual(recordedRespWriter.Body.Bytes(), expectedBody) {
		t.Errorf("Expected status to be written to responseWriter as\n'%v'\nbut was\n'%v'\n", string(expectedBody), string(recordedRespWriter.Body.Bytes()))
	}
}
