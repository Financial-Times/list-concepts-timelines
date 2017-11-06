package main

import (
	"net/http"
	log "github.com/sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/Financial-Times/list-concepts-timelines/handlers"
	"github.com/Financial-Times/list-concepts-timelines/stats"
	"github.com/jawher/mow.cli"
)

func main() {

	app := cli.App("list-concepts-timelines","list-concepts-timelines")

	dbUrl := app.String(cli.StringOpt{
		Name:   "db-url",
		Desc:   "redshift-url",
		EnvVar: "DB_URL",
	})

	stats.InitDB(*dbUrl)

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.PathPrefix("/topConcepts").Handler(handlers.NewTopConceptsHandler())
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.WithError(err).Fatal()
	}
}