package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	log "github.com/sirupsen/logrus"
)

func Open(host, dbname, user, password string) *sql.DB {
	params := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable", host, dbname, user, password)

	db, err := sql.Open("postgres", params)
	if err != nil {
		log.WithError(err).Errorf("opening db connection failed")
	}
	defer db.Close()

	return db
}
