package repository

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

var (
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
func GetPasswordHash(le *logrus.Entry, client *aerospike.Client, username string) (string, error) {
	var key *aerospike.Key

	skey := username
	ikey, err := strconv.ParseInt(skey, 10, 64)
	if err == nil {
		key, err = aerospike.NewKey(namespace, set, ikey)
		panicOnError(err)
	} else {
		key, err = aerospike.NewKey(namespace, set, skey)
		panicOnError(err)
	}

	policy := aerospike.NewPolicy()
	policy.SleepBetweenRetries = 50 * time.Millisecond
	policy.MaxRetries = 10
	policy.SleepMultiplier = 2.0

	rec, err := client.Get(policy, key, "pass")
	if err != nil {
		le.WithError(err).Error("aerospike query failed")
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

func extractString(bins aerospike.BinMap, bin string) (string, error) {
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

func CreateUser(le *logrus.Entry, client *aerospike.Client, username, hash, mail, createdBy string, permisID, subdivID int) error {
	var key *aerospike.Key

	skey := username
	ikey, err := strconv.ParseInt(skey, 10, 64)
	if err == nil {
		key, err = aerospike.NewKey(namespace, set, ikey)
		if err != nil {
			return err
		}
	} else {
		key, err = aerospike.NewKey(namespace, set, skey)
		if err != nil {
			return err
		}
	}

	// NOTE: bin name must be less than 16 characters
	bin1 := aerospike.NewBin("username", username)
	bin2 := aerospike.NewBin("pass", hash)
	bin3 := aerospike.NewBin("mail", mail)
	bin4 := aerospike.NewBin("createdBy", createdBy)
	bin5 := aerospike.NewBin("permisID", permisID)
	bin6 := aerospike.NewBin("subdivID", subdivID)
	bin7 := aerospike.NewBin("createdTs", time.Now().Unix())

	policy := aerospike.NewWritePolicy(0, 0)

	err = client.PutBins(policy, key, bin1, bin2, bin3, bin4, bin5, bin6, bin7)
	if err != nil {
		return err
	}

	le.Debugf("record inserted: namespace=%s set=%s key=%v", key.Namespace(), key.SetName(), key.Value())

	return nil
}

func UpdatePassword(aClient *aerospike.Client, username, hash string) error {
	return nil
}

func SetUserStatus(le *logrus.Entry, username string, active bool) error {
	return nil
}

func GetUserCount(aClient *aerospike.Client) *UserSummary {
	return &UserSummary{
		Total:  0,
		Active: 0,
	}
}

type Metrics struct {
	count int
	total int
}

var setMap = make(map[string]Metrics)

func GetUsers(le *logrus.Entry, aclient *aerospike.Client) ([]*User, error) {
	result := make([]*User, 0)

	recs := make([]*aerospike.Record, 0)

	policy := aerospike.NewScanPolicy()
	policy.RecordsPerSecond = 1000

	nodeList := aclient.GetNodes()
	begin := time.Now()

	for _, node := range nodeList {
		le.Debug("scan node ", node.GetName())
		recordset, err := aclient.ScanNode(policy, node, namespace, set)
		if err != nil {
			return result, err
		}

	L:
		for {
			select {
			case rec := <-recordset.Records:
				if rec == nil {
					break L
				}

				metrics, exists := setMap[rec.Key.SetName()]

				if !exists {
					metrics = Metrics{}
				}
				metrics.count++
				metrics.total++
				setMap[rec.Key.SetName()] = metrics

				recs = append(recs, rec)

			case <-recordset.Errors:
				// if there was an error, stop
				panicOnError(err)
			}
		}

		for k, v := range setMap {
			log.Println("Node ", node, " set ", k, " count: ", v.count)
			v.count = 0
		}
	}

	for _, v := range recs {
		// username, err := extractString(v.Bins, "username")
		// if err != nil {
		// 	return result, err
		// }

		mail, err := extractString(v.Bins, "mail")
		if err != nil {
			return result, err
		}

		createdBy, err := extractString(v.Bins, "createdBy")
		if err != nil {
			return result, err
		}

		user := &User{
			// Name:   username,
			Mail:   mail,
			CreaBy: createdBy,
		}

		result = append(result, user)
	}

	end := time.Now()
	seconds := float64(end.Sub(begin)) / float64(time.Second)
	log.Println("Elapsed time: ", seconds, " seconds")

	return result, nil
}
