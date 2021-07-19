package repository

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

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
	policy.SleepBetweenRetries = 50 * time.Millisecond
	policy.MaxRetries = 10
	policy.SleepMultiplier = 2.0

	rec, err := client.Get(policy, key, "pass")
	if err != nil {
		logrus.Errorf("aerospike query failed: %v", err)
		return "", err
	}

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
}

func printError(format string, a ...interface{}) {
	fmt.Printf("error: "+format+"\n", a...)
}

func extractString(bins as.BinMap, bin string) (string, error) {
	passI, ok := bins[bin]
	if ok {
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
	client, err := as.NewClient(host, port)
	if err != nil {
		return err
	}

	var key *as.Key

	skey := username
	ikey, err := strconv.ParseInt(skey, 10, 64)
	if err == nil {
		key, err = as.NewKey(namespace, set, ikey)
		if err != nil {
			return err
		}
	} else {
		key, err = as.NewKey(namespace, set, skey)
		if err != nil {
			return err
		}
	}

	// NOTE: bin name must be less than 16 characters
	bin1 := as.NewBin("username", username)
	bin2 := as.NewBin("hash", hash)
	bin3 := as.NewBin("mail", mail)
	bin4 := as.NewBin("createdBy", createdBy)
	bin5 := as.NewBin("permisID", permisID)
	bin6 := as.NewBin("subdivID", subdivID)
	bin7 := as.NewBin("createdTS", time.Now().Unix())

	policy := as.NewWritePolicy(0, 0)

	err = client.PutBins(policy, key, bin1, bin2, bin3, bin4, bin5, bin6, bin7)
	if err != nil {
		return err
	}

	le.Debugf("record inserted: namespace=%s set=%s key=%v", key.Namespace(), key.SetName(), key.Value())

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
