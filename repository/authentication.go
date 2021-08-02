// insert into tacacs.authentications(PK, id, account, device_ip, device_name, event_type, user_ip, user_fqdn) values ('yVZfYfB76p7wTHCHKgRb', 'yVZfYfB76p7wTHCHKgRb', 't1000', '10.10.10.10', 'device04', 'user_not_found', '45.67.89.90', 'reverse.example.com')
package repository

import (
	"tacasa-web/types"

	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

const authenticationsSet = "authentications"

func GetAuthentications(le *logrus.Entry, aclient *aerospike.Client) ([]*types.Authentication, error) {
	recs, err := getAllRecords(aclient, authenticationsSet)
	if err != nil {
		return nil, err
	}

	result := make([]*types.Authentication, 0)

	for _, v := range recs {
		auth, err := extractAuthentication(v.Bins)
		if err != nil {
			le.WithError(err).Error("extracting bins failed")
			return nil, err
		}

		result = append(result, auth)
	}

	return result, nil
}

func extractAuthentication(bins aerospike.BinMap) (*types.Authentication, error) {
	auth := &types.Authentication{}

	var err error

	auth.ID, err = extractString(bins, "id")
	if err != nil {
		return nil, err
	}

	auth.Timestamp, err = extractString(bins, "ts")
	if err != nil {
		return nil, err
	}

	auth.AccountName, err = extractString(bins, "account")
	if err != nil {
		return nil, err
	}

	auth.DeviceIP, err = extractString(bins, "device_ip")
	if err != nil {
		return nil, err
	}

	auth.DeviceName, err = extractString(bins, "device_name")
	if err != nil {
		return nil, err
	}

	auth.EventType, err = extractString(bins, "event_type")
	if err != nil {
		return nil, err
	}

	auth.UserIP, err = extractString(bins, "user_ip")
	if err != nil {
		return nil, err
	}

	auth.UserFQDN, err = extractString(bins, "user_fqdn")
	if err != nil {
		return nil, err
	}

	return auth, nil
}
