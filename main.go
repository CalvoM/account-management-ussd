package main

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	log.Info("USSD app server starting")
	r := mux.NewRouter()
	r.Path("/ussd/end_note/").HandlerFunc(USSDEndNotificationHandler).Methods("POST")
	r.Path("/ussd/").HandlerFunc(USSDHandler).Methods("POST")
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8083",
	}
	log.Fatal(srv.ListenAndServe())
}
