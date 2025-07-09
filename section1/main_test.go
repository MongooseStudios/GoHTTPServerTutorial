package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleHello(t *testing.T) {
	// we call this w because it will take the place of the http.ResponseWriter which is conventionally set to w
	w := httptest.NewRecorder()

	handleHello(w, nil)

	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v but got %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Hello, World!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got: %q, expected %q", w.Body.Bytes(), expectedMessage)
	}
}

func TestHandleGoodbye(t *testing.T) {
	// we call this w because it will take the place of the http.ResponseWriter which is conventionally set to w
	w := httptest.NewRecorder()

	handleGoodbye(w, nil)

	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v but got %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}
	
	expectedMessage := []byte("Goodbye!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got: %q, expected %q", w.Body.Bytes(), expectedMessage)
	}
}
