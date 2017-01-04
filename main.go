package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s", r.UserAgent())
}

func main() {
	r := http.NewServeMux()

	r.Handle("/", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(rootHandler)))
	r.Handle("/ping", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(pingHandler)))

	http.ListenAndServe(":8080", r)
}
