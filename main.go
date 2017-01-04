package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s", r.UserAgent())
}

func main() {
	r := http.NewServeMux()

	r.Handle("/", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(rootHandler)))

	http.ListenAndServe(":8080", r)
}
