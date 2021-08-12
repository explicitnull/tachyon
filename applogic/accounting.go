package applogic

import (
	"tacacs-webconsole/repository"
	"tacacs-webconsole/types"
	"time"

	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

func SearchAccounting(le *logrus.Entry, field, value string, begin, end time.Time, aclient *aerospike.Client) []types.AccountingRecord {
	items := make([]types.AccountingRecord, 0)

	var err error

	if value != "" {
		items, err = repository.GetAccountingWithEqualFilter(le, aclient, field, value)
		if err != nil {
			le.WithError(err).Error("searching accounting failed")
			return nil
		}
	} else {
		items, err = repository.GetAccountingWithTimeFilter(le, aclient, begin, end)
		if err != nil {
			le.WithError(err).Error("searching accounting failed")
			return nil
		}
	}

	return items
}
