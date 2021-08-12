package applogic

import (
	"tacacs-webconsole/repository"
	"tacacs-webconsole/types"

	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

func SearchAccounting(le *logrus.Entry, field, value, from, to string, aclient *aerospike.Client) []types.AccountingRecord {
	items := make([]types.AccountingRecord, 0)

	var err error

	if value != "" {
		items, err = repository.GetAccountingWithEqualFilter(le, aclient, field, value)
		if err != nil {
			le.WithError(err).Error("searching accounting failed")
			return nil
		}
	} else {
		items, err = repository.GetAccountingWithTimeFilter(le, aclient, from, to)
		if err != nil {
			le.WithError(err).Error("searching accounting failed")
			return nil
		}
	}

	return items
}
