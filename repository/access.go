package repository

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

const q = `SELECT exists (SELECT role FROM extaccess WHERE uid IN (SELECT uid FROM usr WHERE username=$1))`

// CheckExtendedAccess - checks if user has extended control of tacplus
func CheckExtendedAccess(db *sql.DB, user string) string {
	chkRole, err := db.Prepare(q)
	if err != nil {
		log.Error(err)
	}
	defer chkRole.Close()

	var rExists bool
	err = chkRole.QueryRow(user).Scan(&rExists)
	if err != nil {
		log.Error(err)
	}

	if rExists {
		getRole, err := db.Prepare("SELECT role FROM extaccess WHERE uid IN (SELECT uid FROM usr WHERE username=$1)")
		if err != nil {
			log.Error(err)
		}
		defer getRole.Close()

		var role string
		err = getRole.QueryRow(user).Scan(&role)
		return role
	}

	return "none"
}
