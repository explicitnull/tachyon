package repository

import (
	"database/sql"
	"fmt"

	as "github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

var (
	host      = "13.48.3.15"
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
	// 	client, err := as.NewClient(host, port)
	// 	if err != nil {
	// 		return "", nil
	// 	}

	// 	var key *as.Key

	// 	skey := username
	// 	ikey, err := strconv.ParseInt(skey, 10, 64)
	// 	if err == nil {
	// 		key, err = as.NewKey(namespace, set, ikey)
	// 		panicOnError(err)
	// 	} else {
	// 		key, err = as.NewKey(namespace, set, skey)
	// 		panicOnError(err)
	// 	}

	// 	policy := as.NewPolicy()

	// 	rec, err := client.Get(policy, key, "pass")
	// 	if err != nil {
	// 		logrus.Errorf("aerospike query failed: %v", err)
	// 		return "", err
	// 	}

	// 	if rec != nil {
	// 		printOK("%v", rec.Bins)
	// 		return extractString(rec.Bins, "pass")

	// 	} else {
	// 		printError("record not found: namespace=%s set=%s key=%v", key.Namespace(), key.SetName(), key.Value())
	// 	}

	// 	return "", nil
	return "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", nil
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func printOK(format string, a ...interface{}) {
	fmt.Printf("ok: "+format+"\n", a...)
}

func printError(format string, a ...interface{}) {
	fmt.Printf("error: "+format+"\n", a...)
}

func extractString(bins as.BinMap, bin string) (string, error) {
	passI, ok := bins[bin]
	if !ok {
		pass, ok := passI.(string)
		if ok {
			return pass, nil
		} else {
			fmt.Println("BinMap value is not string")
		}
	} else {
		fmt.Println("failed to get value from BinMap")
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
