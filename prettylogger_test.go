package prettylogger

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func successHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func badRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

func internalErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func TestSingleOutputSuccess(t *testing.T) {
	var b bytes.Buffer
	p := NewPrettyLogger(&b)
	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	fn := p.PrettyLoggerMiddleWare(successHandler)
	fn.ServeHTTP(w, r)

	logOutput := "GET /test " + green + "200" + reset
	if !strings.Contains(b.String(), logOutput) {
		t.Errorf("incorrect log output. Expected '%s' Recieved '%s'", logOutput, b.String())
	}
}

func TestSingleOutputBadRequest(t *testing.T) {
	var b bytes.Buffer
	p := NewPrettyLogger(&b)
	r := httptest.NewRequest(http.MethodPost, "/test", nil)
	w := httptest.NewRecorder()

	fn := p.PrettyLoggerMiddleWare(badRequestHandler)
	fn.ServeHTTP(w, r)

	logOutput := "POST /test " + yellow + "400" + reset
	if !strings.Contains(b.String(), logOutput) {
		t.Errorf("incorrect log output. Expected '%s' Recieved '%s'", logOutput, b.String())
	}
}

func TestSingleOutputInternalError(t *testing.T) {
	var b bytes.Buffer
	p := NewPrettyLogger(&b)
	r := httptest.NewRequest(http.MethodGet, "/anothertest", nil)
	w := httptest.NewRecorder()

	fn := p.PrettyLoggerMiddleWare(internalErrorHandler)
	fn.ServeHTTP(w, r)

	logOutput := "GET /anothertest " + red + "500" + reset
	if !strings.Contains(b.String(), logOutput) {
		t.Errorf("incorrect log output. Expected '%s' Recieved '%s'", logOutput, b.String())
	}
}

func TestDoubleOutputSuccess(t *testing.T) {
	var b1 bytes.Buffer
	var b2 bytes.Buffer
	p := NewPrettyLogger(&b1, &b2)
	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	fn := p.PrettyLoggerMiddleWare(successHandler)
	fn.ServeHTTP(w, r)

	logOutput := "GET /test " + green + "200" + reset
	if !strings.Contains(b1.String(), logOutput) {
		t.Errorf("incorrect log output for log 1. Expected '%s' Recieved '%s'", logOutput, b1.String())
	}

	if !strings.Contains(b2.String(), logOutput) {
		t.Errorf("incorrect log output for log 2. Expected '%s' Recieved '%s'", logOutput, b2.String())
	}
}
