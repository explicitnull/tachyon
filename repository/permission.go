package repository

import (
	"tachyon-web/types"
	"time"

	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

const permissionsSet = "permissions"

func GetPermId(le *logrus.Entry, aClient *aerospike.Client, prm string) (int, error) {
	return 2, nil
}

func GetPermissionsList(le *logrus.Entry, aclient *aerospike.Client) []string {
	return []string{"rw", "ro"}
}

func CreatePermission(le *logrus.Entry, client *aerospike.Client, p *types.Permission) error {
	var key *aerospike.Key

	skey := p.Name

	key, err := aerospike.NewKey(namespace, set, skey)
	if err != nil {
		return err
	}

	// insert into tacacs.permissions (PK, name, description, created_by, created_ts, status) values ('test01', 'test01', 'n4bQgYhMfWWaL+qgxVrQFaO/TxsrC4Is0V1sFbDwCgg', 'ma@ti.ru', 2, 10, 'admin', '2009-01-02 18:00', 'active', '2021-07-01 13:00')

	// NOTE: bin name must be less than 16 characters
	bin1 := aerospike.NewBin("name", p.Name)
	bin2 := aerospike.NewBin("description", p.Description)
	bin3 := aerospike.NewBin("created_by", p.CreatedBy)
	bin4 := aerospike.NewBin("created_ts", time.Now().Unix())
	bin5 := aerospike.NewBin("status", "active")

	policy := aerospike.NewWritePolicy(0, 0)

	err = client.PutBins(policy, key, bin1, bin2, bin3, bin4, bin5)
	if err != nil {
		return err
	}

	le.Debugf("record inserted: namespace=%s set=%s key=%v", key.Namespace(), key.SetName(), key.Value())

	return nil
}
