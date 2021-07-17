package repository

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

// 	GetSubdivisionID returns subdivision ID for DB normalization
func GetSubdivisionID(le *logrus.Entry, db *sql.DB, div string) int {
	return 0
}
