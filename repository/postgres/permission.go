package postgres

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

func GetPermId(le *logrus.Entry, db *sql.DB, prm string) (int, error) {
	row, err := db.Prepare("SELECT prm_id FROM prm WHERE prm=$1")
	if err != nil {
		return 0, err
	}
	defer row.Close()

	var id int
	err = row.QueryRow(prm).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, err
}
