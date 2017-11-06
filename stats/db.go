package stats

import (
	"database/sql"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)
var db *sql.DB

func InitDB(connStr string) {
	database, err := sql.Open("postgres",connStr)
	if err!= nil{
		log.WithError(err).Fatal()
	}
	db = database
	err = db.Ping()
	if err!= nil{
		log.WithError(err).Fatal("Failed to ping the DB")
	}
	log.Info("Connected to DB")
}



