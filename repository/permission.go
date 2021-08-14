package repository

import (
	"errors"
	"tacacs-webconsole/types"
	"time"

	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

const permissionsSet = "permissions"

func GetPermissions(le *logrus.Entry, aclient *aerospike.Client) ([]*types.Permission, error) {
	recs, err := getAllRecords(aclient, permissionsSet)
	if err != nil {
		return nil, err
	}

	result := make([]*types.Permission, 0)

	for _, v := range recs {
		perm, err := extractPermission(v.Bins)
		if err != nil {
			le.WithError(err).Error("extracting bins failed")
			return nil, err
		}

		result = append(result, perm)
	}

	return result, nil
}

func GetPermId(le *logrus.Entry, aClient *aerospike.Client, prm string) (int, error) {
	return 2, nil
}

func GetPermissionsList(le *logrus.Entry, aclient *aerospike.Client) []string {
	return []string{"rw", "ro"}
}

func CreatePermission(le *logrus.Entry, client *aerospike.Client, p *types.Permission) error {
	var key *aerospike.Key

	skey := p.Name

	key, err := aerospike.NewKey(namespace, permissionsSet, skey)
	if err != nil {
		return err
	}

	// insert into tacacs.permissions (PK, name, description, status, created_by, created_ts) values ('asia-rw', 'asia-rw', 'read only for asia', 'active', 'superuser', '2009-01-02 18:00')

	// NOTE: bin name must be less than 16 characters
	bins := aerospike.BinMap{
		"name":        p.Name,
		"description": p.Description,
		"status":      "active",
		"created_by":  p.CreatedBy,
		"created_ts":  time.Now().Unix(),
	}

	policy := aerospike.NewWritePolicy(0, 0)

	err = client.Put(policy, key, bins)
	if err != nil {
		return err
	}

	le.Debugf("record inserted: namespace=%s set=%s key=%v", key.Namespace(), key.SetName(), key.Value())

	return nil
}

func GetPermissionByName(le *logrus.Entry, client *aerospike.Client, name string) (*types.Permission, error) {
	var key *aerospike.Key

	key, err := aerospike.NewKey(namespace, permissionsSet, name)
	if err != nil {
		return nil, err
	}

	policy := aerospike.NewPolicy()

	rec, err := client.Get(policy, key)
	if err != nil {
		return nil, err
	}

	if rec == nil {
		return nil, errors.New("record not found")
	}

	perm, err := extractPermission(rec.Bins)
	if err != nil {
		return nil, err
	}

	return perm, nil
}

func SetPermissionStatus(le *logrus.Entry, name string, status string) error {
	return nil
}

func SetPermissionDescription(le *logrus.Entry, name string, description string) error {
	return nil
}

func extractPermission(bins aerospike.BinMap) (*types.Permission, error) {
	perm := &types.Permission{}

	var err error

	perm.Name, err = extractString(bins, "name")
	if err != nil {
		return nil, err
	}

	perm.Description, err = extractString(bins, "description")
	if err != nil {
		return nil, err
	}

	perm.Status, err = extractString(bins, "status")
	if err != nil {
		return nil, err
	}

	perm.CreatedBy, err = extractString(bins, "created_by")
	if err != nil {
		return nil, err
	}

	perm.CreatedTimestamp, err = extractString(bins, "created_ts")
	if err != nil {
		return nil, err
	}

	return perm, nil
}
