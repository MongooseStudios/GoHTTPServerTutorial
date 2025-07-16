package main

import (
	"bytes"
	"log"
	"log/slog"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/{$}", handleRoot)
	mux.HandleFunc("/goodbye/", handleGoodbye)
	mux.HandleFunc("/hello/", handleHelloParameterized)
	mux.HandleFunc("/responses/{user}/hello/", handleUserResponsesHello)
	mux.HandleFunc("/user/hello/", handleHelloHeader)

	log.Fatal(http.ListenAndServe(":8080", mux))
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

func handleHelloHeader(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("user")
	if username == "" {
		http.Error(w, "invalid username provided", http.StatusBadRequest)
		return
	}

	handleHello(w, username)
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

// These are the original functions before the refactor
//func handleHelloParameterized(w http.ResponseWriter, r *http.Request) {
//	params := r.URL.Query()
//	username := "User"
//	userList := params["user"]
//	if len(userList) > 0 {
//		username = userList[0]
//	}
//	var output bytes.Buffer
//	output.WriteString("Hello, ")
//	output.WriteString(username)
//	output.WriteString("!\n")
//
//	_, err := w.Write(output.Bytes())
//	if err != nil {
//		slog.Error("error writing response", "err", err)
//		return
//	}
//}
//
//func handleUserResponsesHello(w http.ResponseWriter, r *http.Request) {
//	username := r.PathValue("user")
//
//	var output bytes.Buffer
//	output.WriteString("Hello, ")
//	output.WriteString(username)
//	output.WriteString("!\n")
//
//	_, err := w.Write(output.Bytes())
//	if err != nil {
//		slog.Error("error writing response", "err", err)
//		return
//	}
//}
//
//func handleHelloHeader(w http.ResponseWriter, r *http.Request) {
//	username := r.Header.Get("user")
//
//	var output bytes.Buffer
//	output.WriteString("Hello, ")
//	output.WriteString(username)
//	output.WriteString("!\n")
//
//	_, err := w.Write(output.Bytes())
//	if err != nil {
//		slog.Error("error writing response", "err", err)
//		return
//	}
//}
