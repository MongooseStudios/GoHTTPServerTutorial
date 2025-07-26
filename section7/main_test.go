package main

import (
	"bytes"
	"encoding/json"
	"mycoolserver/internal/users"
	"net/http"
	"net/http/httptest"
	"net/mail"
	"reflect"
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
	req.Header.Set("userFirst", "Test")
	req.Header.Set("userLast", "Man")

	// we call this w because it will take the place of the http.ResponseWriter which is conventionally set to w
	w := httptest.NewRecorder()

	testManager := users.NewManager()
	err := testManager.AddUser("Test", "Man", "testman@example.com")
	if err != nil {
		t.Fatalf("error creating test user: %v", err)
	}

	testServer := server{
		userManager: testManager,
	}

	testServer.handleHelloHeader(w, req)

	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected: %v but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Hello, Test Man!  Your email is: testman@example.com\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got: %q, expected %q", w.Body.String(), expectedMessage)
	}
}

func TestHelloHeaderNoHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user/hello/", nil)

	// we call this w because it will take the place of the http.ResponseWriter which is conventionally set to w
	w := httptest.NewRecorder()

	testServer := server{}

	testServer.handleHelloHeader(w, req)

	desiredCode := http.StatusBadRequest
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected: %v but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("invalid first name provided\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got: %q, expected %q", w.Body.String(), expectedMessage)
	}
}

func TestHandleJSON(t *testing.T) {
	testRequest := UserData{
		FirstName: "TestMan",
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
		FirstName: "",
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

func TestAddUser(t *testing.T) {
	testUser := UserData{
		FirstName: "Test",
		LastName:  "Man",
		Email:     "TestMan@example.com",
	}

	marshalledRequestBody, err := json.Marshal(testUser)
	if err != nil {
		t.Fatalf("error marshalling test data: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(marshalledRequestBody))
	req.Header.Set("Content-Type", "application/json")

	// we call this w because it will take the place of the http.ResponseWriter which is conventionally set to w
	w := httptest.NewRecorder()

	testManager := users.NewManager()
	testServer := server{
		userManager: testManager,
	}

	testServer.addUser(w, req)

	desiredCode := http.StatusCreated
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected: %v but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	resultUser, err := testManager.GetUserByName(testUser.FirstName, testUser.LastName)
	if err != nil {
		t.Fatalf("error getting test user back out of manager: %v", err)
	}

	// convert to UserData so we can compare
	convertedResult := convertUserToUserData(resultUser)
	if !reflect.DeepEqual(&testUser, convertedResult) {
		t.Errorf("bad retrieved user\nwanted: %v\ngot: %v\n", &testUser, convertedResult)
	}

}

func TestAddUserBadHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/user", nil)

	// we call this w because it will take the place of the http.ResponseWriter which is conventionally set to w
	w := httptest.NewRecorder()

	testManager := users.NewManager()
	testServer := server{
		userManager: testManager,
	}

	testServer.addUser(w, req)

	desiredCode := http.StatusUnsupportedMediaType
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected: %v but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}
}

func TestGetUser(t *testing.T) {
	testFirstName := "Test"
	testLastName := "Man"
	testEmail := "TestMan@example.com"

	testManager := users.NewManager()
	testServer := server{
		userManager: testManager,
	}

	err := testManager.AddUser(testFirstName, testLastName, testEmail)
	if err != nil {
		t.Fatalf("error inserting test user: %v", err)
	}

	testQuery := UserData{
		FirstName: testFirstName,
		LastName:  testLastName,
	}

	marshalledRequestBody, err := json.Marshal(testQuery)
	if err != nil {
		t.Fatalf("error marshalling test data: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(marshalledRequestBody))
	req.Header.Set("Content-Type", "application/json")

	// we call this w because it will take the place of the http.ResponseWriter which is conventionally set to w
	w := httptest.NewRecorder()

	testServer.getUser(w, req)

	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected: %v but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	decoder := json.NewDecoder(w.Body)

	var decodedResult UserData
	err = decoder.Decode(&decodedResult)
	if err != nil {
		t.Fatalf("error decoding response body: %v", err)
	}

	expectedData := UserData{
		FirstName: testFirstName,
		LastName:  testLastName,
		Email:     testEmail,
	}

	if !reflect.DeepEqual(decodedResult, expectedData) {
		t.Errorf("bad result\ngot: %v\nwanted: %v\n", decodedResult, expectedData)
	}

}

func TestGetUserBadHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/user", nil)

	// we call this w because it will take the place of the http.ResponseWriter which is conventionally set to w
	w := httptest.NewRecorder()

	testManager := users.NewManager()
	testServer := server{
		userManager: testManager,
	}

	testServer.addUser(w, req)

	desiredCode := http.StatusUnsupportedMediaType
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected: %v but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedBody := []byte("unsupported Content-Type header: \"\"\n")
	if !bytes.Equal(expectedBody, w.Body.Bytes()) {
		t.Errorf("bad response body, should be: %q but got: %q", expectedBody, w.Body.String())
	}
}

func TestGetUserNoUser(t *testing.T) {
	testFirstName := "Test"
	testLastName := "Man"
	testEmail := "TestMan@example.com"

	testManager := users.NewManager()
	testServer := server{
		userManager: testManager,
	}

	err := testManager.AddUser(testFirstName, testLastName, testEmail)
	if err != nil {
		t.Fatalf("error inserting test user: %v", err)
	}

	testQuery := UserData{
		FirstName: "foo",
		LastName:  "bar",
	}

	marshalledRequestBody, err := json.Marshal(testQuery)
	if err != nil {
		t.Fatalf("error marshalling test data: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(marshalledRequestBody))
	req.Header.Set("Content-Type", "application/json")

	// we call this w because it will take the place of the http.ResponseWriter which is conventionally set to w
	w := httptest.NewRecorder()

	testServer.getUser(w, req)

	desiredCode := http.StatusNotFound
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected: %v but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedBody := "no users found\n"
	if w.Body.String() != expectedBody {
		t.Errorf("bad response body, should be %q but got %q", expectedBody, w.Body.String())
	}

}

func TestConvertUserToUserData(t *testing.T) {
	testFirstName := "Test"
	testLastName := "User"
	testEmail, err := mail.ParseAddress("testuser@example.com")
	if err != nil {
		t.Fatalf("error parsing test email")
	}

	testUser := users.User{
		FirstName: testFirstName,
		LastName:  testLastName,
		Email:     *testEmail,
	}

	result := convertUserToUserData(&testUser)

	ExpectedUser := &UserData{
		FirstName: testFirstName,
		LastName:  testLastName,
		Email:     testEmail.Address,
	}

	if !reflect.DeepEqual(ExpectedUser, result) {
		t.Errorf("bad conversion\nwant: %+v\ngot: %+v\n", ExpectedUser, result)
	}
}
