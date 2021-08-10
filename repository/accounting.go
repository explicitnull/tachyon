// Feb  2 18:41:38	192.168.0.202	jathan	ttyp1	192.168.0.10	stop	task_id=3	service=shell	process*mgd[7741]	cmd=configure exclusive <cr>
// insert into tacacs.accounting(PK, id, ts, device_ip, device_name, account, user_ip, user_fqdn, command) values ('wTHCHKgRtpav', 'wTHCHKgRtpav', '2021-01-01 01:23', '10.10.10.10', 'device04', 'jathan', '45.67.89.90', 'reverse.example.com', 'configure exclusive <cr>')
package repository

import (
	"tacacs-webconsole/types"

	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

const accountingSet = "accounting"

func GetAccounting(le *logrus.Entry, aclient *aerospike.Client) ([]*types.AccountingRec, error) {
	recs, err := getAllRecords(aclient, accountingSet)
	if err != nil {
		return nil, err
	}

	result := make([]*types.AccountingRec, 0)

	for _, v := range recs {
		acct, err := extractAccountingRec(v.Bins)
		if err != nil {
			le.WithError(err).Error("extracting bins failed")
			return nil, err
		}

		result = append(result, acct)
	}

	return result, nil
}

func extractAccountingRec(bins aerospike.BinMap) (*types.AccountingRec, error) {
	acct := &types.AccountingRec{}

	var err error

	// fields ordered according to tacplus log files
	acct.ID, err = extractString(bins, "id")
	if err != nil {
		return nil, err
	}

	acct.Timestamp, err = extractString(bins, "ts")
	if err != nil {
		return nil, err
	}

	acct.DeviceIP, err = extractString(bins, "device_ip")
	if err != nil {
		return nil, err
	}

	acct.DeviceName, err = extractString(bins, "device_name")
	if err != nil {
		return nil, err
	}

	acct.AccountName, err = extractString(bins, "account")
	if err != nil {
		return nil, err
	}

	acct.UserIP, err = extractString(bins, "user_ip")
	if err != nil {
		return nil, err
	}

	acct.UserFQDN, err = extractString(bins, "user_fqdn")
	if err != nil {
		return nil, err
	}

	acct.Command, err = extractString(bins, "command")
	if err != nil {
		return nil, err
	}

	return acct, nil
}
