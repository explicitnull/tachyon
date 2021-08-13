// insert into tacacs.authentications(PK, id, ts, account, device_ip, device_name, event_type, user_ip, user_fqdn) values ('76p7wTHCHK', '76p7wTHCHK', 1628772995, 'root', '10.10.10.10', 'device04', 'user_not_found', '45.67.89.90', 'reverse.example.com')
package repository

import (
	"tacacs-webconsole/types"
	"time"

	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

const authenticationsSet = "authentications"

func GetAuthentications(le *logrus.Entry, aclient *aerospike.Client) ([]types.Authentication, error) {
	recs, err := getAllRecords(aclient, authenticationsSet)
	if err != nil {
		return nil, err
	}

	result := make([]types.Authentication, 0)

	for _, v := range recs {
		auth, err := extractAuthentication(v.Bins)
		if err != nil {
			le.WithError(err).Error("extracting bins failed")
			return nil, err
		}

		result = append(result, auth)
	}

	le.Debugf("authentication found: %d", len(result))
	return result, nil
}

func GetAuthenticationWithTimeFilter(le *logrus.Entry, aclient *aerospike.Client, begin, end time.Time) ([]types.Authentication, error) {
	le.Debugf("repo time begin: %s, end: %s", begin, end)

	res := make([]types.Authentication, 0)

	beginTS := begin.Unix()
	endTS := end.Unix()

	le.Debugf("repo stamp begin: %d, end: %d", beginTS, endTS)

	records, err := getRecordsWithRangeFilter(aclient, authenticationsSet, "ts", beginTS, endTS)
	if err != nil {
		return res, err
	}

	for _, v := range records {
		aut, err := extractAuthentication(v.Bins)
		if err != nil {
			return nil, err
		}

		res = append(res, aut)
	}

	return res, nil
}

func extractAuthentication(bins aerospike.BinMap) (types.Authentication, error) {
	auth := types.Authentication{}

	var err error

	auth.ID, err = extractString(bins, "id")
	if err != nil {
		return types.Authentication{}, err
	}

	ts, err := extractInt(bins, "ts")
	if err != nil {
		return types.Authentication{}, err
	}
	tm := time.Unix(int64(ts), 0)
	auth.Timestamp = tm.Format(types.TimeFormatSeconds)

	auth.AccountName, err = extractString(bins, "account")
	if err != nil {
		return types.Authentication{}, err
	}

	auth.DeviceIP, err = extractString(bins, "device_ip")
	if err != nil {
		return types.Authentication{}, err
	}

	auth.DeviceName, err = extractString(bins, "device_name")
	if err != nil {
		return types.Authentication{}, err
	}

	auth.EventType, err = extractString(bins, "event_type")
	if err != nil {
		return types.Authentication{}, err
	}

	auth.UserIP, err = extractString(bins, "user_ip")
	if err != nil {
		return types.Authentication{}, err
	}

	auth.UserFQDN, err = extractString(bins, "user_fqdn")
	if err != nil {
		return types.Authentication{}, err
	}

	return auth, nil
}
