package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleRoot(t *testing.T) {
	// we call this w because it will take the place of the http.ResponseWriter which is conventionally set to w
	w := httptest.NewRecorder()

	handleRoot(w, nil)

	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected: %v but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Welcome to our homepage!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got: %q, expected %q", w.Body.String(), expectedMessage)
	}
}

func TestHandleGoodbye(t *testing.T) {
	// we call this w because it will take the place of the http.ResponseWriter which is conventionally set to w
	w := httptest.NewRecorder()

	handleGoodbye(w, nil)

	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected: %v but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Goodbye!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got: %q, expected %q", w.Body.String(), expectedMessage)
	}
}

func TestHandleHelloParameterized(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello?user=TestMan", nil)

	// we call this w because it will take the place of the http.ResponseWriter which is conventionally set to w
	w := httptest.NewRecorder()

	handleHelloParameterized(w, req)

	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected: %v but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Hello, TestMan!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got: %q, expected %q", w.Body.String(), expectedMessage)
	}
}

func TestHandleHelloParameterizedNoParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello/", nil)

	// we call this w because it will take the place of the http.ResponseWriter which is conventionally set to w
	w := httptest.NewRecorder()

	handleHelloParameterized(w, req)

	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected: %v but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Hello, User!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got: %q, expected %q", w.Body.String(), expectedMessage)
	}
}

func TestHandleHelloParameterizedWrongParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello?foo=bar", nil)

	// we call this w because it will take the place of the http.ResponseWriter which is conventionally set to w
	w := httptest.NewRecorder()

	handleHelloParameterized(w, req)

	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected: %v but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Hello, User!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got: %q, expected %q", w.Body.String(), expectedMessage)
	}
}

func TestHandleUserResponsesHello(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/responses/TestMan/hello/", nil)
	req.SetPathValue("user", "TestMan")

	// we call this w because it will take the place of the http.ResponseWriter which is conventionally set to w
	w := httptest.NewRecorder()

	handleUserResponsesHello(w, req)

	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected: %v but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Hello, TestMan!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got: %q, expected %q", w.Body.String(), expectedMessage)
	}
}

func TestHelloHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user/hello/", nil)
	req.Header.Set("user", "Test Man")

	// we call this w because it will take the place of the http.ResponseWriter which is conventionally set to w
	w := httptest.NewRecorder()

	handleHelloHeader(w, req)

	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected: %v but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Hello, Test Man!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got: %q, expected %q", w.Body.String(), expectedMessage)
	}
}

func TestHelloHeaderNoHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user/hello/", nil)

	// we call this w because it will take the place of the http.ResponseWriter which is conventionally set to w
	w := httptest.NewRecorder()

	handleHelloHeader(w, req)

	desiredCode := http.StatusBadRequest
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected: %v but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("invalid username provided\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got: %q, expected %q", w.Body.String(), expectedMessage)
	}
}

func TestHandleJSON(t *testing.T) {
	testRequest := UserData{
		Name: "TestMan",
	}

	marshalledRequestBody, err := json.Marshal(testRequest)
	if err != nil {
		t.Fatalf("error marshalling test data: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/json", bytes.NewBuffer(marshalledRequestBody))

	// we call this w because it will take the place of the http.ResponseWriter which is conventionally set to w
	w := httptest.NewRecorder()

	handleJSON(w, req)

	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected: %v but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Hello, TestMan!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got: %q, expected %q", w.Body.String(), expectedMessage)
	}
}

func TestHandleJSONEmptyBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/json", nil)

	// we call this w because it will take the place of the http.ResponseWriter which is conventionally set to w
	w := httptest.NewRecorder()

	handleJSON(w, req)

	desiredCode := http.StatusBadRequest
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected: %v but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("bad request body\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got: %q, expected %q", w.Body.String(), expectedMessage)
	}
}

func TestHandleJSONEmptyNameField(t *testing.T) {
	testRequest := UserData{
		Name: "",
	}

	marshalledRequestBody, err := json.Marshal(testRequest)
	if err != nil {
		t.Fatalf("error marshalling test data: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/json", bytes.NewBuffer(marshalledRequestBody))

	// we call this w because it will take the place of the http.ResponseWriter which is conventionally set to w
	w := httptest.NewRecorder()

	handleJSON(w, req)

	desiredCode := http.StatusBadRequest
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected: %v but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("invalid username provided\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got: %q, expected %q", w.Body.String(), expectedMessage)
	}
}
