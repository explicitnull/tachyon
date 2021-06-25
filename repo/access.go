package database

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

// CheckExtendedAccess - checks if user has extended control of tacplus
func CheckExtendedAccess(db *sql.DB, user string) string {
	chkRole, err := db.Prepare("select exists (select role from extaccess where uid in (select uid from usr where username=$1))")
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
		getRole, err := db.Prepare("select role from extaccess where uid in (select uid from usr where username=$1)")
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
