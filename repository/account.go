// insert into tacacs.users (PK, id, username, pass, subdiv_id, permis_id, mail, status, created_ts, created_by, pass_set_ts) values ('furai', 'QgYhMfWWaLqg', 'furai', 'n4bQgYhMfWWaL+qgxVrQFaO/TxsrC4Is0V1sFbDwCgg', 2, 10, 'dmitry.zhargalov@rt.ru', 'pass_not_chd', '2021-01-01 18:00', 'superuser', '')
// insert into tacacs.users (PK, id, username, pass, subdiv_id, permis_id, mail, status, created_ts, created_by, pass_chd_ts, ui_level) values ('test09', 'srC4Is0V1sFb', 'test09', 'n4bQgYhMfWWaL+qgxVrQFaO/TxsrC4Is0V1sFbDwCgg', 2, 10, 'anton.moroz@rt.ru', 'suspended', '2009-01-02 18:00', 'furai', '2021-07-01 13:00', "manager")
package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"tachyon/types"
	"time"

	"github.com/aerospike/aerospike-client-go"
	"github.com/dchest/uniuri"
)

const (
	accountsSet = "users"
)

func GetAccounts(ctx context.Context, aclient *aerospike.Client) ([]*types.Account, error) {
	// begin := time.Now()

	recs, err := getAllRecords(aclient, accountsSet)
	if err != nil {
		return nil, err
	}

	result := make([]*types.Account, 0)

	for _, v := range recs {
		acc, err := extractAccount(v.Bins)
		if err != nil {
			return nil, fmt.Errorf("extracting bins failed: %v", err)
		}

		result = append(result, acc)
	}

	// end := time.Now()
	// seconds := float64(end.Sub(begin)) / float64(time.Second)
	// le.Debug("Elapsed time: ", seconds, " seconds")

	return result, nil
}

func GetAccountByName(ctx context.Context, client *aerospike.Client, name string) (*types.Account, error) {
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
		return nil, fmt.Errorf("aerospike query failed: %v", err)
	}

	if rec == nil {
		printError("record not found: namespace=%s accountsSet=%s key=%v", key.Namespace(), key.SetName(), key.Value())
		return nil, errors.New("record not found")
	}

	acc, err := extractAccount(rec.Bins)
	if err != nil {
		return nil, fmt.Errorf("extracting bins failed: %v", err)
	}

	return acc, nil
}

// GetPasswordHash searches for password hash of given user
func GetPasswordHash(ctx context.Context, client *aerospike.Client, username string) (string, error) {
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
		return "", fmt.Errorf("aerospike query failed: %v", err)
	}

	if rec != nil {
		printOK("%v", rec.Bins)
		return extractString(rec.Bins, "pass")

	} else {
		printError("record not found: namespace=%s accountsSet=%s key=%v", key.Namespace(), key.SetName(), key.Value())
	}

	return "", nil
}

// GetAccessLevel checks if user can see or change tacplus configuration
func GetAccessLevel(ctx context.Context, aclient *aerospike.Client, username string) (string, error) {
	key, err := aerospike.NewKey(namespace, accountsSet, username)
	if err != nil {
		return "", err
	}

	policy := aerospike.NewPolicy()

	rec, err := aclient.Get(policy, key, "ui_role")
	if err != nil {
		return "", fmt.Errorf("aerospike query failed: %v", err)
	}

	if rec == nil {
		return "", errors.New("aerospike record not found")
	}

	level, err := extractString(rec.Bins, "ui_role")
	if err != nil {
		return "", err
	}

	// le.Debugf("access level - %s", level)

	return level, nil
}

func CreateUser(ctx context.Context, client *aerospike.Client, username, hash, mail, createdBy string, permisID, subdivID int) error {
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

	// le.Debugf("record inserted: namespace=%s accountsSet=%s key=%v", key.Namespace(), key.SetName(), key.Value())

	return nil
}

func SetPassword(ctx context.Context, aclient *aerospike.Client, acname, hash string) error {
	err := setBinString(aclient, accountsSet, acname, "pass", hash)
	if err != nil {
		return err
	}
	return nil
}

func SetAccountStatus(ctx context.Context, aclient *aerospike.Client, acname string, status string) error {
	err := setBinString(aclient, accountsSet, acname, "status", status)
	if err != nil {
		return err
	}

	return nil
}

func SetSubdivision(ctx context.Context, aclient *aerospike.Client, acname string, subdiv int) error {
	err := setBinInt(aclient, accountsSet, acname, "subdiv_id", subdiv)
	if err != nil {
		return err
	}

	return nil
}

func SetPermission(ctx context.Context, aclient *aerospike.Client, acname string, perm int) error {
	err := setBinInt(aclient, accountsSet, acname, "permis_id", perm)
	if err != nil {
		return err
	}

	return nil
}

func SetMail(ctx context.Context, aclient *aerospike.Client, acname string, mail string) error {
	err := setBinString(aclient, accountsSet, acname, "mail", mail)
	if err != nil {
		return err
	}

	return nil
}

func SetUILevel(ctx context.Context, aclient *aerospike.Client, acname string, role string) error {
	err := setBinString(aclient, accountsSet, acname, "ui_role", role)
	if err != nil {
		return err
	}

	return nil
}

func DeleteAccount(ctx context.Context, aclient *aerospike.Client, acname, authenticatedUsername string) error {
	_, err := deleteRecord(aclient, accountsSet, acname)
	if err != nil {
		return err
	}

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

	uiLevel, err := extractString(bins, "ui_role")
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
		UILevel:                  uiLevel,
	}

	return acc, nil
}
