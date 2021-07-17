package postgres

import (
	"database/sql"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

const (
	getUserCount       = `SELECT COUNT(*) FROM usr`
	getActiveUserCount = `SELECT COUNT(*) FROM usr WHERE act=true`
	getUsers           = `SELECT u.uid, u.username, p.prm, u.mail, d.subdiv, u.created, u.act, u.pass_chd, u.created_by FROM usr u JOIN prm p ON (u.prm_id = p.prm_id) JOIN subdiv d ON (u.subdiv_id = d.l2id) ORDER BY p.prm, u.username`
	createUserQuery    = `INSERT INTO usr(username, mail, prm_id, pass, subdiv_id, created_by) values ($1, $2, $3, $4, $5, $6)`
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

func CreateUser(le *logrus.Entry, db *sql.DB, u User, permissionID, subdivisionID int) error {
	stmt, err := db.Prepare(createUserQuery)
	if err != nil {
		le.WithError(err).Error("preparing statement failed")
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.Name, u.Mail, permissionID, u.Hash, subdivisionID, u.CreaBy)
	if err != nil {
		le.WithError(err).Error("executing statement failed")
		return err
	}

	return nil
}

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

const updatePasswordQuery = `UPDATE usr SET pass=$1 WHERE username=$2`

const getPasswordHash = `SELECT pass FROM usr WHERE username=$1`

// GetPasswordHash searches and returns password hash for given username
func GetPasswordHash(db *sql.DB, username string) (string, error) {
	prepared, err := db.Prepare(getPasswordHash)
	if err != nil {
		log.Error(err)
		return "", err
	}
	defer prepared.Close()

	var hash string
	err = prepared.QueryRow(username).Scan(&hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		log.Error(err)
	}

	return hash, err
}

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
