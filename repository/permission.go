package repository

import (
	"log"
	"tachyon-web/types"
	"time"

	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

const permissionsSet = "permissions"

func GetPermissions(le *logrus.Entry, aclient *aerospike.Client) ([]*types.Permission, error) {
	policy := aerospike.NewScanPolicy()
	policy.RecordsPerSecond = 1000

	nodeList := aclient.GetNodes()
	begin := time.Now()

	recs := make([]*aerospike.Record, 0)

	for _, node := range nodeList {
		le.Debug("scan node ", node.GetName())
		recordset, err := aclient.ScanNode(policy, node, namespace, permissionsSet)
		if err != nil {
			return nil, err
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
			log.Println("Node ", node, " permissionsSet ", k, " count: ", v.count)
			v.count = 0
		}
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

	end := time.Now()
	seconds := float64(end.Sub(begin)) / float64(time.Second)
	log.Println("Elapsed time: ", seconds, " seconds")

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

	// insert into tacacs.permissions (PK, name, description, status, created_by, created_ts) values ('asia-rw', 'asia-rw', 'read only for asia', 'active', 'admin', '2009-01-02 18:00')

	// NOTE: bin name must be less than 16 characters
	bin1 := aerospike.NewBin("name", p.Name)
	bin2 := aerospike.NewBin("description", p.Description)
	bin5 := aerospike.NewBin("status", "active")
	bin3 := aerospike.NewBin("created_by", p.CreatedBy)
	bin4 := aerospike.NewBin("created_ts", time.Now().Unix())

	policy := aerospike.NewWritePolicy(0, 0)

	err = client.PutBins(policy, key, bin1, bin2, bin3, bin4, bin5)
	if err != nil {
		return err
	}

	le.Debugf("record inserted: namespace=%s set=%s key=%v", key.Namespace(), key.SetName(), key.Value())

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
