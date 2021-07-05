package repository

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

// 	GetDivID returns subdivision ID for DB normalization
func GetDivID(le *logrus.Entry, db *sql.DB, div string) int {
	row, err := db.Prepare("SELECT l2id FROM subdiv WHERE subdiv=$1")
	if err != nil {
		le.WithError(err).Error("")
		return 0
	}
	defer row.Close()

	var id int

	err = row.QueryRow(div).Scan(&id)
	if err != nil {
		le.WithError(err).Error("")
		return 0
	}

	return id
}
