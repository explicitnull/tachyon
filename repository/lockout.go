// insert into tacacs.lockouts (PK, ip, fqdn, attempts, first_ts, last_ts, last_device, last_account) values ('38.253.3.135', '3.example.com', 41, '2021-07-12 19:00', '2021-07-23 10:00', 'stkm-core01', 'admin')
package repository

import (
	"tacacs-webconsole/types"

	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

const lockoutsSet = "lockouts"

func GetLockouts(le *logrus.Entry, aclient *aerospike.Client) ([]*types.Lockout, error) {
	recs, err := getAllRecords(aclient, lockoutsSet)
	if err != nil {
		return nil, err
	}

	result := make([]*types.Lockout, 0)

	for _, v := range recs {
		lock, err := extractLockout(v.Bins)
		if err != nil {
			le.WithError(err).Error("extracting bins failed")
			return nil, err
		}

		result = append(result, lock)
	}

	return result, nil
}

func extractLockout(bins aerospike.BinMap) (*types.Lockout, error) {
	lock := &types.Lockout{}

	var err error

	lock.IP, err = extractString(bins, "ip")
	if err != nil {
		return nil, err
	}

	lock.FQDN, err = extractString(bins, "fqdn")
	if err != nil {
		return nil, err
	}

	lock.Attempts, err = extractInt(bins, "attempts")
	if err != nil {
		return nil, err
	}

	lock.FirstAttemptTimestamp, err = extractString(bins, "first_ts")
	if err != nil {
		return nil, err
	}

	lock.LastAttemptTimestamp, err = extractString(bins, "last_ts")
	if err != nil {
		return nil, err
	}

	lock.LastDevice, err = extractString(bins, "last_device")
	if err != nil {
		return nil, err
	}

	lock.LastAccountName, err = extractString(bins, "last_account")
	if err != nil {
		return nil, err
	}

	return lock, nil
}
