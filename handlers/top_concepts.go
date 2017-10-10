package handlers

import (
	"html/template"
	"net/http"

	"github.com/Financial-Times/list-concepts-timelines/stats"
	log "github.com/sirupsen/logrus"
	"time"
)

type TopConceptsHandler struct {
}

func NewTopConceptsHandler() *TopConceptsHandler {
	return &TopConceptsHandler{}
}

func (h *TopConceptsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	startString := ""
	if len(values["start"]) > 0 {
		startString = values["start"][0]
	}

	start, err := time.Parse(stats.DateLayout, startString)
	if err != nil {
		start = time.Now().Add(-24 * time.Hour)
	}

	stopString := ""
	if len(values["stop"]) > 0 {
		stopString = values["stop"][0]
	}

	stop, err := time.Parse(stats.DateLayout, stopString)
	if err != nil {
		stop = time.Now()
	}

	topConcepts := stats.NewTopConcepts(start, stop)

	t, err := template.ParseFiles("templates/topConcepts.html")
	if err != nil {
		log.WithError(err).Error("Error in parsing concept template")
	}
	t.Execute(w, topConcepts)
}
