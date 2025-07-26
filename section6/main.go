package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"mycoolserver/internal/users"
	"net/http"
)

type UserData struct {
	FirstName string
	LastName  string
	Email     string
}

type server struct {
	userManager *users.Manager
}

func main() {
	manager := users.NewManager()
	s := server{
		userManager: manager,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/{$}", handleRoot)
	mux.HandleFunc("/goodbye/", handleGoodbye)
	mux.HandleFunc("/hello/", handleHelloParameterized)
	mux.HandleFunc("/responses/{user}/hello/", handleUserResponsesHello)
	mux.HandleFunc("POST /user/hello/", s.handleHelloHeader)
	mux.HandleFunc("POST /json", handleJSON)
	mux.HandleFunc("POST /add-user", s.addUser)
	mux.HandleFunc("POST /get-user", s.getUser)

	// an example of handling the endpoint with a closure
	// mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) { exampleAddUser(w, r, manager) })

	// an example of a closure with switch on the method
	//mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
	//	switch r.Method {
	//	case http.MethodPost:
	//		exampleAddUser(w, r, manager)
	//	case http.MethodGet:
	//	// Do things that you do on GET
	//	default:
	//		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	//	}
	//})

	// an example of handling the endpoint with a handler
	//mux.Handle("/user", exampleAddUser(manager))

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func (s *server) handleHelloHeader(w http.ResponseWriter, r *http.Request) {
	firstName := r.Header.Get("userFirst")
	if firstName == "" {
		http.Error(w, "invalid first name provided", http.StatusBadRequest)
		return
	}

	lastName := r.Header.Get("userLast")
	if lastName == "" {
		http.Error(w, "invalid last name provided", http.StatusBadRequest)
		return
	}
	user, err := s.userManager.GetUserByName(firstName, lastName)
	if err != nil {
		if errors.Is(err, users.ErrNoResultsFound) {
			http.Error(w, "no users found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("error retrieving user: %v\n", err), http.StatusInternalServerError)
		}
		return
	}

	result := fmt.Sprintf("Hello, %s %s!  Your email is: %s\n", user.FirstName, user.LastName, user.Email.Address)

	_, err = w.Write([]byte(result))
	if err != nil {
		slog.Error("error writing response body", "err", err)
		return
	}
}

func (s *server) addUser(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, fmt.Sprintf("unsupported Content-Type header: %q", contentType), http.StatusUnsupportedMediaType)
		return
	}

	// limit to 1MB
	requestBody := http.MaxBytesReader(w, r.Body, 1048576)

	decoder := json.NewDecoder(requestBody)
	decoder.DisallowUnknownFields()

	var u UserData

	err := decoder.Decode(&u)
	if err != nil {
		slog.Error("error decoding addUser request body", "err", err)
		http.Error(w, "bad request body", http.StatusBadRequest)
		return
	}

	err = s.userManager.AddUser(u.FirstName, u.LastName, u.Email)
	if err != nil {
		http.Error(w, fmt.Sprintf("error adding user: %v\n", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *server) getUser(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, fmt.Sprintf("unsupported Content-Type header %q", contentType), http.StatusUnsupportedMediaType)
		return
	}

	// limit to 1MB
	requestBody := http.MaxBytesReader(w, r.Body, 1048576)

	decoder := json.NewDecoder(requestBody)
	decoder.DisallowUnknownFields()

	var u UserData

	err := decoder.Decode(&u)
	if err != nil {
		http.Error(w, fmt.Sprintf("error decoding request body: %v\n", err), http.StatusBadRequest)
		return
	}

	user, err := s.userManager.GetUserByName(u.FirstName, u.LastName)
	if err != nil {
		if errors.Is(err, users.ErrNoResultsFound) {
			http.Error(w, "no users found", http.StatusNotFound)
		} else {
			slog.Error("error retrieving user", "err", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	converted := convertUserToUserData(user)

	marshalled, err := json.Marshal(converted)
	if err != nil {
		slog.Error("error marshalling getUser response", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(marshalled)
	if err != nil {
		// headers are set by write call, best we can do is log an error
		slog.Error("error writing getUser response body", "err", err)
	}

	return
}

func handleRoot(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Welcome to our homepage!\n"))
	if err != nil {
		slog.Error("error writing response", "err", err)
		return
	}

	return
}

func handleGoodbye(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Goodbye!\n"))
	if err != nil {
		slog.Error("error writing response", "err", err)
		return
	}

	return
}

func handleHelloParameterized(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	username := "User"
	userList := params["user"]
	if len(userList) > 0 {
		username = userList[0]
	}

	handleHello(w, username)
}

func handleUserResponsesHello(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("user")

	handleHello(w, username)
}

func handleJSON(w http.ResponseWriter, r *http.Request) {
	byteData, err := io.ReadAll(r.Body)
	if err != nil || len(byteData) < 1 {
		slog.Error("error reading request body", "err", err)
		http.Error(w, "bad request body", http.StatusBadRequest)
		return
	}

	var reqData UserData
	err = json.Unmarshal(byteData, &reqData)
	if err != nil {
		slog.Error("error unmarshalling request body", "err", err)
		http.Error(w, "error parsing request JSON", http.StatusBadRequest)
		return
	}

	if reqData.FirstName == "" {
		http.Error(w, "invalid username provided", http.StatusBadRequest)
		return
	}

	handleHello(w, reqData.FirstName)
}

func handleHello(w http.ResponseWriter, username string) {
	var output bytes.Buffer
	output.WriteString("Hello, ")
	output.WriteString(username)
	output.WriteString("!\n")

	_, err := w.Write(output.Bytes())
	if err != nil {
		slog.Error("error writing response body", "err", err)
		return
	}
}

func convertUserToUserData(u *users.User) *UserData {
	converted := UserData{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email.Address,
	}

	return &converted
}

func validateContentType(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, fmt.Sprintf("unsupported Content-Type header %q", contentType), http.StatusUnsupportedMediaType)
		return
	}
}

// example for use with a closure
//func exampleAddUser(w http.ResponseWriter, r *http.Request, manager *users.Manager) {
//	http.Error(w, "not implemented", http.StatusInternalServerError)
//}

// example for use with a handler
//func exampleAddUser(manager *users.Manager) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		http.Error(w, "not implemented", http.StatusInternalServerError)
//	}
//}
