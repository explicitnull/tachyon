package repository

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

const getPasswordHash = `SELECT pass FROM usr WHERE username=$1`

// GetPasswordHash searches and returns password hash for given username
func GetPasswordHash(db *sql.DB, username string) (string, error) {
	prepared, err := db.Prepare(getPasswordHash)
	if err != nil {
		log.Error(err)
		return "", err
	}
	defer prepared.Close()

	var hash string
	err = prepared.QueryRow(username).Scan(&hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		log.Error(err)
	}

	return hash, err
}
