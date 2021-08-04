// insert into tacacs.users (PK, id, username, pass, subdiv_id, permis_id, mail, status, created_ts, created_by, pass_set_ts) values ('furai', 'QgYhMfWWaLqg', 'furai', 'n4bQgYhMfWWaL+qgxVrQFaO/TxsrC4Is0V1sFbDwCgg', 2, 10, 'dmitry.zhargalov@rt.ru', 'pass_not_chd', '2021-01-01 18:00', 'admin', '')
// insert into tacacs.users (PK, id, username, pass, subdiv_id, permis_id, mail, status, created_ts, created_by, pass_chd_ts, ui_role) values ('test09', 'srC4Is0V1sFb', 'test09', 'n4bQgYhMfWWaL+qgxVrQFaO/TxsrC4Is0V1sFbDwCgg', 2, 10, 'anton.moroz@rt.ru', 'suspended', '2009-01-02 18:00', 'furai', '2021-07-01 13:00', "manager")
package repository

import (
	"errors"
	"log"
	"strconv"
	"tacacs-webconsole/types"
	"time"

	"github.com/aerospike/aerospike-client-go"
	"github.com/dchest/uniuri"
	"github.com/sirupsen/logrus"
)

const (
	accountsSet = "users"
)

func GetAccounts(le *logrus.Entry, aclient *aerospike.Client) ([]*types.Account, error) {
	begin := time.Now()

	recs, err := getAllRecords(aclient, accountsSet)
	if err != nil {
		return nil, err
	}

	result := make([]*types.Account, 0)

	for _, v := range recs {
		acc, err := extractAccount(v.Bins)
		if err != nil {
			le.WithError(err).Error("extracting bins failed")
			return nil, err
		}

		result = append(result, acc)
	}

	end := time.Now()
	seconds := float64(end.Sub(begin)) / float64(time.Second)
	log.Println("Elapsed time: ", seconds, " seconds")

	return result, nil
}

// GetPasswordHash searches for password hash of given user
func GetPasswordHash(le *logrus.Entry, client *aerospike.Client, username string) (string, error) {
	var key *aerospike.Key

	skey := username
	ikey, err := strconv.ParseInt(skey, 10, 64)
	if err == nil {
		key, err = aerospike.NewKey(namespace, accountsSet, ikey)
		panicOnError(err)
	} else {
		key, err = aerospike.NewKey(namespace, accountsSet, skey)
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
		printError("record not found: namespace=%s accountsSet=%s key=%v", key.Namespace(), key.SetName(), key.Value())
	}

	return "", nil
}

func CreateUser(le *logrus.Entry, client *aerospike.Client, username, hash, mail, createdBy string, permisID, subdivID int) error {
	var key *aerospike.Key

	skey := username

	key, err := aerospike.NewKey(namespace, accountsSet, skey)
	if err != nil {
		return err
	}

	id := uniuri.NewLen(12)
	// TODO: set status depending on checkbox "need password change"
	// NOTE: bin name must be less than 16 characters
	bins := aerospike.BinMap{
		"id":        id,
		"username":  username,
		"pass":      hash,
		"subdivID":  subdivID,
		"permisID":  permisID,
		"mail":      mail,
		"status":    "active",
		"createdTs": time.Now().Unix(),
		"createdBy": createdBy,
	}

	policy := aerospike.NewWritePolicy(0, 0)

	err = client.Put(policy, key, bins)
	if err != nil {
		return err
	}

	le.Debugf("record inserted: namespace=%s accountsSet=%s key=%v", key.Namespace(), key.SetName(), key.Value())

	return nil
}

func GetAccountByName(le *logrus.Entry, client *aerospike.Client, name string) (*types.Account, error) {
	var key *aerospike.Key

	skey := name
	key, err := aerospike.NewKey(namespace, accountsSet, skey)
	panicOnError(err)

	policy := aerospike.NewPolicy()
	policy.SleepBetweenRetries = 50 * time.Millisecond
	policy.MaxRetries = 10
	policy.SleepMultiplier = 2.0

	rec, err := client.Get(policy, key)
	if err != nil {
		le.WithError(err).Error("aerospike query failed")
		return nil, err
	}

	if rec == nil {
		printError("record not found: namespace=%s accountsSet=%s key=%v", key.Namespace(), key.SetName(), key.Value())
		return nil, errors.New("record not found")
	}

	acc, err := extractAccount(rec.Bins)
	if err != nil {
		le.WithError(err).Error("extracting bins failed")
		return nil, err
	}

	return acc, nil
}

func SetPassword(aClient *aerospike.Client, username, hash string) error {
	return nil
}

func SetAccountStatus(le *logrus.Entry, name string, status string) error {
	return nil
}

func SetSubdivision(name string, subdiv string) error {
	return nil
}

func SetPermission(name string, perm string) error {
	return nil
}

func SetMail(name string, mail string) error {
	return nil
}

func extractAccount(bins aerospike.BinMap) (*types.Account, error) {
	username, err := extractString(bins, "username")
	if err != nil {
		return nil, err
	}

	mail, err := extractString(bins, "mail")
	if err != nil {
		return nil, err
	}

	subdivID, err := extractInt(bins, "subdiv_id")
	if err != nil {
		return nil, err
	}

	permisID, err := extractInt(bins, "permis_id")
	if err != nil {
		return nil, err
	}

	createdBy, err := extractString(bins, "created_by")
	if err != nil {
		return nil, err
	}

	createdTimestamp, err := extractString(bins, "created_ts")
	if err != nil {
		return nil, err
	}

	status, err := extractString(bins, "status")
	if err != nil {
		return nil, err
	}

	passwordChangedTimestamp, err := extractString(bins, "pass_chd_ts")
	if err != nil {
		return nil, err
	}

	uiRole, err := extractString(bins, "ui_role")
	if err != nil {
		return nil, err
	}

	acc := &types.Account{
		Name:                     username,
		Mail:                     mail,
		Subdivision:              strconv.Itoa(subdivID),
		Permission:               strconv.Itoa(permisID),
		CreatedBy:                createdBy,
		CreatedTimestamp:         createdTimestamp,
		Status:                   status,
		PasswordChangedTimestamp: passwordChangedTimestamp,
		UIRole:                   uiRole,
	}

	return acc, nil
}
