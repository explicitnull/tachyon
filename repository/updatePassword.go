package repository

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

const updatePasswordQuery = `UPDATE usr SET pass=$1 WHERE username=$2`

func UpdatePassword(db *sql.DB, usr, hash string) error {
	stmt, err := db.Prepare(updatePasswordQuery)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(hash, usr)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	return err
}
