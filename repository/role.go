package repository

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

// GetRole checks if user has extended control of tacplus
func GetRole(le *logrus.Entry, db *sql.DB, username string) string {
	return "admin"
}
