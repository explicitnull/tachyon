// Feb  2 18:41:38	192.168.0.202	jathan	ttyp1	192.168.0.10	stop	task_id=3	service=shell	process*mgd[7741]	cmd=configure exclusive <cr>
// insert into tacacs.accounting(PK, id, ts, device_ip, device_name, account, user_ip, user_fqdn, command) values ('1jcnrdeslA', '1jcnrdeslA', 1610000000, '12.12.12.12', 'device34', 'superuser', '45.67.89.90', 'reverse.example.com', 'login <cr>')
// insert into tacacs.accounting(PK, id, ts, device_ip, device_name, account, user_ip, user_fqdn, command) values ('aFCHxkuria', 'aFCHxkuria', 1600000000, '10.10.10.10', 'device01', 'jathan', '45.67.89.90', 'reverse.example.com', 'configure exclusive <cr>')
// insert into tacacs.accounting(PK, id, ts, device_ip, device_name, account, user_ip, user_fqdn, command) values ('abCHKRtpav', 'abCHKRtpav', 1610000000, '11.11.11.11', 'device22', 'superuser', '45.67.89.90', 'reverse.example.com', 'set system login class view-only permissions [ view view-configuration ] <cr>')
// insert into tacacs.accounting(PK, id, ts, device_ip, device_name, account, user_ip, user_fqdn, command) values ('abCHKRtpaz', 'abCHKRtpaz', 1628772960, '13.13.13.13', 'rtr23', 'furai', '3.4.5.6', 'reverse.example.com', 'login')

package repository

import (
	"tacacs-webconsole/types"
	"time"

	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

const accountingSet = "accounting"

func GetAccounting(le *logrus.Entry, aclient *aerospike.Client) ([]types.AccountingRecord, error) {
	recs, err := getAllRecords(aclient, accountingSet)
	if err != nil {
		return nil, err
	}

	result := make([]types.AccountingRecord, 0)

	for _, v := range recs {
		acct, err := extractAccountingRecord(v.Bins)
		if err != nil {
			le.WithError(err).Error("extracting bins failed")
			return nil, err
		}

		result = append(result, acct)
	}

	return result, nil
}

func GetAccountingWithEqualFilter(le *logrus.Entry, aclient *aerospike.Client, field, value string) ([]types.AccountingRecord, error) {
	res := make([]types.AccountingRecord, 0)

	records, err := getRecordsWithEqualFilter(aclient, accountingSet, field, value)
	if err != nil {
		return res, err
	}

	for _, v := range records {
		acct, err := extractAccountingRecord(v.Bins)
		if err != nil {
			return nil, err
		}

		res = append(res, acct)
	}

	return res, nil
}

func GetAccountingWithTimeFilter(le *logrus.Entry, aclient *aerospike.Client, begin, end time.Time) ([]types.AccountingRecord, error) {
	le.Debugf("repo time begin: %s, end: %s", begin, end)

	res := make([]types.AccountingRecord, 0)

	beginTS := begin.Unix()
	endTS := end.Unix()

	le.Debugf("repo stamp begin: %d, end: %d", beginTS, endTS)

	records, err := getRecordsWithRangeFilter(aclient, accountingSet, "ts", beginTS, endTS)
	if err != nil {
		return res, err
	}

	for _, v := range records {
		acct, err := extractAccountingRecord(v.Bins)
		if err != nil {
			return nil, err
		}

		res = append(res, acct)
	}

	return res, nil
}

func extractAccountingRecord(bins aerospike.BinMap) (types.AccountingRecord, error) {
	acct := types.AccountingRecord{}

	var err error

	// fields ordered according to tacplus log files

	acct.ID, err = extractString(bins, "id")
	if err != nil {
		return acct, err
	}

	ts, err := extractInt(bins, "ts")
	if err != nil {
		return acct, err
	}
	tm := time.Unix(int64(ts), 0)
	acct.Timestamp = tm.Format(types.TimeFormatSeconds)

	acct.DeviceIP, err = extractString(bins, "device_ip")
	if err != nil {
		return acct, err
	}

	acct.DeviceName, err = extractString(bins, "device_name")
	if err != nil {
		return acct, err
	}

	acct.AccountName, err = extractString(bins, "account")
	if err != nil {
		return acct, err
	}

	acct.UserIP, err = extractString(bins, "user_ip")
	if err != nil {
		return acct, err
	}

	acct.UserFQDN, err = extractString(bins, "user_fqdn")
	if err != nil {
		return acct, err
	}

	acct.Command, err = extractString(bins, "command")
	if err != nil {
		return acct, err
	}

	return acct, nil
}
