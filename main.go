package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s", r.UserAgent())
}

func main() {
	r := mux.NewRouter()

	r.Handle("/", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(rootHandler))).Methods("GET")
	r.Handle("/ping", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(pingHandler))).Methods("GET")

	http.ListenAndServe(":8080", r)
}
