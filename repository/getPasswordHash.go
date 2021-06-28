package repository

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

const getPasswordHash = `SELECT pass FROM usr WHERE username=$1`

// GetPasswordHash searches and returns password hash by given username
func GetPasswordHash(db *sql.DB, user string) (string, error) {
	prepared, err := db.Prepare(getPasswordHash)
	if err != nil {
		log.Error(err)
	}
	defer prepared.Close()

	var hash string
	err = prepared.QueryRow(user).Scan(&hash)
	if err != nil {
		if err == sql.ErrNoRows {
			log.WithField("username", user).Warnf("username not found in database", user)
			return "", err
		}
		log.Error(err)
	}

	return hash, err
}
