package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	account string
	token   string
)

func init() {
	flag.StringVar(&account, "account", "", "PagerDuty Account Name")
	flag.StringVar(&token, "token", "", "PagerDuty API Token")
	flag.Parse()
	if account == "" {
		log.Fatal("You must specify an account name")
	}
	if token == "" {
		log.Fatal("You must specify a token")
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("slacker-pagerduty"))
}

func onCallHandler(w http.ResponseWriter, r *http.Request) {
	data, err := getOnCall()
	if err != nil {
		data = fmt.Sprintf("Error retrieving current on call...")
	}
	w.Write([]byte(data))
}

func incidentsHandler(w http.ResponseWriter, r *http.Request) {
	data, err := getIncidents()
	if err != nil {
		data = fmt.Sprintf("Error retrieving current incidents...")
	}
	w.Write([]byte(data))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/oncall", onCallHandler)
	r.HandleFunc("/incidents", incidentsHandler)

	http.Handle("/", r)
	fmt.Println("Running on :8080...")
	http.ListenAndServe(":8080", nil)
}
