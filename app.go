package main

import (
	"net/http"
	log "github.com/sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/Financial-Times/list-concepts-timelines/handlers"
)

func main() {

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.PathPrefix("/topConcepts").Handler(handlers.NewTopConceptsHandler())
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.WithError(err).Fatal()
	}
}