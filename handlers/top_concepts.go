package handlers

import (
	"net/http"

	"encoding/json"
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
	log.WithField("url", r.URL).Info("http request")
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

	listUUID := ""
	if len(values["listUUID"]) > 0 {
		listUUID = values["listUUID"][0]
	}
	if listUUID == "" {
		listUUID = "520ddb76-e43d-11e4-9e89-00144feab7de"
	}

	topConcepts := stats.GetTopConcepts(start, stop, listUUID)

	json.NewEncoder(w).Encode(topConcepts.GetTimeLine())
}
