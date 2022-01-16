package applogic

import (
	"tachyon/repository"
	"tachyon/types"
	"time"

	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

const acctOffset = 60

func ShowAccounting(le *logrus.Entry, aclient *aerospike.Client) ([]types.AccountingRecord, error) {
	now := time.Now()
	begin := now.Add(-acctOffset * time.Minute)
	end := now

	le.Debugf("applogic begin: %s, end: %s", begin, end)

	items, err := repository.GetAccountingWithTimeFilter(le, aclient, begin, end)
	if err != nil {
		le.WithError(err).Error("getting accounting failed")
		return nil, err
	}

	return items, nil
}

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
