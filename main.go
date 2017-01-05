package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	defaultDogStatsdHost = "127.0.0.1"
	defaultDogStatsdPort = "8125"
	defaultMetricPrefix  = "sendgrid.event."
)

var (
	metricPrefix string
	statsdClient *statsd.Client
)

// SendGridEvents represents the scheme of Event Webhook body
// https://sendgrid.com/docs/API_Reference/Webhooks/event.html#-Event-POST-Example
type SendGridEvents []struct {
	SGMessageID string `json:"sg_message_id"`
	Email       string `json:"email"`
	Timestamp   int    `json:"timestamp"`
	SMTPID      string `json:"smtp-id,omitempty"`
	Event       string `json:"event"`
	Category    string `json:"category,omitempty"`
	URL         string `json:"url,omitempty"`
	AsmGroupID  int    `json:"asm_group_id,omitempty"`
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "https://github.com/dtan4/sendgrid2datadog")
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", err)
		return
	}

	var events SendGridEvents

	if err := json.Unmarshal(body, &events); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", err)
		return
	}

	for _, event := range events {
		if err := statsdClient.Incr(metricPrefix+event.Event, nil, 1); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "%s", err)
			return
		}
	}
}

func main() {
	var dogStatsdHost, dogStatsdPort string

	dogStatsdHost = os.Getenv("DOGSTATSD_HOST")
	if dogStatsdHost == "" {
		dogStatsdHost = defaultDogStatsdHost
	}

	dogStatsdPort = os.Getenv("DOGSTATSD_PORT")
	if dogStatsdPort == "" {
		dogStatsdPort = defaultDogStatsdPort
	}

	dogStatsdAddr := fmt.Sprintf("%s:%s", dogStatsdHost, dogStatsdPort)

	metricPrefix = os.Getenv("METRIC_PREFIX")
	if metricPrefix == "" {
		metricPrefix = defaultMetricPrefix
	}

	var err error

	statsdClient, err = statsd.New(dogStatsdAddr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	r := mux.NewRouter()

	r.Handle("/", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(rootHandler))).Methods("GET")
	r.Handle("/ping", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(pingHandler))).Methods("GET")
	r.Handle("/webhook", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(webhookHandler))).Methods("POST")

	fmt.Println("Server started.")

	http.ListenAndServe(":8080", r)
}
