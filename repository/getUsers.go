package repository

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

const (
	getUserCount       = `SELECT COUNT(*) FROM usr`
	getActiveUserCount = `SELECT COUNT(*) FROM usr WHERE act=true`
	getUsers           = `SELECT u.uid, u.username, p.prm, u.mail, d.subdiv, u.created, u.act, u.pass_chd, u.created_by FROM usr u JOIN prm p ON (u.prm_id = p.prm_id) JOIN subdiv d ON (u.subdiv_id = d.l2id) ORDER BY p.prm, u.username`
)

type User struct {
	Id         int
	Name       string
	Hash       string // Password hash
	Cleartext  string
	Subdiv     string
	Prm        string
	Mail       string
	Active     bool
	ActiveSt   string
	ActiveBox  string // Is "checked" or "" for HTML form
	CreaTime   string // Full time form
	CreaTimeS  string // Short time form
	CreaBy     string
	PassChd    bool
	SubdivList []string
	PrmList    []string
}

type UserSummary struct {
	Total  int
	Active int
}

func GetUserCount(db *sql.DB) *UserSummary {
	sum := new(UserSummary)

	err := db.QueryRow(getUserCount).Scan(&sum.Total)
	if err != nil {
		log.Errorf("quering total users count failed: %v", err)
		return sum
	}

	err = db.QueryRow(getActiveUserCount).Scan(&sum.Active)
	if err != nil {
		log.Errorf("quering active users count failed: %v", err)
		return sum
	}

	return sum
}

func GetUsers(db *sql.DB) []*User {
	res := make([]*User, 0)
	rows, err := db.Query(getUsers)
	if err != nil {
		log.Errorf("quering all users failed: %v", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		u := new(User)

		err = rows.Scan(&u.Id, &u.Name, &u.Prm, &u.Mail, &u.Subdiv, &u.CreaTime, &u.Active, &u.PassChd, &u.CreaBy)
		if err != nil {
			log.WithError(err).Error("scanning sql rows failed")
		}

		res = append(res, u)
	}

	return res
}
