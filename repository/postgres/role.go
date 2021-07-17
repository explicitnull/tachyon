package postgres

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

const (
	existsQuery  = `SELECT EXISTS (SELECT role FROM extaccess WHERE uid IN (SELECT uid FROM usr WHERE username=$1))`
	getRoleQuery = `SELECT role FROM extaccess WHERE uid IN (SELECT uid FROM usr WHERE username=$1)`
)

// CheckRole checks if user has extended control of tacplus
func CheckRole(le *logrus.Entry, db *sql.DB, username string) string {
	preFlightPrepared, err := db.Prepare(existsQuery)
	if err != nil {
		le.Error(err)
		return "none"
	}
	defer preFlightPrepared.Close()

	var roleExists bool

	err = preFlightPrepared.QueryRow(username).Scan(&roleExists)
	if err != nil {
		le.Error(err)
		return "none"
	}

	if roleExists {
		prepared, err := db.Prepare(getRoleQuery)
		if err != nil {
			le.Error(err)
		}
		defer prepared.Close()

		var role string
		err = prepared.QueryRow(username).Scan(&role)
		if err != nil {
			le.Error(err)
			return "none"
		}

		return role
	}

	return "none"
}
