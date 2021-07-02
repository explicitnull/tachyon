package repository

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

const activateUserQuery = `update usr set act='true' where username=$1`

func ActivateUser(db *sql.DB, le *logrus.Entry, username string) error {
	stmt, err := db.Prepare(activateUserQuery)
	if err != nil {
		le.WithError(err).Errorf("preparing statement failed")
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username)
	if err != nil {
		le.WithError(err).Errorf("setting active flag for user failed")
		return err
	}

	return err
}
