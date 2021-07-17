package repository

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	as "github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

var (
	host      = "127.0.0.1"
	port      = 3000
	namespace = "tacacs"
	set       = "users"
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

// GetPasswordHash searches for password hash of given user
func GetPasswordHash(db *sql.DB, username string) (string, error) {
	client, err := as.NewClient(host, port)
	if err != nil {
		return "", nil
	}

	var key *as.Key

	skey := username
	ikey, err := strconv.ParseInt(skey, 10, 64)
	if err == nil {
		key, err = as.NewKey(namespace, set, ikey)
		panicOnError(err)
	} else {
		key, err = as.NewKey(namespace, set, skey)
		panicOnError(err)
	}

	policy := as.NewPolicy()
	rec, err := client.Get(policy, key, "pass")
	panicOnError(err)

	if rec != nil {
		printOK("%v", rec.Bins)
		return extractString(rec.Bins, "pass")

	} else {
		printError("record not found: namespace=%s set=%s key=%v", key.Namespace(), key.SetName(), key.Value())
	}

	return "", nil
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func printOK(format string, a ...interface{}) {
	fmt.Printf("ok: "+format+"\n", a...)
	os.Exit(0)
}

func printError(format string, a ...interface{}) {
	fmt.Printf("error: "+format+"\n", a...)
	os.Exit(1)
}

func extractString(bins as.BinMap, bin string) (string, error) {
	passI, ok := bins[bin]
	if ok {
		pass, ok := passI.(string)
		if ok {
			return pass, nil
		}
	}

	return "", nil
}

func CreateUser(le *logrus.Entry, username, hash, mail, createdBy string, permisID, subdivID int) error {
	return nil
}

func UpdatePassword(db *sql.DB, username, hash string) error {
	return nil
}

func SetUserStatus(le *logrus.Entry, username string, active bool) error {
	return nil
}

func GetUserCount(db *sql.DB) *UserSummary {
	return &UserSummary{}
}

func GetUsers(db *sql.DB) []*User {
	res := make([]*User, 0)

	return res
}
